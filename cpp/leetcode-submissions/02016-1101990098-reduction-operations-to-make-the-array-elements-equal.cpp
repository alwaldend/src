class Solution {
public:
    int reductionOperations(vector<int>& nums) {
        sort(nums.begin(), nums.end());
        int ans = 0;
        int up = 0;
        auto length {nums.size()};
        for (int i = 1; i < length; ++i) {
            if (nums[i] != nums[i - 1]) {
                ++up;
            }
            ans += up;
        }
        return ans;
    }
};