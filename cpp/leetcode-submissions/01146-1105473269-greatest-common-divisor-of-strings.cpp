class Solution {
public:
    bool valid(string str1, string str2, size_t k) {
        size_t len1 {str1.size()}, len2 {str2.size()};
        if (len1 % k != 0 || len2 % k != 0) {
            return false;
        }
        string base = str1.substr(0, k);
        size_t n1 {len1 / k}, n2 {len2 / k};
        if (n1 == n2) {
            return str1 == str2 && joinWords(base, n1) == str1;
        }
        return str1 == joinWords(base, n1) && str2 == joinWords(base, n2);
    }

    string joinWords(string str, size_t k) {
        string ans = "";
        for (size_t i = 0; i < k; ++i) {
            ans += str;
        }
        return ans;
    }
    
    
    string gcdOfStrings(string str1, string str2) {
        for (size_t i = min(str1.size(), str2.size()); i > 0; --i) {
            if (valid(str1, str2, i)) {
                return str1.substr(0, i);
            }
        }
        return "";
    }
};