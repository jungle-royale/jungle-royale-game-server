package physical

type RotRectangle struct {
	rect     Rectangle
	rotation float64 // degree
}

func NewRotRectangle(
	x float32,
	y float32,
	width float32,
	length float32,
	rotation float64,
) *RotRectangle {
	return &RotRectangle{
		*NewRectangle(x, y, width, length),
		rotation,
	}
}
