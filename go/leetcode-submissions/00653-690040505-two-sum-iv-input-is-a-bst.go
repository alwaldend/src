import "sort"

func construct(root *TreeNode, results map[int]int) {
	if root == nil {
		return
	}
	if results[root.Val] < 2 {
		results[root.Val] += 1
	}
	construct(root.Left, results)
	construct(root.Right, results)
}

func sorted(input_map map[int]int) (results []int) {
	for key := range input_map {
		results = append(results, key)
	}
	return
}

func findTarget(root *TreeNode, target int) bool {
	if root == nil {
		return false
	}
	results := map[int]int{}
	construct(root, results)
	if target%2 == 0 && results[target/2] == 2 {
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
