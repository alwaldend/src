import "sort"

func construct(root *TreeNode, results map[int]bool, doubles map[int]bool) {
	if root == nil {
		return
	}
	_, exists := results[root.Val]
	if !exists {
		results[root.Val] = true
	} else if _, double := doubles[root.Val]; exists && double {
		doubles[root.Val] = true
	}
	construct(root.Left, results, doubles)
	construct(root.Right, results, doubles)
}

func sorted(input_map map[int]bool) (results []int) {
	for key := range input_map {
		results = append(results, key)
	}
	return
}

func findTarget(root *TreeNode, target int) bool {
	if root == nil {
		return false
	}
	results := map[int]bool{}
	doubles := map[int]bool{}
	construct(root, results, doubles)
	if _, double_exists := doubles[target/2]; target%2 == 0 && double_exists {
		return true
	}
	results_sorted := sorted(results)
	results_sorted_len := len(results_sorted)
	sort.Ints(results_sorted)
	for index, number := range results_sorted {
		for index_inner := index + 1; index_inner < results_sorted_len; index_inner++ {
			if number+results_sorted[index_inner] == target {
				return true
			}
		}
	}
	return false
}