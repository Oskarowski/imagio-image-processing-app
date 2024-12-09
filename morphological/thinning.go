package morphological

var StructuralElements = []BinaryImage{
	{
		{1, 1, 1},
		{-1, 1, -1},
		{0, 0, 0},
	},
	{
		{-1, 0, 0},
		{1, 1, 0},
		{1, 1, -1},
	},
	{
		{1, -1, 0},
		{1, 1, 0},
		{1, -1, 0},
	},
	{
		{1, 1, -1},
		{1, 1, 0},
		{-1, 0, 0},
	},
	{
		{1, 1, 1},
		{-1, 1, -1},
		{0, 0, 0},
	},
	{
		{-1, 1, 1},
		{0, 1, 1},
		{0, 0, -1},
	},
	{
		{0, -1, 1},
		{0, 1, 1},
		{0, -1, 1},
	},
	{
		{0, 0, -1},
		{0, 1, 1},
		{-1, 0, 1},
	},
}

func matchesStructuralElement(img BinaryImage, x, y int, se BinaryImage) bool {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			seVal := se[dy+1][dx+1]
			if seVal == -1 {
				continue // Don't care condition
			}
			if img[y+dy][x+dx] != seVal {
				return false
			}
		}
	}
	return true
}

// Reference: https://www.ee.nthu.edu.tw/clhuang/09420EE368000DIP/chapter09.pdf
// Reference: https://homepages.inf.ed.ac.uk/rbf/HIPR2/thin.htm
func Thinning(image BinaryImage, structElems []BinaryImage) BinaryImage {
	height := len(image)
	width := len(image[0])

	// Helper function to apply a thinning operation
	applyThinning := func(img BinaryImage, se BinaryImage) (BinaryImage, bool) {
		changed := false
		result := make(BinaryImage, height)
		for i := range result {
			result[i] = make([]int, width)
			copy(result[i], img[i])
		}

		for y := 1; y < height-1; y++ {
			for x := 1; x < width-1; x++ {
				if img[y][x] == 1 && matchesStructuralElement(img, x, y, se) {
					result[y][x] = 0
					changed = true
				}
			}
		}

		return result, changed
	}

	changed := true
	for changed {
		changed = false
		for _, se := range structElems {
			var seChanged bool
			image, seChanged = applyThinning(image, se)
			changed = changed || seChanged
		}
	}

	return image
}
