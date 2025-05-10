class Solution {
public:
    int numberOfMatches(int n) {
        int ans = 0;
        while (n > 1) {
            if (n % 2 == 0) {
                int matches {n / 2};
                ans += matches;
                n = matches;
            } else {
                int matches {(n - 1) / 2};
                ans += matches;
                n = matches + 1;
            }
        }
        
        return ans;
    }
};