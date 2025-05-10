class Solution {
public:
    string reverseWords(string s) {
        string ans = "";
        string temp = "";
        int length = s.length();
        int j = 0;
        for (j = 0; s[j] == ' ' && j < length; ++j) { }
        for (int i = length - 1; i >= j; --i) {
            if(s[i] == ' '){
                if(temp != ""){                  
                    ans = ans + temp + ' ';
                    temp = "";
                }
                continue;
            }
            else{
                temp = s[i] + temp;
            }
        }
        return ans + temp;
    }
};