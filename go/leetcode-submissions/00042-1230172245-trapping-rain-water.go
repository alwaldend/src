func trap(height []int) int {
	var waterStored int
	leftElevationPoint, rightElevationPoint := 0, len(height)-1                                    // pointers for left and right elevations
	maxLeftElevation, maxRightElevation := height[leftElevationPoint], height[rightElevationPoint] // max value of left and right elevation points

	for leftElevationPoint < rightElevationPoint { // until left elevation point reaches right elevation point or right elevation point reaches left elevation point
		if maxLeftElevation < maxRightElevation { // if max left elevation is smaller than max right elevation
			leftElevationPoint++                               // move the left elevation point by 1
			if height[leftElevationPoint] > maxLeftElevation { // if height of current left elevation after moving is greater than max left elevation
				maxLeftElevation = height[leftElevationPoint] // swap the max left elevation with current left elevation point value
			}
			waterStored += maxLeftElevation - height[leftElevationPoint] // store water with delta of max left elevation and current left elevation
		} else {
			rightElevationPoint--                                // decrease the right elevation by 1
			if height[rightElevationPoint] > maxRightElevation { // if height of current right elevation after decreasing is greater than max right elevation
				maxRightElevation = height[rightElevationPoint] // swap the max right elevation with current right elevation point value
			}
			waterStored += maxRightElevation - height[rightElevationPoint] // store water with delta of max right elevation and current right elevation
		}
	}

	return waterStored
}