class Solution {
public:
    int compress(vector<char>& chars) {
        int i = 0, res = 0;
        int length = chars.size();
        while (i < length) {
            int groupLength = 1;
            while (i + groupLength < length && chars[i + groupLength] == chars[i]) {
                groupLength++;
            }
            chars[res++] = chars[i];
            if (groupLength > 1) {
                for (char c : to_string(groupLength)) {
                    chars[res++] = c;
                }
            }
            i += groupLength;
        }
        return res;
    }
};