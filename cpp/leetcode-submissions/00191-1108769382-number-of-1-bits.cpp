class Solution {
public:
    int hammingWeight(uint32_t n) {
        int ans {};
        while (n) {
            if (n & 1) {
                ++ans;
            }
            n >>= 1;
        }
        return ans;
    }
};