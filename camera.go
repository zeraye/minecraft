package minecraft

type camera struct {
	fNear, fFar, fFov, fAspectRatio, fFovRad float64
	vec                                      vec3d
}
