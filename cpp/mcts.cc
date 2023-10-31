#include "mcts.h"

#include <torch/torch.h>

#include <algorithm>
#include <cassert>
#include <cmath>
#include <vector>

#include "util.h"

namespace hexz {

namespace {
std::mt19937 rng{std::random_device{}()};
}

Node::Node(Node* parent, int player, Move move)
    : parent_{parent}, player_{player}, move_{move} {}

float Node::Puct() const {
  Perfm::Scope ps(Perfm::Puct);

  constexpr float uct_c = 5.0;  // Constant weight of the exploration term.
  float q = 0.0;
  if (visit_count_ > 0) {
    q = wins_ / visit_count_;
  }
  float pr = 0.0;
  {
    Perfm::Scope ps(Perfm::PuctMoveProbs);
    pr =
        parent_->move_probs_.index({move_.typ, move_.r, move_.c}).item<float>();
  }
  return q + uct_c * pr * std::sqrt(parent_->visit_count_) / (1 + visit_count_);
}

Node* Node::MaxPuctChild() {
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
    if (n->Player() == 0 && result > 0) {
      n->wins_ += result;
    } else if (n->Player() == 1 && result < 0) {
      n->wins_ += -result;
    }
    n = n->parent_;
  }
}

void Node::CreateChildren(int player, const std::vector<Move>& moves) {
  assert(children_.empty());
  children_.reserve(moves.size());
  for (int i = 0; i < moves.size(); i++) {
    children_.emplace_back(std::make_unique<Node>(this, player, moves[i]));
  }
  std::shuffle(children_.begin(), children_.end(), rng);
}

NeuralMCTS::NeuralMCTS(torch::jit::script::Module module) : module_{module} {
  module_.eval();
}

NeuralMCTS::Prediction NeuralMCTS::Predict(int player, const Board& board) {
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
  return NeuralMCTS::Prediction{
      .move_probs = logits.reshape({2, 11, 10}).exp(),
      .value = value,
  };
}

bool NeuralMCTS::Run(Node* root, const Board& board) {
  Board b(board);
  Node* n = root;
  // Move to leaf node.
  auto t_start = UnixMicros();
  {
    Perfm::Scope ps(Perfm::FindLeaf);
    while (!n->IsLeaf()) {
      n = n->MaxPuctChild();
      b.MakeMove(n->Player(), n->GetMove());
    }
  }
  // Expand leaf node. Usually it's the opponent's turn.
  int player = 1 - n->Player();
  auto moves = b.NextMoves(player);
  if (moves.empty()) {
    // Opponent has no valid moves left. Try other player.
    player = 1 - player;
    moves = b.NextMoves(player);
  }
  if (moves.empty()) {
    // No player can make a move => game over.
    n->Backpropagate(b.Result());
    return n != root;  // Return if we made any progress at all in this run.
  }
  n->CreateChildren(player, moves);
  auto pred = Predict(player, b);
  n->SetMoveProbs(pred.move_probs);
  n->Backpropagate(pred.value);
  return true;
}

std::vector<hexzpb::TrainingExample> NeuralMCTS::PlayGame(const Board& board,
                                                          int runs_per_move,
                                                          int max_moves) {
  Board b(board);
  std::vector<hexzpb::TrainingExample> examples;
  int64_t started_micros = UnixMicros();
  int n = 0;
  int player = 0;
  auto root =
      std::make_unique<Node>(nullptr, 1 - player, Move{-1, -1, -1, -1.0});
  for (; n < max_moves; n++) {
    std::cout << "Move " << n << " after "
              << (float(UnixMicros() - started_micros) / 1000000) << "s\n";
    // Root's children have the player whose turn it actually is.
    // So root has the opposite player.
    bool progress = true;
    // The first moves are the most important ones. Think twice as hard for
    // those.
    int limit = n >= 6 ? runs_per_move : 2 * runs_per_move;
    for (int i = 0; i < limit && progress; i++) {
      progress = Run(root.get(), b);
    }
    if (root->IsLeaf()) {
      // Game over
      break;
    }
    std::unique_ptr<Node> best_child = root->MostVisitedChildAsRoot();
    b.MakeMove(best_child->Player(), best_child->GetMove());
    player = 1 - best_child->Player();
    root = std::move(best_child);
  }
  return examples;
}

}  // namespace hexz
