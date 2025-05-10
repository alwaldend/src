class Solution {
public:
    vector<int> twoSum(vector<int>& nums, int target) {
        std::map<int, int> numToIndex;
        int length = nums.size();
        std::vector<int> ans;
        for (int i = 0; i < length; ++i) {
            int num = nums[i];
            int diff = target - num;
            if (numToIndex.contains(diff)) {
                ans = {numToIndex[diff], i};
                break;
            }
            numToIndex[num] = i;
        }
        return ans;
    }
};