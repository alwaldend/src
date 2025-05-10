func minMutation(startGene string, endGene string, bank []string) int {
    if startGene == endGene {
        return 0
    }

    bank = append(bank, startGene)

    isMut := func (gene1 string, gene2 string) bool {
        foundDiff := false
        for i := 0; i < len(gene1); i++ {
            if gene1[i] == gene2[i] {
                continue
            }
            if foundDiff {
                return false
            }
            foundDiff = true
        } 
        return foundDiff
    }

    graph := map[string][]string{}

    for i, gene1 := range bank {
        for _, gene2 := range bank[i+1:] {
            if !isMut(gene1, gene2) {
                continue
            }
            if _, ok := graph[gene1]; !ok {
                graph[gene1] = []string{}
            }
            if _, ok := graph[gene2]; !ok {
                graph[gene2] = []string{}
            }
            graph[gene1] = append(graph[gene1], gene2)
            graph[gene2] = append(graph[gene2], gene1)
        }
    }

    queue, ok := graph[endGene]
    if !ok {
        return -1
    }
    visited := map[string]struct{}{endGene: struct{}{}}
    num := 1
    for length := len(queue); length > 0; length = len(queue) {
        for i := 0; i < length; i++ {
            gene := queue[i]
            if _, ok := visited[gene]; ok {
                continue
            }
            if gene == startGene {
                return num 
            }
            visited[gene] = struct{}{}
            queue = append(queue, graph[gene]...)
        }
        num += 1
        queue = queue[length:]
    }

    return -1
}