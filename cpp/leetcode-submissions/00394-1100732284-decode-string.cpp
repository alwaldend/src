class Solution {
public:
    string decodeString(string s) {
        return std::get<1>(decode(0, s));
    }
    
    std::tuple<int, string> decode(int pos, string s) {
        int num {0};
        string word {""};
        unsigned long length {s.size()};
        for(; pos < length; ++pos) {
            char c = s[pos];
            if(c == '[') {
                auto [newPos, repeat] {decode(++pos, s)};
                for(; num > 0; --num) {
                    word += repeat;
                }
                pos = newPos;
            } else if (c >= '0' && c <='9') {
                num = num * 10 + (c - '0');
            } else if (c == ']') {
                return {pos, word};
            } else {
                word += c;
            }
        }
        return {pos, word};
    }
};