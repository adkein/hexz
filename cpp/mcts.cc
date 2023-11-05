#include "mcts.h"

#include <absl/log/absl_log.h>
#include <absl/status/statusor.h>
#include <torch/torch.h>

#include <algorithm>
#include <cassert>
#include <cmath>
#include <ostream>
#include <sstream>
#include <vector>

#include "base.h"
#include "perfm.h"

namespace hexz {

float Node::uct_c = 5.0;

Node::Node(Node* parent, int turn, Move move)
    : parent_{parent},
      turn_{turn},
      move_{move},
      flat_idx_{move.typ * 11 * 10 + move.r * 10 + move.c} {}

int Node::NumVisitedChildren() const noexcept {
  int n = 0;
  for (const auto& c : children_) {
    if (c->visit_count_ > 0) {
      n++;
    }
  }
  return n;
}

std::string Node::Stats() const {
  return absl::StrCat("nchildren:", children_.size(),
                      " visited_children:", NumVisitedChildren(),
                      " wins:", wins_);
}

float Node::Puct() const noexcept {
  float q = 0.0;
  if (visit_count_ > 0) {
    q = wins_ / visit_count_;
  }
  float pr = parent_->move_probs_[flat_idx_];
  return q + Node::uct_c * pr * std::sqrt(parent_->visit_count_) /
                 (1 + visit_count_);
}

Node* Node::MaxPuctChild() const {
  Perfm::Scope ps(Perfm::MaxPuctChild);
  if (children_.empty()) {
    return nullptr;
  }
  int best_i = 0;
  float best_puct = children_[best_i]->Puct();
  for (int i = 1; i < children_.size(); i++) {
    float puct = children_[i]->Puct();
    if (puct > best_puct) {
      best_i = i;
      best_puct = puct;
    }
  }
  return children_[best_i].get();
}

std::unique_ptr<Node> Node::MostVisitedChildAsRoot() {
  assert(!children_.empty());
  int best_i = 0;
  for (int i = 1; i < children_.size(); i++) {
    if (children_[i]->visit_count_ > children_[best_i]->visit_count_) {
      best_i = i;
    }
  }
  std::unique_ptr<Node> best_child = std::move(children_[best_i]);
  best_child->parent_ = nullptr;
  return best_child;
}

void Node::Backpropagate(float result) {
  Node* n = this;
  while (n != nullptr) {
    n->visit_count_++;
    if (n->turn() == 0 && result > 0) {
      n->wins_ += result;
    } else if (n->turn() == 1 && result < 0) {
      n->wins_ += -result;
    }
    n = n->parent_;
  }
}

void Node::CreateChildren(int turn, const std::vector<Move>& moves) {
  ABSL_DCHECK(children_.empty())
      << "Must not add children if they already exist.";
  children_.reserve(moves.size());
  for (int i = 0; i < moves.size(); i++) {
    children_.emplace_back(std::make_unique<Node>(this, turn, moves[i]));
  }
}

void Node::AppendDebugString(std::ostream& os,
                             const std::string& indent) const {
  os << indent << "Node(\n";
  os << indent << "  turn: " << turn_ << "\n";
  os << indent << "  move: (" << move_.typ << ", " << move_.r << ", " << move_.c
     << ", " << move_.value << ")\n";
  os << indent << "  wins: " << wins_ << "\n";
  os << indent << "  visit_count: " << visit_count_ << "\n";
  if (parent_ != nullptr) {
    os << indent << "  puct: " << Puct() << "\n";
  }
  if (!children_.empty()) {
    os << indent << "  children: [\n";
    for (const auto& c : children_) {
      c->AppendDebugString(os, indent + "    ");
      os << indent << ",\n";
    }
    os << indent << "  ]\n";
  }
  os << indent << ")";
}

std::string Node::DebugString() const {
  std::ostringstream os;
  AppendDebugString(os, "");
  return os.str();
}

TorchModel::Prediction TorchModel::Predict(int player, const Board& board) {
  Perfm::Scope ps(Perfm::Predict);
  torch::NoGradGuard no_grad;
  auto input = board.Tensor(player).unsqueeze(0);
  std::vector<torch::jit::IValue> inputs{
      input,
  };
  auto output = module_.forward(inputs);

  // The model should output two values: the move likelihoods as a [1, 220]
  // tensor of logits, and a single float value prediction.
  assert(output.isTuple());
  const auto output_tuple = output.toTuple();
  const auto logits = output_tuple->elements()[0].toTensor();
  const auto dim = logits.sizes();
  assert(dim.size() == 2 && dim[0] == 1 && dim[1] == 2 * 11 * 10);
  const auto value = output_tuple->elements()[1].toTensor().item<float>();
  return Prediction{
      .move_probs = logits.reshape({2, 11, 10}).exp(),
      .value = value,
  };
}

NeuralMCTS::NeuralMCTS(Model& model, const Config& config)
    : model_{model},
      runs_per_move_{config.runs_per_move},
      max_moves_per_game_{config.max_moves_per_game},
      runs_per_move_gradient_{config.runs_per_move_gradient} {}

bool NeuralMCTS::Run(Node& root, Board& board) {
  Perfm::Scope ps(Perfm::NeuralMCTS_Run);
  Node* n = &root;
  // Move to leaf node.
  auto t_start = UnixMicros();
  {
    Perfm::Scope ps(Perfm::FindLeaf);
    while (!n->IsLeaf()) {
      Node* child = n->MaxPuctChild();
      board.MakeMove(n->turn(), child->move());
      n = child;
    }
  }
  // Expand leaf node. Usually it's the turn as indicated by the node.
  int turn = n->turn();
  auto moves = board.NextMoves(turn);
  if (moves.empty()) {
    // Player has no valid moves left. Try if opponent can proceed.
    turn = 1 - turn;
    n->SetTurn(turn);
    moves = board.NextMoves(turn);
  }
  if (moves.empty()) {
    // No player can make a move => game over.
    n->Backpropagate(board.Result());
    return n != &root;  // Return if we made any progress at all in this run.
  }
  // Initially we assume that it's the opponent's turn. If that turns out to be
  // false, the turn gets updated when trying to find next moves (see above).
  n->CreateChildren(1 - turn, moves);
  n->ShuffleChildren();  // Avoid selection bias.
  auto pred = model_.Predict(n->turn(), board);
  n->SetMoveProbs(pred.move_probs);
  // Backpropagate the model prediction. Need to reorient it s.t. 1 means player
  // 0 won.
  n->Backpropagate(n->turn() == 0 ? pred.value : -pred.value);
  return true;
}

int NeuralMCTS::NumRuns(int move) const noexcept {
  return std::max(
      1, static_cast<int>(std::round(runs_per_move_ *
                                     (1 + move * runs_per_move_gradient_))));
}

absl::StatusOr<std::vector<hexzpb::TrainingExample>> NeuralMCTS::PlayGame(
    Board& board, int max_runtime_seconds) {
  std::vector<hexzpb::TrainingExample> examples;
  int64_t started_micros = UnixMicros();
  int n = 0;
  // Root's children have the player whose turn it actually is.
  // Every game starts with player 0, so root must use player 1.
  auto root =
      std::make_unique<Node>(nullptr, /*turn=*/0, Move{-1, -1, -1, -1.0});
  float result = 0.0;
  bool game_over = false;
  const int64_t max_micros =
      max_runtime_seconds > 0
          ? started_micros +
                static_cast<int64_t>(max_runtime_seconds) * 1'000'000
          : std::numeric_limits<int64_t>::max();
  for (; n < max_moves_per_game_; n++) {
    int64_t move_started = UnixMicros();
    if (move_started > max_micros) {
      return absl::DeadlineExceededError(
          "max_runtime_seconds exceeded before the game was finished");
    }
    bool progress = true;
    // The first moves are the most important ones. Think twice as hard for
    // those.
    const int runs = NumRuns(n);
    for (int i = 0; i < runs && progress; i++) {
      Board b(board);
      progress = Run(*root, b);
    }
    if (root->IsLeaf()) {
      result = board.Result();
      ABSL_LOG(INFO) << "Game over. Final score: " << board.Score()
                     << ". Result: " << result;
      game_over = true;
      break;
    }

    // Add example.
    hexzpb::TrainingExample example;
    int64_t move_ready = UnixMicros();
    example.set_unix_micros(move_ready);
    example.set_turn(root->turn());
    example.mutable_stats()->set_duration_micros(move_ready - move_started);
    example.mutable_stats()->set_move(n);
    example.mutable_stats()->set_valid_moves(root->NumChildren());
    example.mutable_stats()->set_visit_count(root->visit_count());
    example.set_encoding(hexzpb::TrainingExample::PYTORCH);
    // The board in the example must be viewed from the perspective of the
    // player whose turn it is.
    auto enc_b = torch::pickle_save(board.Tensor(root->turn()));
    example.mutable_board()->assign(enc_b.begin(), enc_b.end());
    auto enc_pr = torch::pickle_save(root->NormVisitCounts());
    example.mutable_move_probs()->assign(enc_pr.begin(), enc_pr.end());
    examples.push_back(example);

    std::string stats = root->Stats();
    if (n < 5 || n % 10 == 0) {
      ABSL_LOG(INFO) << "Move " << n << " (turn: " << root->turn() << ") ("
                     << board.DebugString() << ") after "
                     << (float)(UnixMicros() - started_micros) / 1000000
                     << "s. stats: " << stats;
    } else {
      ABSL_DLOG(INFO) << "Move " << n << " (turn: " << root->turn() << ") ("
                      << board.DebugString() << ") after "
                      << (float)(UnixMicros() - started_micros) / 1000000
                      << "s. stats: " << stats;
    }

    int turn = root->turn();
    // NOTE: Must not access root after this step!
    std::unique_ptr<Node> best_child = root->MostVisitedChildAsRoot();
    board.MakeMove(turn, best_child->move());
    root = std::move(best_child);
  }
  if (game_over) {
    for (auto& ex : examples) {
      ex.set_result(ex.turn() == 0 ? result : -result);
    }
  }
  return examples;
}

}  // namespace hexz
