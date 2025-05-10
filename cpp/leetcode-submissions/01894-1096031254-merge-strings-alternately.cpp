class Solution {
public:
    string mergeAlternately(string word1, string word2) {
        std::stringstream ans;
        int l1 = word1.size(), l2 = word2.size();
        int lMax = max(l1, l2);
        for (int i = 0; i < lMax; ++i) {
            if (i < l1) {
                ans << word1[i];
            } else {
                ans << word2.substr(i);
                break;
            }
            if (i < l2) {
                ans << word2[i];
            } else {
                ans << word1.substr(i + 1);
                break;
            }
        }
        return ans.str();
    }
};