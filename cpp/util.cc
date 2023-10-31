#include "util.h"

#include <algorithm>
#include <chrono>
#include <string>
#include <vector>

namespace hexz {

Perfm::CumulativeStats Perfm::stats_[Perfm::StatsSize];

void Perfm::Init() {
  for (int i = 0; i < StatsSize; i++) {
    stats_[i].label = static_cast<Perfm::Label>(i);
  }
}

void Perfm::PrintStats() {
  std::vector<Perfm::CumulativeStats> stats;
  int scope_len = 3;
  for (int i = 0; i < StatsSize; i++) {
    if (Perfm::stats_[i].count == 0) continue;
    stats.push_back(Perfm::stats_[i]);
    auto s = Perfm::LabelName(static_cast<Perfm::Label>(i)).size() + 3;
    if (s > scope_len) {
      scope_len = s;
    }
  }
  // Sort by elapsed time, descending.
  std::sort(stats.begin(), stats.end(), [](const auto& lhs, const auto& rhs) {
    return lhs.elapsed_nanos > rhs.elapsed_nanos;
  });

  std::printf("%-*s %10s %10s %12s\n", scope_len, "scope", "total_time",
              "count", "ops/s");

  for (const auto& s : stats) {
    std::printf("%-*s %9.3fs %10lld %12.3f\n", scope_len,
                Perfm::LabelName(s.label).c_str(),
                double(s.elapsed_nanos) / 1e9, s.count,
                s.count * 1e9 / s.elapsed_nanos);
  }
}

std::string GetEnv(const std::string& name) {
  const char* value = std::getenv(name.c_str());
  if (value == nullptr) {
    return "";
  }
  return std::string(value);
}

int GetEnvAsInt(const std::string& name, int default_value) {
  const char* value = std::getenv(name.c_str());
  if (value == nullptr) {
    return default_value;
  }
  return std::atoi(value);
}

int64_t UnixMicros() {
  return std::chrono::duration_cast<std::chrono::microseconds>(
             std::chrono::steady_clock::now().time_since_epoch())
      .count();
}

}  // namespace hexz
