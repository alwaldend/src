func numRescueBoats(people []int, limit int) int {
    sort.Ints(people)

    numberOfBouts := 0

    start := 0
    end := len(people)-1

    for start < end {
        if people[start] + people[end] <= limit {
            numberOfBouts++
            start++
        }else{
            numberOfBouts++
        }
        end--
    }

    if start == end {
        numberOfBouts++
    }

    return numberOfBouts
}