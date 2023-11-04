#ifndef __HEXZ_BOARD_H__
#define __HEXZ_BOARD_H__

#include <absl/strings/str_cat.h>
#include <torch/torch.h>

#include <utility>
#include <vector>

namespace hexz {

namespace internal {

// Hashable indexes into the hexz board.
struct Idx {
  int r;
  int c;
  bool operator==(const Idx& other) const {
    return r == other.r && c == other.c;
  }
  bool IsValid() const noexcept {
    return r >= 0 && r < 11 && c >= 0 && c < 10 - r % 2;
  }
};

const std::vector<Idx>& NeighborsOf(const Idx& k);

}  // namespace internal
}  // namespace hexz

namespace std {
template <>
struct hash<hexz::internal::Idx> {
  size_t operator()(const hexz::internal::Idx& k) const noexcept {
    size_t h1 = hash<int>{}(k.r);
    size_t h2 = hash<int>{}(k.c);
    return h1 ^ (h2 << 1);
  }
};
}  // namespace std

namespace hexz {

struct Move {
  int typ;
  int r;
  int c;
  float value;
  std::string DebugString() const {
    return absl::StrCat("Move(", typ, ", ", r, ", ", c, ", ", value, ")");
  }
};

// Torch representation of a hexz board.
// A board is represented by an (9, 11, 10) tensor. Each 11x10 channel is
// a one-hot encoding of the presence of specific type of piece/obstacle/etc.
// The channels are:
//
// * 0: flags by P0
// * 1: cell value 1-5 for P0
// * 2: cells blocked for P0 (any occupied cell or a cell next to a 5)
// * 3: next value for P0
// * 4: flags by P1
// * 5: cell value 1-5 for P1
// * 6: cells blocked for P1
// * 7: next value for P1
// * 8: grass cells with value 1-5
//
// An action is specified by a (2, 11, 10) numpy array. The first 11x10 channel
// represents a flag move, the second one represents a regular cell move. A
// flag move must have a single 1 set, a normal move must have a single value
// 1-5 set.
class Board {
 public:
  enum Channel {
    kFlag = 0,
    kValue = 1,
    kBlocked = 2,
    kNextValue = 3,
    kGrass = 8
  };

  Board();
  static Board RandomBoard();

  // Copy c'tor.
  Board(const Board& other);

  std::pair<float, float> Score() const;
  float Result() const;

  torch::Tensor Tensor(int player) const;

  int Flags(int player) const;

  void MakeMove(int player, const Move& move);
  std::vector<Move> NextMoves(int player) const;

  std::string DebugString() const;

  // CellValue returns the board's value in cell (r, c) and channel ch.
  // This method can be used to access any cell in any channel of the board.
  // It is not optimized for performance. Get a PyTorch accessor in that case.
  float CellValue(int player, Channel ch, int r, int c) {
    int ch_idx = static_cast<int>(ch == kGrass ? ch : ch + 4 * player);
    return b_.index({ch_idx, r, c}).item<float>();
  }
  // These methods can be used to create specific board setups, e.g. for
  // testing. They place the rock or grass unconditionally on the given field,
  // without any checks or propagation. Should only be used on an otherwise
  // empty board.
  void SetCellValue(int player, Channel ch, int r, int c, float value) {
    int ch_idx = static_cast<int>(ch == kGrass ? ch : ch + 4 * player);
    b_.index_put_({ch_idx, r, c}, value);
  }
  void SetRemainingFlags(int player, int n_flags) { nflags_[player] = n_flags; }

 private:
  torch::Tensor b_;
  int nflags_[2];
};

}  // namespace hexz
#endif  // __HEXZ_BOARD_H__
