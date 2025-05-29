func numSubarraysWithSum(nums []int, goal int) int {

    hashmap := make(map[int]int)

    hashmap[0] = 1

    sum := 0

    count := 0



    for _, num := range nums {

        sum += num

        rem := sum - goal

        if val, ok := hashmap[rem]; ok {

            count += val

        }

        hashmap[sum]++

    }



    return count

}