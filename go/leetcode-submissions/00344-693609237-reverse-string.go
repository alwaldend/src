func reverseString(characters []byte) {
	length := len(characters)
	for index := 0; index < length/2; index++ {
		index_last := length - index -1 
		current, last := characters[index], characters[index_last]
		characters[index_last] = current
		characters[index] = last
	}
}
