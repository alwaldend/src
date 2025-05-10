class Solution {
public:
    string removeStars(string s) {
        std::string ans;
        unsigned long length {s.size()};
        int remove {0};
        for (int i = length - 1; i >= 0; --i) {
            char c = s[i];
            if (c == '*') {
                ++remove;
                continue;
            }
            if (remove == 0) {
                ans += c;
            } else {
                --remove;
            }
        }
        std::reverse(ans.begin(), ans.end());
        return ans;
    }
};