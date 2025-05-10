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
        if (root == nullptr) {
            return {};
        }
        std::vector<int> ans{};
        std::deque<TreeNode*> q{root};
        while (q.size() != 0) {
            auto node = q.front();
            auto left{node->left}, right{node->right};
            q.pop_front();
            ans.push_back(node->val);
            if (right != nullptr) {
                q.push_front(right);
            }
            if (left != nullptr) {
                q.push_front(left);
            }
        }
        return ans;
    }
};