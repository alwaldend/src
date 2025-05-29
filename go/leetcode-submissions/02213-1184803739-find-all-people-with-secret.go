package main

import (
	"sort"
)

func find(groups []int, index int) int {
	for index != groups[index] {
		index = groups[index]
	}
	return index
}

func findAllPeople(n int, meetings [][]int, firstPerson int) []int {
	groups := make([]int, 100000)
	var result []int
	var temp []int

	for i := 0; i < n; i++ {
		groups[i] = i
	}
	groups[firstPerson] = 0

	sort.Slice(meetings, func(i, j int) bool {
		return meetings[i][2] < meetings[j][2]
	})

	i := 0
	for i < len(meetings) {
		currentTime := meetings[i][2]
		temp = temp[:0]
		for i < len(meetings) && meetings[i][2] == currentTime {
			g1 := find(groups, meetings[i][0])
			g2 := find(groups, meetings[i][1])
			groups[max(g1, g2)] = min(g1, g2)
			temp = append(temp, meetings[i][0], meetings[i][1])
			i++
		}
		for _, j := range temp {
			if find(groups, j) != 0 {
				groups[j] = j
			}
		}
	}

	for j := 0; j < n; j++ {
		if find(groups, j) == 0 {
			result = append(result, j)
		}
	}

	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
