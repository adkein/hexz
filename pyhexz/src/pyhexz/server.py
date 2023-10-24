import base64
import collections
from contextlib import contextmanager
import datetime
from flask import Flask, current_app, make_response, request
import json
import logging
import os
import pytz
import queue
import time
import typing

from pyhexz.board import Board
from pyhexz.config import TrainingConfig
from pyhexz.errors import HexzError
from pyhexz.hexz import HexzNeuralNetwork, NeuralMCTS
from pyhexz import hexz_pb2
from pyhexz import sconv
from pyhexz import svg
from pyhexz.modelserver import LocalModelRepository
from pyhexz import training


def suggest_move(
    model: HexzNeuralNetwork, state: hexz_pb2.GameEngineState, think_time: float
) -> tuple[dict[str, typing.Any], int]:
    """Runs the ML model to obtain a move suggestion.

    Returns:
        A tuple of (SuggestMoveResponse JSON response, http status code).

    Args:
        model: the model to use in neural MCTS.
        state: the game state, including the player whose turn it is, and the board.
        think_time: thinking time in seconds.
    """
    board = Board.from_numpy(sconv.convert_board(state.flagz.board))
    try:
        board.validate()
    except ValueError as e:
        return str(e), 400
    turn = state.flagz.board.turn - 1  # Go impl uses (1, 2), we use (0, 1).
    print(f"Board info: flags:{board.nflags}, turn:{turn}")
    mcts = NeuralMCTS(
        board, model, game_id=time.strftime("CPU-%Y%m%d-%H%M%S"), turn=turn
    )
    n = 0
    started = time.perf_counter()
    while True:
        if n & 63 == 0 and time.perf_counter() - started >= think_time:
            break
        mcts.run()
        n += 1
    best_child = mcts.root.best_child()
    if not best_child:
        raise ValueError("No next move")
    typ, r, c, _ = best_child.move
    # Return a SuggestMoveResponse JSON.
    resp = {
        "move": {
            "move": state.flagz.board.move,
            "row": int(r),
            "col": int(c),
            "type": 0 if typ == 1 else 5,  # 5==Flag
        }
    }
    svg.export(
        "./latest.html",
        [board, board],
        [f"Model move probs (value: {mcts.value:.3f})", "MCTS move likelihoods"],
        [mcts.root.move_probs, mcts.root.move_likelihoods()],
    )
    print(f"child qs: {[(c.wins, c.visit_count) for c in mcts.root.children]}")
    print(
        f"suggested move: {resp}, iterations: {n}, tree size: {mcts.size()}, best child: vc:{best_child.visit_count} {best_child.wins} {best_child.puct()}."
    )
    return resp, 200


def create_app():
    app = Flask(__name__)
    app.logger.setLevel(logging.INFO)
    config = TrainingConfig.from_env()
    app.model_queue = queue.SimpleQueue()
    app.hexz_config = config

    if config.model_repo_base_dir:
        app.training_task_queue = queue.SimpleQueue()
        app.model_repo = LocalModelRepository(config.model_repo_base_dir)
        if not config.model_name:
            raise HexzError(
                "No model_name specified. Did you forget to set HEXZ_MODEL_NAME?"
            )
        t = training.TrainingTask(
            repo=app.model_repo,
            config=config,
            queue=app.training_task_queue,
        )
        t.daemon = True
        t.start()
        app.training_task = t
        app.logger.info("Started training task")

    @contextmanager
    def get_model():
        """Gets a model from the model_queue, if one is available. Otherwise, loads a new one.

        The idea is that each thread processing a request needs exclusive access to a model.
        """
        try:
            model = current_app.model_queue.get_nowait()
        except queue.Empty:
            model = app.model_repo.get_model(app.hexz_config.model_name)
        try:
            yield model
        finally:
            current_app.model_queue.put(model)

    # The path /hexz/cpu/suggest must be identical to the one the Go client uses.
    @app.post("/hexz/cpu/suggest")
    def suggestmove():
        req = request.json
        think_time_ns = req["maxThinkTime"]
        enc_state = req["gameEngineState"]
        ge_state = hexz_pb2.GameEngineState().FromString(base64.b64decode(enc_state))
        if not current_app.model_path:
            return "Missing model_path", 500
        with get_model() as model:
            return suggest_move(model, ge_state, think_time_ns / 1e9)

    @app.get("/")
    def index():
        now = datetime.datetime.now(tz=pytz.UTC).isoformat()
        resp = make_response(f"Hello from Python hexz at {now}!\n")
        resp.headers["Content-Type"] = "text/plain"
        return resp

    @app.get("/status")
    def status():
        now = datetime.datetime.now(tz=pytz.UTC).isoformat(timespec="milliseconds")
        resp = make_response(
            json.dumps(
                {
                    "timestamp": now,
                    "model": {
                        "path": current_app.model_path,
                        "pool_size": current_app.model_queue.qsize(),
                    },
                },
                indent=True,
            )
        )
        resp.headers["Content-Type"] = "application/json"
        return resp

    @app.post("/examples")
    def examples():
        reply_q = queue.SimpleQueue()
        current_app.training_task_queue.put({
            "type": "AddTrainingExamplesRequest",
            "data": request.data,
            "reply_q": reply_q,
            })
        return reply_q.get(timeout=5)

    @app.get("/models")
    def model():
        reply_q = queue.SimpleQueue()
        current_app.training_task_queue.put({
            "type": "GetModelKey",
            "reply_q": reply_q,
            })
        model_key: hexz_pb2.ModelKey = reply_q.get(timeout=5)
        return {
            "model_name": model_key.name,
            "checkpoint": model_key.checkpoint,
        }

    # For debugging bad requests:
    # @app.errorhandler(400)
    # def handle_bad_request(e):
    #     print("Really bad request!!!", e)
    #     return 'bad request!', 400

    return app
