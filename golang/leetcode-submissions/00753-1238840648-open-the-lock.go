import "strconv"

func openLock(deadends []string, target string) int {
    pow10 := []int{1, 10, 100, 1000}
    visit := make([]int, 10000) // 0: not visited, 1: visited through forward direction, -1: visited through backward direction, 2: deadends
    for _, dead := range deadends {
        num, _ := strconv.Atoi(dead)
        visit[num] = 2
    }
    src := 0
    dest, _ := strconv.Atoi(target)
    steps := 0
    dir := 1
    if visit[src] == 2 || visit[dest] == 2 {
        return -1
    }
    if src == dest {
        return 0
    }
    forward := make([]int, 0)
    backward := make([]int, 0)
    forward = append(forward, src)
    visit[src] = 1
    backward = append(backward, dest)
    visit[dest] = -1
    for len(forward) > 0 && len(backward) > 0 {
        if len(forward) > len(backward) {
            forward, backward = backward, forward
            dir = -dir
        }
        steps++
        size := len(forward)
        for j := 0; j < size; j++ {
            cur := forward[0]
            forward = forward[1:]
            for _, p := range pow10 {
                d := (cur / p) % 10
                for _, i := range []int{-1, 1} {
                    z := d + i
                    if z == -1 {
                        z = 9
                    } else if z == 10 {
                        z = 0
                    }
                    next := cur + (z-d)*p
                    if visit[next] == -dir {
                        return steps
                    }
                    if visit[next] == 0 {
                        forward = append(forward, next)
                        visit[next] = dir
                    }
                }
            }
        }
    }
    return -1
}
