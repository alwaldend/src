func reverseWords(_string string) string {
	index_word_start, length, result := 0, len(_string), []rune(_string)
	for index, character := range result {
		// ignore normal characters
		if character != ' ' {
			continue
		}
		// word ended -> reverse characters from the start of the word to the
		//end of it
		reverse_word(result, length, index_word_start, index)
		index_word_start = index + 1
	}
	reverse_word(result, length, index_word_start, length)
	return string(result)
}

func reverse_word(_string []rune, length int, start int, end int) {
	//fmt.Println("reversing", string(_string), start, "->", end)
	length_word := end - start
	for index := start; index < start+length_word/2; index++ {
		index_last := end - (index - start) - 1
		current, last := _string[index], _string[index_last]
		_string[index], _string[index_last] = last, current
	}
	//fmt.Println("result", string(_string), start, "->", end)
}
