class Solution {
public:
    string predictPartyVictory(string senate) {
        std::queue<int> rad, dir;
        unsigned long length {senate.size()};
        for (int i = 0; i < length; ++i){
            if (senate[i] == 'R'){
                rad.push(i);
            } else {
                dir.push(i);
            }
        }
        while (!rad.empty() && !dir.empty()) {
            if (rad.front() < dir.front()) {
                rad.push(length);
            } else {
                dir.push(length);
            }
            ++length;
            rad.pop();
            dir.pop();
        }
        return rad.empty() ? "Dire" : "Radiant";
    }
};