class Solution {
public:
    int maxCoins(vector<int>& piles) {
        std::sort(piles.begin(), piles.end());
        size_t length {piles.size()};
        size_t picks {length / 3};
        int count {};
        for (size_t i {length - 2}; picks > 0; i -= 2, --picks) {
            count += piles[i];
        }
        return count;
    }
};