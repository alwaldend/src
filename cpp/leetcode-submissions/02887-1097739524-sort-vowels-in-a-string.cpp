class Solution {
public:
    string sortVowels(string s) {
        std::map<char, int> vowels{
            {'A', 0}, {'E', 0}, {'I', 0}, {'O', 0}, {'U', 0},
            {'a', 0}, {'e', 0}, {'i', 0}, {'o', 0}, {'u', 0}
        };
        unsigned long length{ s.size() };
        for (char ch : s) {
            if (vowels.contains(ch)) {
                ++vowels[ch];
            }
        }
        for (int i = 0; i < length; ++i) {
            const char ch = s[i];
            if (!vowels.contains(ch)) {
                continue;
            }
            for (const auto& [orderChar, count] : vowels) {
                if (count > 0) {
                    s[i] = orderChar;
                    --vowels[orderChar];
                    break;
                }
            }
        }
        return s;
    }
};