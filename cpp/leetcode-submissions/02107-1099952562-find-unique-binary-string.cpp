class Solution {
public:
    string findDifferentBinaryString(vector<string>& nums) {
        unordered_set<int> integers;
        for (string num : nums) {
            integers.insert(stoi(num, 0, 2));
        }

        int n = nums.size();
        string ans;
        for (int num = 0; num <= n; num++) {
            if (integers.find(num) == integers.end()) {
                ans = bitset<16>(num).to_string();
                break;
            }
        }
        return ans.substr(16 - n);
    }
};