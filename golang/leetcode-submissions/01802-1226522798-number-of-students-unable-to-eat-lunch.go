func countStudents(students []int, sandwiches []int) int {
    req := make([]int, 2)
    for _, student := range students {
        req[student]++
    }

    for _, sandwich := range sandwiches {
        if req[sandwich] == 0 {
            return req[1 - sandwich]
        } else {
            req[sandwich]--
        }
    }

    return 0
}