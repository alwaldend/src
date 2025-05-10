/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

type Codec struct {
    
}

func Constructor() Codec {
    return Codec{}
}

// Serializes a tree to a single string.
func (this *Codec) serialize(root *TreeNode) string {
    if root == nil {
        return ""
    }
    var buffer bytes.Buffer
    queue := []*TreeNode{root}
    for length := len(queue); length > 0; length = len(queue) {
        for i := 0; i < length; i++ {
            node := queue[i]
            if node != nil {
                buffer.WriteString(strconv.Itoa(node.Val))
                queue = append(queue, node.Left, node.Right)
            }
            buffer.WriteString(",")
        }
        queue = queue[length:]
    }
    ans := buffer.String()
    return ans[:len(ans)-1]
}

// Deserializes your encoded data to tree.
func (this *Codec) deserialize(data string) *TreeNode {    
    nodes := []*TreeNode{}
    for _, str := range strings.Split(data, ",") {
        if str == "" {
            nodes = append(nodes, nil)
        } else {
            num, _ := strconv.Atoi(str)
            nodes = append(nodes, &TreeNode{num, nil, nil})
        }
    }
    start := 0
    length := len(nodes)
    for _, node := range nodes {
        if node == nil {
            continue
        }
        left, right := 2 * start + 1, 2 * start + 2
        if left < length {
            node.Left = nodes[left]
        }
        if right < length {
            node.Right = nodes[right]
        }
        start++
    }
    return nodes[0]
}


/**
 * Your Codec object will be instantiated and called as such:
 * ser := Constructor();
 * deser := Constructor();
 * data := ser.serialize(root);
 * ans := deser.deserialize(data);
 */