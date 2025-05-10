class Solution {
public:
    string removeStars(string s) {
        std::string ans;
        int remove {0};
        for (int i = s.size() - 1; i >= 0; --i) {
            char c = s[i];
            if (c == '*') {
                ++remove;
            } else if (remove == 0) {
                ans += c;
            } else {
                --remove;
            }
        }
        std::reverse(ans.begin(), ans.end());
        return ans;
    }
};