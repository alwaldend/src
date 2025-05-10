class Solution {
public:
    int largestAltitude(vector<int>& gain) {
        int att {0}, maxAtt {0};
        for (int num : gain) {
            att += num;
            maxAtt = max(maxAtt, att);
        }
        return maxAtt;
    }
};