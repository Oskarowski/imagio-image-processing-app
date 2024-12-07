package morphological

type BinaryImage [][]int

type StructuringElement struct {
	Data             [][]int
	OriginX, OriginY int
}

func ReflectSE(se StructuringElement) StructuringElement {
	rows := len(se.Data)
	cols := len(se.Data[0])
	reflected := make([][]int, rows)
	for i := range reflected {
		reflected[i] = make([]int, cols)
		for j := range reflected[i] {
			reflected[i][j] = se.Data[rows-i-1][cols-j-1]
		}
	}
	return StructuringElement{Data: reflected, OriginX: se.OriginX, OriginY: se.OriginY}
}

func Fits(image BinaryImage, se StructuringElement, x, y int) bool {
	rows := len(image)
	cols := len(image[0])

	for i := 0; i < len(se.Data); i++ {
		for j := 0; j < len(se.Data[i]); j++ {
			if se.Data[i][j] == 1 {
				newX := x + i - se.OriginX
				newY := y + j - se.OriginY
				if newX < 0 || newX >= rows || newY < 0 || newY >= cols || image[newX][newY] == 0 {
					return false
				}
			}
		}
	}
	return true
}
