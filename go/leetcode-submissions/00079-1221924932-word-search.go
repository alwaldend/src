func exist(board [][]byte, word string) bool {
    if nLetters := len(board) * len(board[0]); nLetters < len(word) {
        return false
    }

    for i := range board {
        for j := range board[i] {
            if dfs(i, j, 0, board, word, make(map[pair]struct{})) {
                return true
            }
        }
    }

    return false
}

type pair struct {
    r, c int
}

func dfs(r, c, i int, board [][]byte, word string, visited map[pair]struct{}) bool {
    if i == len(word) {
        return true
    }

    inBounds := r >= 0 && r < len(board) && c >= 0 && c < len(board[0])
    if _, ok := visited[pair{r, c}]; !inBounds || ok || word[i] != board[r][c] {
        return false
    }

    visited[pair{r, c}] = struct{}{}

    up := dfs(r+1, c, i+1, board, word, visited)
    down := dfs(r-1, c, i+1, board, word, visited)
    right := dfs(r, c+1, i+1, board, word, visited)
    left := dfs(r, c-1, i+1, board, word, visited)

    delete(visited, pair{r, c})

    return up || down || right || left
}
