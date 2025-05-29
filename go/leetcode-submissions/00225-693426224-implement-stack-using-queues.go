
type MyStack struct{ queue *list.List }

func Constructor() MyStack           { return MyStack{&list.List{}} }
func (this *MyStack) Push(value int) { this.queue.PushFront(value) }
func (this *MyStack) Top() int       { return this.queue.Front().Value.(int) }
func (this *MyStack) Empty() bool    { return this.queue.Len() == 0 }
func (this *MyStack) Pop() int {
	return this.queue.Remove(this.queue.Front()).(int)
}