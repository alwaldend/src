func moveZeroes(nums []int)  {
    	if len(nums) < 2 {
		return
	}

	for z, p := 0, 1; p < len(nums) && z < len(nums); {
		if nums[z] == 0 && nums[p] != 0 {
			if p > z {
				nums[z], nums[p] = nums[p], nums[z]
				z++
			}
			p = z+1
		} else {
			if nums[z] != 0 {
				z++
			}
			if nums[p] == 0 {
				p++
			}
		}
	}

}