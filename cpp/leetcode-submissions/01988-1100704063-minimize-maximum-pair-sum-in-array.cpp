class Solution {
public:
    int minPairSum(vector<int>& nums) {
        std::sort(nums.begin(), nums.end());
        unsigned long length {nums.size()};
        int maxSum {0};
        for (int i = 0; i < length / 2; ++i) {
            maxSum = max(maxSum, nums[i] + nums[length - i - 1]);
        }
        return maxSum;
    }
};