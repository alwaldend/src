type SeatManager struct {
    seats []bool
}


func Constructor(n int) SeatManager {
    return SeatManager{make([]bool, n)}
}


func (this *SeatManager) Reserve() int {
    for i, num := range this.seats {
        if !num {
            this.seats[i] = true
            return i + 1
        }
    }
    return -1
}


func (this *SeatManager) Unreserve(seatNumber int)  {
    this.seats[seatNumber-1] = false
}


/**
 * Your SeatManager object will be instantiated and called as such:
 * obj := Constructor(n);
 * param_1 := obj.Reserve();
 * obj.Unreserve(seatNumber);
 */