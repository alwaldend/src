type MyHashSet struct {
    m map[int]struct{}
}


func Constructor() MyHashSet {
    return MyHashSet{map[int]struct{}{}}
}


func (this *MyHashSet) Add(key int)  {
    this.m[key] = struct{}{}
}


func (this *MyHashSet) Remove(key int)  {
    delete(this.m, key)
}


func (this *MyHashSet) Contains(key int) bool {
    _, ok := this.m[key]
    return ok
}


/**
 * Your MyHashSet object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Add(key);
 * obj.Remove(key);
 * param_3 := obj.Contains(key);
 */