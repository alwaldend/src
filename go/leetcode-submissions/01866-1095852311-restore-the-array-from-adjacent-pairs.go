func restoreArray(adjacentPairs [][]int) []int {
    // [[2,1],[3,4],[3,2]]
    // 1: [2]
    // 2: [1, 3]
    // 3: [2, 4]
    // 4: [3]
    // Output: [1,2,3,4]
    graph := map[int][]int{}
    for _, pair := range adjacentPairs {
        num1, num2 := pair[0], pair[1]
        graph[num1] = append(graph[num1], num2)
        graph[num2] = append(graph[num2], num1)
    }
    length := len(adjacentPairs) + 1
    ans := make([]int, length)
    for node, edges := range graph {
        if len(edges) == 1 {
            ans[0], ans[1] = node, edges[0]
            break
        }
    }
    for i := 2; i < length; i++ {
        cur, prev := ans[i-1], ans[i-2]
        for _, target := range graph[cur] {
            if target != prev {
                ans[i] = target
            }
        }
    }
    return ans
}