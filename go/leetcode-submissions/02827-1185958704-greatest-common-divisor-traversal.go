func canTraverseAllPairs(nums []int) bool {
    if len(nums) == 1 {
        return true
    }
    n := len(nums)
    maxElement := nums[0]
    minElement := nums[0]
    for _, num := range nums {
        if num > maxElement {
            maxElement = num
        }
        if num < minElement {
            minElement = num
        }
    }
    if minElement == 1 {
        return false
    }
    factorArray := factorsCalculator(maxElement)

    parent := make([]int, maxElement+1)
    rank := make([]int, maxElement+1)
    for i := 0; i <= maxElement; i++ {
        parent[i] = i
        rank[i] = 1
    }

    for _, num := range nums {
        x := num
        for x > 1 {
            p := factorArray[x]
            union(parent, rank, p, num)
            for x%p == 0 {
                x = x / p
            }
        }
    }

    p := find(parent, nums[0])
    for i := 1; i < n; i++ {
        if find(parent, nums[i]) != p {
            return false
        }
    }

    return true
}

func factorsCalculator(n int) []int {
    dp := make([]int, n+2)
    for i := 0; i < len(dp); i++ {
        dp[i] = i
    }
    for i := 2; i <= n; i++ {
        if dp[i] == i {
            for j := i * 2; j <= n; j += i {
                if dp[j] == j {
                    dp[j] = i
                }
            }
        }
    }
    return dp
}

func find(parent []int, a int) int {
    if parent[a] == a {
        return a
    }
    parent[a] = find(parent, parent[a])
    return parent[a]
}

func union(parent []int, rank []int, a int, b int) {
    a = find(parent, a)
    b = find(parent, b)
    if a == b {
        return
    }
    if rank[a] < rank[b] {
        a, b = b, a
    }
    parent[b] = a
    rank[a] += rank[b]
}