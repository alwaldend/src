class Solution {
public:
    int titleToNumber(string columnTitle) {
        int ans{0};
        for (const char& ch : columnTitle) {
            ans = (ans * 26) + (ch - 'A') + 1;
        }
        return ans;
    }
};