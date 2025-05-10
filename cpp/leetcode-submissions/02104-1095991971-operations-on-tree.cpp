class LockingTree {
public:
    std::vector<int> parents, locked;
    std::vector<std::vector<int>> children;
    LockingTree(vector<int>& parent) {
        int length = parent.size();
        this->children.resize(length);
        this->locked.resize(length, -1);
        this->parents = parent;
        for (int i = 1; i < length; ++i) {
            this->children[this->parents[i]].push_back(i);
        }
    }
    
    bool lock(int num, int user) {
        if (this->locked[num] != -1) {
            return false;
        }
        this->locked[num] = user;
        return true;
    }
    
    bool unlock(int num, int user) {
        if (this->locked[num] != user) {
            return false;
        }
        this->locked[num] = -1;
        return true;
    }

    bool upgrade(int num, int user) {
        if (this->locked[num] != -1) {
            return false;
        }
        int parent = this->parents[num];
        while (parent != -1) {
            if (this->locked[parent] != -1) {
                return false;
            }
            parent = this->parents[parent];
        }
        if (!this->unlockDesc(num)) {
            return false;
        }
        this->locked[num] = user;
        return true;
    }

    bool unlockDesc(int parent) {
        bool hasLocked = false;
        for (auto const& child : this->children[parent]) {
            if (this->locked[child] != -1) {
                this->locked[child] = -1;
                hasLocked = true;
            }
            bool descHasLocked = unlockDesc(child);
            if (descHasLocked) {
                hasLocked = true;
            }
        }
        return hasLocked;
    }
};

/**
 * Your LockingTree object will be instantiated and called as such:
 * LockingTree* obj = new LockingTree(parent);
 * bool param_1 = obj->lock(num,user);
 * bool param_2 = obj->unlock(num,user);
 * bool param_3 = obj->upgrade(num,user);
 */