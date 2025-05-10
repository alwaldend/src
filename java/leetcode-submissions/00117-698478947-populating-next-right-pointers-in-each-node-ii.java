class Solution {
    public Node connect(Node root) {
        Node leftMost = root;
        while (leftMost != null) {
            Node cur = leftMost;
            leftMost = null;
            Node pre = null;
            while (cur != null) {
                if (leftMost == null) {
                    leftMost = cur.left == null ? cur.right: cur.left;
                }
                if (pre != null) {
                    pre.next = cur.left == null ? cur.right : cur.left;
                }
                if (cur.left != null && cur.right != null) {
                    cur.left.next = cur.right;
                }
                pre = cur.right == null ? (cur.left == null ? pre : cur.left) : cur.right;
                cur = cur.next;
            }
        }
        return root;
    }
}