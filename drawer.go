package minecraft

import (
	"image/color"
	"math"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type meshDrawer struct {
	widget.BaseWidget
	lines      []fyne.CanvasObject
	meshObject *mesh
	eng        *engine
}

func (md *meshDrawer) CreateRenderer() fyne.WidgetRenderer {
	return &meshRenderer{md}
}

func newMeshDrawer(meshObject *mesh, eng *engine) *meshDrawer {
	drawer := &meshDrawer{
		lines:      linesFromMesh(meshObject, eng),
		meshObject: meshObject,
		eng:        eng,
	}
	drawer.ExtendBaseWidget(drawer)
	return drawer
}

func linesFromMesh(meshObject *mesh, eng *engine) []fyne.CanvasObject {
	lines := []fyne.CanvasObject{}
	frameTris := []*triangle{}
	for _, tri := range meshObject.tris {
		render, frameTri := linesFromTriangle(&tri, eng)
		if render {
			frameTris = append(frameTris, frameTri)
		}
	}

	// sort frameTris by z to draw by correct order
	sort.Slice(frameTris, func(i, j int) bool {
		z1 := (frameTris[i].p[0].z + frameTris[i].p[1].z + frameTris[i].p[2].z) / 3
		z2 := (frameTris[j].p[0].z + frameTris[j].p[1].z + frameTris[j].p[2].z) / 3
		return z1 < z2
	})

	for _, frameTri := range frameTris {
		// drawing frames
		// lines = append(lines, []fyne.CanvasObject{
		// 	newLineBetween(frameTri.p[0], frameTri.p[1], theme.PrimaryColor(), 3),
		// 	newLineBetween(frameTri.p[1], frameTri.p[2], theme.PrimaryColor(), 3),
		// 	newLineBetween(frameTri.p[2], frameTri.p[0], theme.PrimaryColor(), 3),
		// }...)
		lines = append(lines, fillLinesFromTriangle(frameTri, eng)...)
	}
	return lines
}

func linesFromTriangle(tri *triangle, eng *engine) (bool, *triangle) {
	matRotZ := createRotationZMatrix(eng.cam.fTheta)
	matRotX := createRotationXMatrix(eng.cam.fTheta)

	matProj := createProjectionMatrix(eng.cam.fAspectRatio, eng.cam.fFovRad, eng.cam.fFar, eng.cam.fNear)

	triProjected := triangle{}
	triTranslated := triangle{}
	triRotatedZ := triangle{}
	triRotatedZX := triangle{}

	// rotate in Z-axis
	triRotatedZ.p[0] = multiplyMatrixVector(tri.p[0], matRotZ)
	triRotatedZ.p[1] = multiplyMatrixVector(tri.p[1], matRotZ)
	triRotatedZ.p[2] = multiplyMatrixVector(tri.p[2], matRotZ)

	// rotate in X-axis
	triRotatedZX.p[0] = multiplyMatrixVector(triRotatedZ.p[0], matRotX)
	triRotatedZX.p[1] = multiplyMatrixVector(triRotatedZ.p[1], matRotX)
	triRotatedZX.p[2] = multiplyMatrixVector(triRotatedZ.p[2], matRotX)

	// offset into the screen
	triTranslated = triRotatedZX
	triTranslated.p[0].z = triRotatedZX.p[0].z + 12
	triTranslated.p[1].z = triRotatedZX.p[1].z + 12
	triTranslated.p[2].z = triRotatedZX.p[2].z + 12

	// use Cross-Product to get surface normal
	line1 := vec3d{}
	line2 := vec3d{}
	normal := vec3d{}

	line1.x = triTranslated.p[1].x - triTranslated.p[0].x
	line1.y = triTranslated.p[1].y - triTranslated.p[0].y
	line1.z = triTranslated.p[1].z - triTranslated.p[0].z

	line2.x = triTranslated.p[2].x - triTranslated.p[0].x
	line2.y = triTranslated.p[2].y - triTranslated.p[0].y
	line2.z = triTranslated.p[2].z - triTranslated.p[0].z

	// normalise normal
	normal.x = line1.y*line2.z - line1.z*line2.y
	normal.y = line1.z*line2.x - line1.x*line2.z
	normal.z = line1.x*line2.y - line1.y*line2.x

	l := math.Sqrt(normal.x*normal.x + normal.y*normal.y + normal.z*normal.z)
	normal.x /= l
	normal.y /= l
	normal.z /= l

	if normal.x*(triTranslated.p[0].x-eng.cam.vec.x)+
		normal.y*(triTranslated.p[0].y-eng.cam.vec.y)+
		normal.z*(triTranslated.p[0].z-eng.cam.vec.z) < 0 {
		// illumination
		light := eng.light.vec
		l = math.Sqrt(light.x*light.x + light.y*light.y + light.z*light.z)
		light.x /= l
		light.y /= l
		light.z /= l

		// how similar is normal to light direction
		dp := normal.x*light.x + normal.y*light.y + normal.z*light.z

		// project triangles from 3D to 2D
		triProjected.p[0] = multiplyMatrixVector(triTranslated.p[0], matProj)
		triProjected.p[1] = multiplyMatrixVector(triTranslated.p[1], matProj)
		triProjected.p[2] = multiplyMatrixVector(triTranslated.p[2], matProj)

		// scale into view
		triProjected.p[0].x += 1
		triProjected.p[1].x += 1
		triProjected.p[2].x += 1

		triProjected.p[0].y += 1
		triProjected.p[1].y += 1
		triProjected.p[2].y += 1

		triProjected.p[0].x *= 0.5 * float64(eng.width)
		triProjected.p[1].x *= 0.5 * float64(eng.width)
		triProjected.p[2].x *= 0.5 * float64(eng.width)

		triProjected.p[0].y *= 0.5 * float64(eng.height)
		triProjected.p[1].y *= 0.5 * float64(eng.height)
		triProjected.p[2].y *= 0.5 * float64(eng.height)
		triProjected.c = GetColor(dp)

		return true, &triProjected
	}

	return false, &triangle{}
}

func fillLinesFromTriangle(tri *triangle, eng *engine) []fyne.CanvasObject {
	vec1 := tri.p[0]
	vec2 := tri.p[1]
	vec3 := tri.p[2]

	// sort vertices ascending by y
	if vec1.y > vec2.y {
		vec1, vec2 = vec2, vec1
	}
	if vec1.y > vec3.y {
		vec1, vec3 = vec3, vec1
	}
	if vec2.y > vec3.y {
		vec2, vec3 = vec3, vec2
	}

	lines := []fyne.CanvasObject{}

	// draw accuracy
	dy := 0.5

	if vec2.y == vec3.y {
		lines = append(lines, fillLinesFromBottomFlatTriangle(vec1, vec2, vec3, eng, tri.c, dy)...)
	} else if vec1.y == vec2.y {
		lines = append(lines, fillLinesFromBottomFlatTriangle(vec1, vec2, vec3, eng, tri.c, dy)...)
	} else {
		vec4 := vec3d{x: vec1.x + ((vec2.y-vec1.y)/(vec3.y-vec1.y))*(vec3.x-vec1.x), y: vec2.y}
		lines = append(lines, fillLinesFromBottomFlatTriangle(vec1, vec2, vec4, eng, tri.c, dy)...)
		lines = append(lines, fillLinesFromTopFlatTriangle(vec2, vec4, vec3, eng, tri.c, dy)...)
	}

	return lines
}

func fillLinesFromBottomFlatTriangle(vec1, vec2, vec3 vec3d, eng *engine, c color.Color, dy float64) []fyne.CanvasObject {
	lines := []fyne.CanvasObject{}
	invslope1 := (vec2.x - vec1.x) / (vec2.y - vec1.y)
	invslope2 := (vec3.x - vec1.x) / (vec3.y - vec1.y)
	curx1 := vec1.x
	curx2 := vec1.x

	for scanlineY := vec1.y; scanlineY <= vec2.y; scanlineY += dy {
		lines = append(lines, newLineBetween(vec3d{x: curx1, y: scanlineY}, vec3d{x: curx2, y: scanlineY}, c, float32(dy)))
		curx1 += invslope1 * dy
		curx2 += invslope2 * dy
	}

	return lines
}

func fillLinesFromTopFlatTriangle(vec1, vec2, vec3 vec3d, eng *engine, c color.Color, dy float64) []fyne.CanvasObject {
	lines := []fyne.CanvasObject{}
	invslope1 := (vec3.x - vec1.x) / (vec3.y - vec1.y)
	invslope2 := (vec3.x - vec2.x) / (vec3.y - vec2.y)
	curx1 := vec3.x
	curx2 := vec3.x

	for scanlineY := vec3.y; scanlineY > vec1.y; scanlineY -= dy {
		lines = append(lines, newLineBetween(vec3d{x: curx1, y: scanlineY}, vec3d{x: curx2, y: scanlineY}, c, 1))
		curx1 -= invslope1 * dy
		curx2 -= invslope2 * dy
	}

	return lines
}

func GetColor(dp float64) color.Color {
	return color.Gray{uint8(math.MaxUint8 * dp)}
}

func newLineBetween(vec1, vec2 vec3d, color color.Color, width float32) *canvas.Line {
	return &canvas.Line{
		Position1:   fyne.NewPos(float32(vec1.x), float32(vec1.y)),
		Position2:   fyne.NewPos(float32(vec2.x), float32(vec2.y)),
		StrokeColor: color,
		StrokeWidth: width,
	}
}
