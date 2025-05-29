func compareVersion(version1 string, version2 string) int {

	rev1 := strings.Split(version1, ".")
	rev2 := strings.Split(version2, ".")
    // fmt.Println(rev1, rev2)

	revInt1 := make([]int, len(rev1))
	revInt2 := make([]int, len(rev2))

	for i := 0; i < len(rev1); i++ {
		revInt1[i], _ = strconv.Atoi(rev1[i])
	}
	for i := 0; i < len(rev2); i++ {
		revInt2[i], _ = strconv.Atoi(rev2[i])
	}
    // fmt.Println(revInt1, revInt2)

	i := 0
	j := 0

	for i < len(revInt1) && j < len(revInt2) {

		if revInt1[i] < revInt2[j] {
			return -1
		} else if revInt1[i] > revInt2[j] {
			return 1
		}
		i++
		j++
	}

	for i < len(revInt1) {
		if revInt1[i] > 0 {
			return 1
		}
		i++
	}
	for j < len(revInt2) {
		if revInt2[j] > 0 {
			return -1
		}
		j++
	}
	return 0
}
