class Solution {
public:
    vector<int> getElementsOnlyInFirstList(vector<int>& nums1, vector<int>& nums2) {
        unordered_set<int> onlyInNums1;
        for (int num : nums1) {
            bool existInNums2 = false;
            for (int x : nums2) {
                if (x == num) {
                    existInNums2 = true;
                    break;
                }
            }
            if (!existInNums2) {
                onlyInNums1.insert(num);
            }
        }
        return vector<int> (onlyInNums1.begin(), onlyInNums1.end());
    }
    
    vector<vector<int>> findDifference(vector<int>& nums1, vector<int>& nums2) {
        return {getElementsOnlyInFirstList(nums1, nums2), getElementsOnlyInFirstList(nums2, nums1)};
    }
};