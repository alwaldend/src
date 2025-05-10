class Solution {
public:
    double findMaxAverage(vector<int>& nums, int k) {
        int length = nums.size();
        double sum = std::accumulate(nums.begin(), nums.begin() + k , 0);
        double maxSum = sum;
        for (int i = k; i < length; ++i) {
            sum += nums[i] - nums[i-k];
            maxSum = max(sum, maxSum);
        }
        return maxSum / k;
    }
};