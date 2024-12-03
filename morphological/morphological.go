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
	rows := len(se.Data)
	cols := len(se.Data[0])
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			imgX := x - se.OriginX + i
			imgY := y - se.OriginY + j
			if imgX < 0 || imgX >= len(image) || imgY < 0 || imgY >= len(image[0]) {
				continue
			}
			if se.Data[i][j] == 1 && image[imgX][imgY] == 0 {
				return false
			}
		}
	}
	return true
}
