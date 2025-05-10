class Solution {
public:
    int pivotIndex(vector<int>& nums) {
        int left {0}, right {std::accumulate(nums.begin() + 1, nums.end(), 0)};
        unsigned long length {nums.size()};
        if (right == 0) {
            return 0;
        }
        for (int i = 1; i < length; ++i) {
            left += nums[i-1];
            right -= nums[i];
            if (left == right) {
                return i;
            }
        }
        return -1;
    }
};