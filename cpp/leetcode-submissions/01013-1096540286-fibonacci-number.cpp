class Solution {
public:
    int fib(int n) {
        if (n < 2) {
            return n;
        }
        int cur = 1, prev = 0;
        for (int i = 2; i <= n; ++i) {
            int newVal = cur + prev;
            prev = cur;
            cur = newVal;
        }
        return cur;
    }
};