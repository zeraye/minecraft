package minecraft

import (
	"bufio"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type vec3d struct {
	x, y, z float64
}

type triangle struct {
	p [3]vec3d
	c color.Color
}

type mesh struct {
	tris []triangle
}

func loadFromObjectFile(filename string) *mesh {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	verts := []vec3d{}
	meshObject := &mesh{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, " ")
		if tokens[0] == "v" {
			vec := vec3d{}
			if value, err := strconv.ParseFloat(tokens[1], 64); err == nil {
				vec.x = value
			} else {
				panic(err)
			}
			if value, err := strconv.ParseFloat(tokens[2], 64); err == nil {
				vec.y = value
			} else {
				panic(err)
			}
			if value, err := strconv.ParseFloat(tokens[3], 64); err == nil {
				vec.z = value
			} else {
				panic(err)
			}
			verts = append(verts, vec)
		}
		if tokens[0] == "f" {
			tria := triangle{}
			if index, err := strconv.ParseInt(tokens[1], 10, 0); err == nil {
				tria.p[0] = verts[index-1]
			} else {
				panic(err)
			}
			if index, err := strconv.ParseInt(tokens[2], 10, 0); err == nil {
				tria.p[1] = verts[index-1]
			} else {
				panic(err)
			}
			if index, err := strconv.ParseInt(tokens[3], 10, 0); err == nil {
				tria.p[2] = verts[index-1]
			} else {
				panic(err)
			}
			meshObject.tris = append(meshObject.tris, tria)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return meshObject
}

type mat4x4 struct {
	m [4][4]float64
}

func multiplyMatrixVector(vec vec3d, mat mat4x4) vec3d {
	oVec := vec3d{}
	oVec.x = vec.x*mat.m[0][0] + vec.y*mat.m[1][0] + vec.z*mat.m[2][0] + mat.m[3][0]
	oVec.y = vec.x*mat.m[0][1] + vec.y*mat.m[1][1] + vec.z*mat.m[2][1] + mat.m[3][1]
	oVec.z = vec.x*mat.m[0][2] + vec.y*mat.m[1][2] + vec.z*mat.m[2][2] + mat.m[3][2]

	w := vec.x*mat.m[0][3] + vec.y*mat.m[1][3] + vec.z*mat.m[2][3] + mat.m[3][3]
	if w != 0 {
		oVec.x /= w
		oVec.y /= w
		oVec.z /= w
	}

	return oVec
}

func createRotationZMatrix(fTheta float64) mat4x4 {
	return mat4x4{
		[4][4]float64{
			{math.Cos(fTheta), math.Sin(fTheta), 0, 0},
			{-math.Sin(fTheta), math.Cos(fTheta), 0, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 1}}}
}

func createRotationXMatrix(fTheta float64) mat4x4 {
	return mat4x4{
		[4][4]float64{
			{1, 0, 0, 0},
			{0, math.Cos(fTheta * 0.5), math.Sin(fTheta * 0.5), 0},
			{0, -math.Sin(fTheta * 0.5), math.Cos(fTheta * 0.5), 0},
			{0, 0, 0, 1}}}
}

func createProjectionMatrix(fAspectRatio, fFovRad, fFar, fNear float64) mat4x4 {
	return mat4x4{
		[4][4]float64{
			{fAspectRatio * fFovRad, 0, 0, 0},
			{0, fFovRad, 0, 0},
			{0, 0, fFar / (fFar - fNear), (-fFar * fNear) / (fFar - fNear)},
			{0, 0, (-fFar * fNear) / (fFar - fNear), 0}}}
}
