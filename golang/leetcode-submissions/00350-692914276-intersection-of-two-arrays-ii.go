
func intersect(numbers_1 []int, numbers_2 []int) []int {
	length_1, length_2 := len(numbers_1), len(numbers_2)
	length_biggest, result := length_1, make([]int, length_1+length_2)
	if length_2 > length_1 {
		length_biggest = length_2
	}
	occurences := make(map[int][]int, length_biggest)

	for index := 0; index < length_biggest; index++ {
		if index < length_1 {
			add_to_occurences(numbers_1[index], 0, occurences)
		}
		if index < length_2 {
			add_to_occurences(numbers_2[index], 1, occurences)
		}
	}
	index := 0
	for number, occurence := range occurences {
		repeat := occurence[0]
		if occurence[1] < repeat {
			repeat = occurence[1]
		}
		for ; repeat > 0; repeat-- {
			result[index] = number
			index++
		}
	}
    return result[0 : index]
}
func add_to_occurences(number int, index int, occurences map[int][]int) {
	if _, occured := occurences[number]; !occured {
		occurences[number] = make([]int, 2)
	}
	occurences[number][index]++
}
