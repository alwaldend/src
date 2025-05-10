func deckRevealedIncreasing(deck []int) []int {
	n := len(deck)
	sort.Ints(deck)

	res := make([]int, n)
	queue := make([]int, n)

	for i := range queue {
		queue[i] = i
	}

	for _, card := range deck {
		idx := queue[0]
		queue = queue[1:]
		res[idx] = card

		if len(queue) > 0 {
			queue = append(queue, queue[0])
			queue = queue[1:]
		}
	}

	return res
}