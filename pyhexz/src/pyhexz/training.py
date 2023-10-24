"""Classes used for training in a training server."""

from collections.abc import Mapping
from google.protobuf.message import DecodeError
import io
import numpy as np
import queue
import threading
import time
import torch
from torch import nn
import torch.nn.functional as F
from typing import Any, Optional
from pyhexz.config import TrainingConfig
from pyhexz import hexz_pb2
from pyhexz.errors import HexzError
from pyhexz.modelserver import ModelRepository


class InMemoryDataset(torch.utils.data.Dataset):
    def __init__(self):
        self.examples = []

    def __getitem__(self, k):
        return self.examples[k]

    def __len__(self):
        return len(self.examples)

    def append(self, example):
        self.examples.append(example)


class TrainingRunnerTask(threading.Thread):
    def __init__(
        self,
        repo: ModelRepository,
        model_key: hexz_pb2.ModelKey,
        config: TrainingConfig,
        batch: InMemoryDataset,
        done_queue: queue.SimpleQueue,
    ):
        super().__init__()
        self.repo = repo
        self.model_key = model_key
        self.done_queue = done_queue
        self.config = config
        self.batch = batch

    def run(self):
        reply = {
            "type": "TrainingRunnerTask",
            "status": "ERROR",
        }
        try:
            started = time.perf_counter()
            self.run_training()
            elapsed = time.perf_counter() - started
            reply["status"] = "OK"
            reply["model_key"] = self.model_key
            reply["elapsed"] = elapsed
        except Exception as e:
            reply["error"] = str(e)
        finally:
            self.done_queue.put(reply)

    def run_training(self):
        batch = self.batch
        device = self.config.device
        model = self.repo.get_model(self.model_key.name, self.model_key.checkpoint)
        model = model.to(device)
        num_epochs = self.config.num_epochs
        batch_size = self.config.batch_size
        loader = torch.utils.data.DataLoader(
            dataset=batch,
            batch_size=batch_size,
            shuffle=self.config.shuffle,
            pin_memory=self.config.pin_memory,
        )
        print(f"torch.autograd.gradcheck:  {torch.autograd.gradcheck}")
        for X, (y_pr, y_val) in loader:
            print(f"Shape of X [N, C, H, W]: {X.shape}: {X.dtype}")
            print(f"Shape of y_pr [N, C, H, W]: {y_pr.shape}: {y_pr.dtype}")
            print(f"Shape of y_val [N, V]: {y_val.shape}: {y_val.dtype}")
            break

        pr_loss_fn = nn.CrossEntropyLoss()
        val_loss_fn = nn.MSELoss()

        optimizer = torch.optim.Adam(model.parameters(), lr=1e-3, weight_decay=1e-4)

        for _ in range(num_epochs):
            for X, (y_pr, y_val) in loader:
                # Send to device.
                X = X.to(device)
                y_pr = y_pr.to(device)
                y_val = y_val.to(device)

                # Predict
                pred_pr, pred_val = model(X)

                # Compute loss
                pr_loss = pr_loss_fn(pred_pr.flatten(1), y_pr.flatten(1))
                val_loss = val_loss_fn(pred_val, y_val)
                loss = pr_loss + val_loss

                # Backpropagation
                loss.backward()
                optimizer.step()
                optimizer.zero_grad(set_to_none=True)
        self.model_key.checkpoint += 1
        self.repo.store_model(
            name=self.model_key.name,
            checkpoint=self.model_key.checkpoint,
            model=model,
        )


class TrainingTask(threading.Thread):
    """TrainingTask runs as a separate thread, accepts training data and
    spaws a separate thread or process to train the model and save new
    checkpoints."""

    _STATE_ACCEPTING = "ACCEPTING"
    _STATE_TRAINING = "TRAINING"

    def __init__(
        self, repo: ModelRepository, config: TrainingConfig, queue: queue.SimpleQueue
    ):
        super().__init__()
        self.repo = repo
        checkpoint = repo.get_latest_checkpoint(config.model_name)
        self.model_key = hexz_pb2.ModelKey(
            name=config.model_name,
            checkpoint=checkpoint,
        )
        self.queue = queue
        self.config = config
        self.batch = InMemoryDataset()
        self.training_runner_task: Optional[TrainingRunnerTask] = None
        self.state = self._STATE_ACCEPTING

    def start_training_runner_task(self):
        t = TrainingRunnerTask(
            self.repo, self.model_key, self.config, self.batch, self.queue
        )
        t.start()
        self.training_runner_task = t
        self.state = self._STATE_TRAINING
        self.batch = InMemoryDataset()

    def handle_training_examples(self, msg: Mapping[str, Any]):
        data = msg["data"]
        reply_q = msg["reply_q"]
        batch = self.batch
        if self.state != self._STATE_ACCEPTING:
            resp = hexz_pb2.AddTrainingExamplesResponse(
                status=hexz_pb2.AddTrainingExamplesResponse.REJECTED,
                error_message=f"Server is currently not accepting examples",
            )
            reply_q.put(resp.SerializeToString())
            return
        try:
            req = hexz_pb2.AddTrainingExamplesRequest.FromString(data)
        except DecodeError as e:
            resp = hexz_pb2.AddTrainingExamplesResponse(
                status=hexz_pb2.AddTrainingExamplesResponse.REJECTED,
                error_message=f"Failed to decode AddTrainingExamplesResponse: {e}",
            )
            reply_q.put(resp.SerializeToString())
            return
        if req.model_key != self.model_key:
            resp = hexz_pb2.AddTrainingExamplesResponse(
                status=hexz_pb2.AddTrainingExamplesResponse.REJECTED,
                latest_model=self.model_key,
            )
            reply_q.put(resp.SerializeToString())
            return
        resp = hexz_pb2.AddTrainingExamplesResponse(
            status=hexz_pb2.AddTrainingExamplesResponse.ACCEPTED,
            latest_model=self.model_key,
        )
        reply_q.put(resp.SerializeToString())

        # Add examples and run training if we have a full batch.
        lim = min(len(req.examples), self.config.batch_size - len(batch))
        for ex in req.examples[:lim]:
            board = np.load(io.BytesIO(ex.board))
            pr = np.load(io.BytesIO(ex.move_probs))
            val = np.array([ex.result], dtype=np.float32)
            batch.append((board, (pr, val)))

        if len(batch) >= self.config.batch_size:
            self.start_training_runner_task()

    def handle_training_result(self, msg: Mapping[str, Any]):
        # Start a new batch, whether training was successful or not.
        self.state = self._STATE_ACCEPTING
        if msg["status"] == "OK":
            self.model_key = msg["model_key"]
            elapsed = msg["elapsed"]
            print(
                f"Finished training batch of size {self.config.batch_size} for {self.config.num_epochs} epochs in {elapsed:.1f}s."
            )
            print(f"Updated model key to {self.model_key.name}:{self.model_key.checkpoint}.")
        else:
            print(f"Training failed: {msg.get('error')}")

    def handle_get_model_key(self, msg):
        reply_q = msg["reply_q"]
        reply_q.put(self.model_key)
        
    def run(self):
        while True:
            # We use bytes as input and output in the communication with this
            # TrainingTask to simplify serialization when running in a separate process.
            msg = self.queue.get()
            if msg["type"] == "TrainingRunnerTask":
                self.handle_training_result(msg)
            elif msg["type"] == "AddTrainingExamplesRequest":
                self.handle_training_examples(msg)
            elif msg["type"] == "GetModelKey":
                self.handle_get_model_key(msg)
            else:
                raise HexzError(
                    f"TrainingTask: received unknown message type {msg['type']}"
                )
