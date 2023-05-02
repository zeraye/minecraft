package minecraft

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type meshRenderer struct {
	md *meshDrawer
}

func (mr *meshRenderer) Destroy() {
}

func (mr *meshRenderer) Layout(s fyne.Size) {
}

func (mr *meshRenderer) MinSize() fyne.Size {
	return fyne.NewSize(float32(mr.md.eng.width), float32(mr.md.eng.height))
}

func (mr *meshRenderer) Objects() []fyne.CanvasObject {
	return mr.md.lines
}

func (mr *meshRenderer) Refresh() {
	mr.md.lines = linesFromMesh(mr.md.meshObject, mr.md.eng)
	canvas.Refresh(mr.md)
}
