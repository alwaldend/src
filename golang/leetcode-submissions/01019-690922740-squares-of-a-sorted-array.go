
func sortedSquares(numbers []int) []int {
	length := len(numbers)
	if length == 0 || length == 1 {
		return square(numbers, false)
	}
	negativesIndex := -1
	for index, number := range numbers {
		if number >= 0 {
			negativesIndex = index
			break
		}
	}
	if negativesIndex == 0 || negativesIndex == -1 {
		return square(numbers, negativesIndex == -1)
	}
	//fmt.Println("negativesIndex", negativesIndex)
	result, resultIndex, negativesIndex := make([]int, length), 0, negativesIndex-1
	for positivesIndex := negativesIndex + 1; resultIndex < length; positivesIndex++ {
		positiveOverflow := positivesIndex >= length
		for {
			if negativesIndex < 0 || resultIndex >= length {
				break
			}
			if !positiveOverflow && numbers[negativesIndex]*-1 > numbers[positivesIndex] {
				break
			}
			
			result[resultIndex] = numbers[negativesIndex]
			resultIndex++
			//fmt.Println("negative", negativesIndex, resultIndex, numbers[negativesIndex], result, numbers)
            negativesIndex--
		}
		if resultIndex < length && !positiveOverflow {
	//		fmt.Println("positivesIndex", positivesIndex, resultIndex, numbers[positivesIndex], result, numbers)
			result[resultIndex] = numbers[positivesIndex]
			resultIndex++
		}
	}
	return square(result, false)
}

func square(array []int, reverse bool) []int {
	if reverse {
		length := len(array)
		reversed := make([]int, length)
		for index := length - 1; index >= 0; index-- {
			reversed[length-index-1] = array[index] * array[index]
		}
		return reversed
	} else {
		for index, number := range array {
			array[index] = number * number
		}
		return array
	}
}
