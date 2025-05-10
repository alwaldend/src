/**
 * Definition for a binary tree node.
 * public class TreeNode {
 *     int val;
 *     TreeNode left;
 *     TreeNode right;
 *     TreeNode() {}
 *     TreeNode(int val) { this.val = val; }
 *     TreeNode(int val, TreeNode left, TreeNode right) {
 *         this.val = val;
 *         this.left = left;
 *         this.right = right;
 *     }
 * }
 */
class Solution {
    public int deepestLeavesSum(TreeNode root) {
        MaxDepthInfo maxDepthInfo = new MaxDepthInfo(0 ,0);
        sumAtLevel(root, 0, maxDepthInfo);
        return maxDepthInfo.getSumAtMaxDepth();
    }

    public void sumAtLevel(TreeNode root, int currentDepth, MaxDepthInfo maxDepthInfo) {
        if (root == null) return;

        if (currentDepth > maxDepthInfo.getMaxDepth()) {
            maxDepthInfo.setMaxDepth(currentDepth);
            maxDepthInfo.setSumAtMaxDepth(root.val);
        }

        else if (currentDepth == maxDepthInfo.getMaxDepth())
            maxDepthInfo.setSumAtMaxDepth(maxDepthInfo.getSumAtMaxDepth() + root.val);

        sumAtLevel(root.left, currentDepth + 1, maxDepthInfo);
        sumAtLevel(root.right, currentDepth + 1, maxDepthInfo);
    }

    public static class MaxDepthInfo {
        private int maxDepth;
        private int sumAtMaxDepth;

        public MaxDepthInfo(int maxDepth, int sumAtMaxDepth) {
            this.maxDepth = maxDepth;
            this.sumAtMaxDepth = sumAtMaxDepth;
        }

        public int getMaxDepth() { return maxDepth;}

        public void setMaxDepth(int maxDepth) { this.maxDepth = maxDepth;}

        public int getSumAtMaxDepth() { return sumAtMaxDepth;}

        public void setSumAtMaxDepth(int sumAtMaxDepth) { this.sumAtMaxDepth = sumAtMaxDepth;}
    }
}