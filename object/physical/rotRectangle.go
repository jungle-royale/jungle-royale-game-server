package physical

type RotRectangle struct {
	rect     Rectangle
	rotation float64 // degree
}

func NewRotRectangle(
	x float64,
	y float64,
	width float64,
	length float64,
	rotation float64,
) *RotRectangle {
	return &RotRectangle{
		*NewRectangle(x, y, width, length),
		rotation,
	}
}
