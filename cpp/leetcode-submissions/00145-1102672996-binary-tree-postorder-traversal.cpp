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
    vector<int> postorderTraversal(TreeNode* root) {
        if (root == nullptr) {
            return {};
        }
        std::vector<int> ans{};
        addNodes(root, ans);
        return ans;
    }

    void addNodes(TreeNode* node, std::vector<int>& ans) {
        auto left{node->left}, right{node->right};
        if (left != nullptr) {
            addNodes(left, ans);
        }
        if (right != nullptr) {
            addNodes(right, ans);
        }
        ans.push_back(node->val);
    }
};