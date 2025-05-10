type Node struct {
	children [26]*Node
	word     string
}

func (n *Node) Insert(word string) {
	cur := n
	for _, c := range word {
		idx := c - 'a'
		if cur.children[idx] == nil {
			cur.children[idx] = &Node{}
		}
		cur = cur.children[idx]
	}
	cur.word = word
}

func (n *Node) IsEmpty() bool {
	for _, child := range n.children {
		if child != nil {
			return false
		}
	}
	return true
}

func (n *Node) Remove(word string) bool {
	if len(word) == 0 {
		n.word = ""
		return n.IsEmpty()
	}
	child := n.children[word[0]-'a']
	if child.Remove(word[1:]) {
		n.children[word[0]-'a'] = nil
		return n.IsEmpty()
	}
	return false
}

func dfs(board [][]byte, r, c int, root, cur *Node, res *[]string) {
	rc := board[r][c]
	board[r][c] = 0
    
	if cur.word != "" {
		*res = append(*res, cur.word)
		root.Remove(cur.word)
	}
	ds := [5]int{0, 1, 0, -1, 0}
	for i := 0; i < len(ds)-1; i++ {
		dr, dc := r+ds[i], c+ds[i+1]
		if dr < 0 || dr >= len(board) || dc < 0 || dc >= len(board[0]) {
			continue
		}
		b := board[dr][dc]
		if b == 0 || cur.children[b-'a'] == nil {
			continue
		}
		dfs(board, dr, dc, root, cur.children[b-'a'], res)
	}
	board[r][c] = rc
}

func findWords(board [][]byte, words []string) []string {
	m, n := len(board), len(board[0])
	res, trie, has := []string{}, &Node{}, map[string]bool{}

	for r := 0; r < m; r++ {
		for c := 0; c < n-1; c++ {
			p := string(board[r][c]) + string(board[r][c+1])
			has[p] = true
		}
	}
	for r := 0; r < m-1; r++ {
		for c := 0; c < n; c++ {
			p := string(board[r][c]) + string(board[r+1][c])
			has[p] = true
		}
	}
	for _, word := range words {
		valid := true
		for i := 0; i < len(word)-1; i++ {
			a, b := string(word[i]), string(word[i+1])
			if !has[a+b] && !has[b+a] {
				valid = false
				break
			}
		}
		if valid {
			trie.Insert(word)
		}
	}
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			b := board[r][c]
			if trie.children[b-'a'] != nil {
				dfs(board, r, c, trie, trie.children[b-'a'], &res)
			}
		}
	}
	return res
}