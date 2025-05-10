class Solution {
public:
    int countNicePairs(vector<int>& nums) {
        std::unordered_map<int, int> diffs{};
        for (const int& num : nums) {
            int revNum{};
            int tempNum{num};
            while(tempNum) {
                revNum = (revNum * 10) + (tempNum % 10);
                tempNum /= 10;
            }
            ++diffs[num - revNum];
        }
        int ans{};
        double mod{std::pow(10, 9) + 7};
        for (const auto& [num, count] : diffs) {
            long pairs{(1L * count * count - count) / 2};
            ans = std::fmod(ans + pairs, mod);
        }
        return ans;
    }

};