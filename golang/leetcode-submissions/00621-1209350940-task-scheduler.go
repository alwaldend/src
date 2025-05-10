func leastInterval(tasks []byte, n int) int {
	if n == 0 {
		return len(tasks)
	}

	cnt := make([]int, 26)
	for _, task := range tasks {
		cnt[task - 'A']++
	}

	var maxCount, sameMaxCount int
	for _, c := range cnt {
		if c > maxCount {
			maxCount = c
			sameMaxCount = 1
		} else if c == maxCount {
			sameMaxCount++
		}
	}

	res := (n + 1) * (maxCount - 1) + sameMaxCount
    if (res > len(tasks)) {
        return res
    } else {
        return len(tasks)
    }
}