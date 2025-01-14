#ifndef __HEXZ_MCTS_H__
#define __HEXZ_MCTS_H__

#include <absl/status/status.h>
#include <absl/status/statusor.h>
#include <torch/script.h>
#include <torch/torch.h>

#include <cassert>
#include <ostream>
#include <string_view>
#include <vector>

#include "base.h"
#include "board.h"
#include "hexz.pb.h"
#include "perfm.h"

namespace hexz {

class Node {
 public:
  Node(Node* parent, int turn, Move move);

  // Simple getters.
  int turn() const { return turn_; }
  const Move& move() const { return move_; }
  float prior() const noexcept { return prior_; }
  int visit_count() const noexcept { return visit_count_; }
  const Node* parent() const noexcept { return parent_; }
  float wins() const noexcept { return wins_; }
  float value() const noexcept { return value_; }
  bool terminal() const noexcept { return terminal_; }
  const std::vector<std::unique_ptr<Node>>& children() const {
    return children_;
  }

  // Returns this node's Q value. If the child has not been visited
  // yet, returns a Q value derived from the parent's Q value.
  float Q() const noexcept;

  // Update the turn.
  void SetTurn(int turn) { turn_ = turn; }

  // Returns this node's PUCT value.
  float Puct() const noexcept;
  // Returns a non-owned pointer to the child with the greatest PUCT value.
  Node* MaxPuctChild() const;
  // Returns the child with the greated visit count,
  // marks it as a root node (by setting its parent_ to nullptr),
  // and clears the vector of child nodes of *this* Node.
  //
  // The Node on which this method was called MUST NOT be used
  // afterwards!
  std::unique_ptr<Node> MostVisitedChildAsRoot();

  // Returns a child as the new root, selected with a probability
  // proportional to its relative visit count.
  // Behaves like MostVisitedChildAsRoot in all other respects.
  std::unique_ptr<Node> SampleChildAsRoot();

  // Backpropagate propagates the given result to this Node and all parents.
  // The given result is interpreted from the perspective of player 0.
  // A value of 1 always means that player 0 won.
  // This is particularly relevant when feeding in model predictions,
  // which are usually from the perspective of the current player.
  void Backpropagate(float result);

  bool IsLeaf() const { return children_.empty(); }

  // Adds children representing the given moves to this node.
  void CreateChildren(int turn, const std::vector<Move>& moves);
  // Shuffles children randomly. This can be used to avoid selection bias.
  void ShuffleChildren();

  // Returns the normalized visit counts (which sum to 1) as a (2, 11, 10)
  // tensor. This value should be used to update the move probs of the model.
  torch::Tensor NormVisitCounts() const;

  // Recursively zeros visit_count_ and wins_ of all nodes in the whole subtree.
  // This method is used for prediction caching: re-use the search tree of
  // the previous move during self-play, but only for re-using the model
  // predictions.
  void ResetTree();

  // Returns a boolean Tensor of shape (2, 11, 10) that is true in each
  // position that represents a valid move and false everywhere else.
  torch::Tensor ActionMask() const;

  // Sets the initial move probabilities ("policy") obtained from the model
  // prediction.
  void SetMoveProbs(torch::Tensor move_probs);
  // Sets the value of this node (typically: as predicted by the model).
  // The value should range between -1 and 1; values greater than 0 predict that
  // it is a win for the player whose turn it is in this node.
  void SetValue(float value) { value_ = value; }
  void SetTerminal(bool b) { terminal_ = b; }

  // Adds Dirichlet noise to the (already set) move probs.
  // The weight parameter must be in the open interval (0, 1).
  // concentration is the Dirichlet concentration ("alpha") parameter.
  void AddDirichletNoise(float weight, float concentration);

  // Returns the number of children that had a nonzero visit_count.
  int NumVisitedChildren() const noexcept;
  int NumChildren() const noexcept { return children_.size(); }
  std::string Stats() const;
  // Weight of the exploration term.
  // Must only be modified at program startup.
  static float uct_c;

  std::string DebugString() const;

 private:
  void AppendDebugString(std::ostream& os, const std::string& indent) const;

  Node* parent_;
  int turn_;
  Move move_;
  float prior_ = 0.0;
  float wins_ = 0.0;
  int visit_count_ = 0;
  float value_ = 0.0;
  bool terminal_ = false;
  std::vector<std::unique_ptr<Node>> children_;
};

// Writes the substree starting at root to path as a GraphViz .dot file.
absl::Status WriteDotGraph(const Node& root, const std::string& path);

// Interface class for a model than can make predictions for
// move likelihoods and board evaluations.
class Model {
 public:
  struct Prediction {
    torch::Tensor move_probs;
    float value;
  };
  virtual Prediction Predict(const Board& board, const Node& node) = 0;
  virtual ~Model() = default;
};

// Implementation of Model that uses an actual PyTorch ScriptModule.
// This is the implementation that should be used everywhere outside of tests.
class TorchModel : public Model {
  using Prediction = Model::Prediction;

 public:
  explicit TorchModel(torch::jit::Module module) : module_{module} {}
  TorchModel(hexzpb::ModelKey key, torch::jit::Module module)
      : key_{key}, module_{module} {}
  Prediction Predict(const Board& board, const Node& node) override;

  const hexzpb::ModelKey& Key() const { return key_; }
  torch::jit::Module& Module() { return module_; }

 private:
  hexzpb::ModelKey key_;
  torch::jit::Module module_;
};

class NeuralMCTS {
 public:
  // The model is not owned. Owners of the NeuralMCTS instance must ensure it
  // outlives this instance.
  NeuralMCTS(Model& model, const Config& config);

  absl::StatusOr<std::vector<hexzpb::TrainingExample>> PlayGame(
      Board& board, int max_runtime_seconds);

  // Returns a pair of (N, record_example), where N is the number of iterations
  // for the MCTS run and record_example indicates whether the resulting example
  // should be recorded.
  // This is only relevant in a full PlayGame cycle.
  std::pair<int, bool> NumRuns(int move) const noexcept;

  // Executes a single run of the MCTS algorithm, starting at root.
  // This method should be called during "normal" play, i.e. outside of
  // training.
  bool Run(Node& root, const Board& board);
  // RunReusingTree is a variant of Run that should be used for self-play
  // during training.
  // If add_noise is true, Dirichlet noise will be added to the root node's
  // move probs.
  bool RunReusingTree(Node& root, const Board& b, bool add_noise);

  // SuggestMove returns the best move suggestion that the NeuralMCTS algorithm
  // comes up with in think_time_millis milliseconds.
  absl::StatusOr<std::unique_ptr<Node>> SuggestMove(int player,
                                                    const Board& board,
                                                    int think_time_millis);

 private:
  Config config_;
  // Not owned.
  Model& model_;
};

}  // namespace hexz
#endif  // __HEXZ_MCTS_H__