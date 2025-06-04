class Solution:
    def imageSmoother(self, img: List[List[int]]) -> List[List[int]]:
        # Save the dimensions of the image.
        m = len(img)
        n = len(img[0])

        # Create a new image of the same dimension as the input image.
        smooth_img = [[0] * n for _ in range(m)]

        # Iterate over the cells of the image.
        for i in range(m):
            for j in range(n):
                # Initialize the sum and count 
                sum = 0
                count = 0

                # Iterate over all plausible nine indices.
                for x in (i - 1, i, i + 1):
                    for y in (j - 1, j, j + 1):
                        # If the indices form valid neighbor
                        if 0 <= x < m and 0 <= y < n:
                            sum += img[x][y]
                            count += 1

                # Store the rounded down value in smooth_img[i][j].
                smooth_img[i][j] = sum // count
        
        # Return the smooth image.
        return smooth_img