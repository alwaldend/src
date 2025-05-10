/**
 * Definition for a binary tree node.
 * struct TreeNode {
 *     int val;
 *     TreeNode *left;
 *     TreeNode *right;
 *     TreeNode() : val(0), left(nullptr), right(nullptr) {}
 *     TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
 *     TreeNode(int x, TreeNode *left, TreeNode *right) : val(x), left(left), right(right) {}
 * };
 */
class Solution {
public:
    vector<int> preorderTraversal(TreeNode* root) {
        std::vector<int> ans{};
        std::deque<TreeNode*> q{root};
        while (q.size() != 0) {
            auto node = q.front();
            q.pop_front();
            if (node == nullptr) {
                continue;
            }
            ans.push_back(node->val);
            q.push_front(node->right);
            q.push_front(node->left);
        }
        return ans;
    }
};