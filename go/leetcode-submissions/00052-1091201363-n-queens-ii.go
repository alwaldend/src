func totalNQueens(n int) int {
    sCol:=make([]bool,n)
    sD1:=make([]bool,2*n)
    sD2:=make([]bool,2*n)
    return helper(0,n,sCol,sD1,sD2)
}
func helper(r,n int,sCol,sD1,sD2 []bool) int{
    if r==n{
        return 1
    }
    res:=0
    for i:=0; i < n; i++ {
        if !sCol[i] && !sD1[i+r] && !sD2[(r-i)+n]{
          // board[r][i]=true
            sCol[i]=true
            sD1[i+r]=true
            sD2[(r-i)+n]=true
            res=res+helper(r+1,n,sCol,sD1,sD2)
          //  board[r][i]=false
            sCol[i]=false
            sD1[i+r]=false
            sD2[(r-i)+n]=false
        }
    }
    return res
}