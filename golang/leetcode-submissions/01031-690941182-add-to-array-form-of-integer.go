func addToArrayForm(number1 []int, add int) []int {
	if add == 0 {
		return number1
	}
	if len(number1) == 0 {
		return convert(add)
	}
	number2 := convert(add)
	length1, length2, carry := len(number1), len(number2), 0
	hightest := length1
	if length2 > length1 {
		hightest = length2
	}
	result := make([]int, hightest)
	index1, index2, indexResult := length1-1, 0, hightest-1
	for {
		index1Valid, index2Valid := index1 >= 0, index2 < length2
		if !index1Valid && !index2Valid && carry==0 {
			break
		}
		digit1, digit2 := 0, 0
		if index1Valid {
			digit1 = number1[index1]
			index1--
		}
		if index2Valid {
			digit2 = number2[index2]
			index2++
		}
		digitResult := digit1 + digit2 + carry
		if digitResult > 9 {
			carry = 1
			digitResult -= 10
		} else {
			carry = 0
		}
		if indexResult == -1 {
			result = append(result, 0)
			copy(result[1:], result[0:hightest])
            indexResult=0
		}
		result[indexResult] = digitResult
		indexResult--
	}
	return result
}

func convert(number int) (result []int) {
	for {
		if number == 0 {
			return
		}
		result = append(result, number%10)
		number /= 10
	}
}
