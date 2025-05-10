class Solution {
public:
    void moveZeroes(vector<int>& nums) {
        int length = nums.size();
        int lastNonZeroFoundAt = 0;
        for (int i = 0; i < length; ++i) {
            if (nums[i] != 0) {
                swap(nums[lastNonZeroFoundAt], nums[i]);
                lastNonZeroFoundAt += 1;
            }
        }
    }
};