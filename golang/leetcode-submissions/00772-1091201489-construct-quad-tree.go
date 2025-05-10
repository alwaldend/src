/**
 * Definition for a QuadTree node.
 * type Node struct {
 *     Val bool
 *     IsLeaf bool
 *     TopLeft *Node
 *     TopRight *Node
 *     BottomLeft *Node
 *     BottomRight *Node
 * }
 */

func construct(grid [][]int) *Node {
	var dfs func(y0, x0, width int) *Node
	dfs = func(y0, x0, width int) *Node {
		if width == 1 {
			return &Node{
				Val:    grid[y0][x0] == 1,
				IsLeaf: true,
			}
		}

		w := width / 2
		topLeft := dfs(y0, x0, w)
		topRight := dfs(y0, x0+w, w)
		bottomLeft := dfs(y0+w, x0, w)
		bottomRight := dfs(y0+w, x0+w, w)
		var node *Node

		if topLeft.Val == topRight.Val && bottomLeft.Val == bottomRight.Val && topLeft.Val == bottomLeft.Val &&
			topLeft.IsLeaf && topRight.IsLeaf && bottomLeft.IsLeaf && bottomRight.IsLeaf {
			node = &Node{
				Val:    topLeft.Val,
				IsLeaf: true,
			}
		} else {
			node = &Node{
				Val:         true,
				IsLeaf:      false,
				TopLeft:     topLeft,
				TopRight:    topRight,
				BottomLeft:  bottomLeft,
				BottomRight: bottomRight,
			}
		}
		return node
	}
	return dfs(0, 0, len(grid))
}
