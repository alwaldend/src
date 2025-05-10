class Solution {
public:
    bool uniqueOccurrences(vector<int>& arr) {
        std::unordered_map<int, int> counts;
        std::unordered_set<int> encountered;
        for (int num : arr) {
            counts[num] += 1;
        }
        for (const auto [num, count] : counts) {
            if (encountered.find(count) == encountered.end()) {
                encountered.insert(count);
            } else {
                return false;
            }
        }
        return true;
    }
};