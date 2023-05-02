package main

import (
	"math"

	"fyne.io/fyne/app"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gonum.org/v1/gonum/mat"
)

var (
	_ fyne.Draggable = (*LineDrawer)(nil)
)

func main() {
	app := app.New()
	window := app.NewWindow("Minecraft")

	points := []*mat.VecDense{
		mat.NewVecDense(3, []float64{0, 0, 0}),
		mat.NewVecDense(3, []float64{1, 0, 0}),
		mat.NewVecDense(3, []float64{0, 1, 0}),
		mat.NewVecDense(3, []float64{1, 1, 0}),
		mat.NewVecDense(3, []float64{0, 0, 1}),
		mat.NewVecDense(3, []float64{1, 0, 1}),
		mat.NewVecDense(3, []float64{0, 1, 1}),
		mat.NewVecDense(3, []float64{1, 1, 1}),
	}

	matrix := mat.NewDense(len(points), 3, nil)

	for i := 0; i < len(points); i++ {
		points[i].ScaleVec(200, points[i])
		matrix.SetRow(i, points[i].RawVector().Data)
	}

	window.SetContent(NewLineDrawer(matrix))
	window.Resize(fyne.NewSize(800, 800))
	window.ShowAndRun()
}

type LineDrawer struct {
	widget.BaseWidget
	lines  []fyne.CanvasObject
	matrix *mat.Dense
}

func NewLineDrawer(matrix *mat.Dense) *LineDrawer {
	draw := &LineDrawer{
		lines:  LinesFromMatrix(matrix),
		matrix: matrix,
	}
	draw.ExtendBaseWidget(draw)
	return draw
}

func LinesFromMatrix(matrix *mat.Dense) []fyne.CanvasObject {
	return []fyne.CanvasObject{
		NewLineBetween(matrix.RawRowView(0), matrix.RawRowView(1)),
		NewLineBetween(matrix.RawRowView(1), matrix.RawRowView(3)),
		NewLineBetween(matrix.RawRowView(3), matrix.RawRowView(2)),
		NewLineBetween(matrix.RawRowView(2), matrix.RawRowView(0)),
		NewLineBetween(matrix.RawRowView(4), matrix.RawRowView(5)),
		NewLineBetween(matrix.RawRowView(5), matrix.RawRowView(7)),
		NewLineBetween(matrix.RawRowView(7), matrix.RawRowView(6)),
		NewLineBetween(matrix.RawRowView(6), matrix.RawRowView(4)),
		NewLineBetween(matrix.RawRowView(0), matrix.RawRowView(4)),
		NewLineBetween(matrix.RawRowView(1), matrix.RawRowView(5)),
		NewLineBetween(matrix.RawRowView(2), matrix.RawRowView(6)),
		NewLineBetween(matrix.RawRowView(3), matrix.RawRowView(7)),
	}
}

func (l *LineDrawer) Dragged(d *fyne.DragEvent) {
	dy := float64(d.Dragged.DY) / 100
	dx := -float64(d.Dragged.DX) / 100

	cosDx := math.Cos(dx)
	sinDx := math.Sin(dx)

	cosDy := math.Cos(dy)
	sinDy := math.Sin(dy)

	// Combined matrix for dragging in both directions.
	data := []float64{
		cosDx, 0, sinDx,
		sinDy * sinDx, cosDy, -sinDy * cosDx,
		-cosDy * sinDx, sinDy, cosDy * cosDx,
	}
	R := mat.NewDense(3, 3, data)
	l.matrix.Mul(l.matrix, R)

	l.lines = LinesFromMatrix(l.matrix)
	l.Refresh()
}

func (l *LineDrawer) DragEnd() {}

type lineRenderer struct {
	lineDrawer *LineDrawer
}

func (lr *lineRenderer) Destroy() {
}

func (lr *lineRenderer) Layout(s fyne.Size) {
}

func (lr *lineRenderer) MinSize() fyne.Size {
	return fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
}

func (lr *lineRenderer) Objects() []fyne.CanvasObject {
	return lr.lineDrawer.lines
}

func (lr *lineRenderer) Refresh() {
	canvas.Refresh(lr.lineDrawer)
}

func (l *LineDrawer) CreateRenderer() fyne.WidgetRenderer {
	return &lineRenderer{lineDrawer: l}
}

// func NewLineBetween(x1, y1, x2, y2 float64) *canvas.Line {
func NewLineBetween(vec1, vec2 []float64) *canvas.Line {
	return &canvas.Line{
		Position1:   fyne.NewPos(float32(vec1[0]+300), float32(vec1[1]+300)),
		Position2:   fyne.NewPos(float32(vec2[0]+300), float32(vec2[1]+300)),
		StrokeColor: theme.PrimaryColor(),
		StrokeWidth: 3,
	}
}
