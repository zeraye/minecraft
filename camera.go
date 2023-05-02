package minecraft

type camera struct {
	fNear, fFar, fFov, fAspectRatio, fFovRad, fTheta float64
	vec                                              vec3d
}
