package minecraft

type vec3d struct {
	x, y, z float64
}

type triange struct {
	p [3]vec3d
}

type mesh struct {
	tris []triange
}

type mat4x4 struct {
	m [4][4]float64
}

type iengine interface {
	MultiplyMatrixVector()
}

type sengine struct {
	meshCube mesh
	matProj  mat4x4
}
