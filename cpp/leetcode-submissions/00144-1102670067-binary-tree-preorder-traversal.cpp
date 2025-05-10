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
        std::vector<TreeNode*> q{root};
        while (q.size() != 0) {
            auto node = q.back();
            auto left{node->left}, right{node->right};
            q.pop_back();
            ans.push_back(node->val);
            if (right != nullptr) {
                q.push_back(right);
            }
            if (left != nullptr) {
                q.push_back(left);
            }
        }
        return ans;
    }
};