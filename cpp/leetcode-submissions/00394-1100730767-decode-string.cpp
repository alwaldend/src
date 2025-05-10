class Solution {
public:
    string decodeString(string s) {
        int pos = 0;
        return decode(pos, s);
    }
    
    string decode(int& pos, string s) {
        int num {0};
        string word {""};
        unsigned long length {s.size()};
        for(; pos < length; ++pos) {
            char c = s[pos];
            if(c == '[') {
                string repeat {decode(++pos, s)};
                for(; num > 0; --num) {
                    word += repeat;
                }
            } else if (c >= '0' && c <='9') {
                num = num * 10 + (c - '0');
            } else if (c == ']') {
                return word;
            } else {
                word += c;
            }
        }
        return word;
    }
};