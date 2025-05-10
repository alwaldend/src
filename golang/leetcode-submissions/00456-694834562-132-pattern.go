
func find132pattern(numbers []int) bool {
	length := len(numbers)
	if length < 3 {
		// if the array doesn't have at least three numbers, it cannot have
		// '123' pattern
		return false
	}
	list, third_element := list.List{}, math.MinInt
	for index := length - 1; index >= 0; index-- {
		current := numbers[index]
		if current < third_element {
			return true
		}
		for list.Len() != 0 && list.Front().Value.(int) < current {
			third_element = list.Front().Value.(int)
			list.Remove(list.Front())
		}
		list.PushFront(current)
	}
	return false
}