class Solution {
public:
    bool canPlaceFlowers(vector<int>& flowerbed, int n) {
        int curZeros = 0, length = flowerbed.size();
        if (flowerbed[0] == 0) {
            curZeros = 2;
        }
        int ans = 0;
        for (int i = 1; i < length; ++i) {
            bool isFlower = flowerbed[i] == 1; 
            if (!isFlower) {
                curZeros += 1;
                continue;
            }
            ans += (curZeros - 1) / 2;
            curZeros = 0;
            if (ans >= n) {
                return true;
            }
        }
        ans += curZeros / 2;
        return ans >= n;
    }
};