package minecraft

import (
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type engine struct {
	app           fyne.App
	window        fyne.Window
	meshes        []*mesh
	meshDrawer    *meshDrawer
	width, height float64
	fps           float64
	cam           *camera
	lastUpdate    time.Time
	light         *lightSource
	start         time.Time
}

func New() *engine {
	app := app.New()
	window := app.NewWindow("Minecraft")

	fps := 15.0
	width := 1280.0
	height := 720.0

	fNear := 0.5
	fFar := 100.0
	fFov := 90.0
	fAspectRatio := height / width
	fFovRad := 1.0 / math.Tan(fFov*0.5/180.0*math.Pi)
	camVec := vec3d{0, 0, 0}
	fTheta := 0.0

	return &engine{
		app:        app,
		window:     window,
		meshes:     []*mesh{},
		meshDrawer: &meshDrawer{},
		width:      width,
		height:     height,
		fps:        fps,
		cam:        &camera{fNear, fFar, fFov, fAspectRatio, fFovRad, fTheta, camVec},
		lastUpdate: time.Now(),
		light:      &lightSource{vec3d{0, 0, -1}},
		start:      time.Now(),
	}
}

func (eng *engine) Run() {
	eng.window.Resize(fyne.NewSize(float32(eng.width), float32(eng.height)))

	// temporary generate cube, remove in the future
	meshShip := loadFromObjectFile("../objects/ship.obj")

	eng.meshes = append(eng.meshes, meshShip)

	combinedMeshObject := &mesh{}

	for _, meshOjbect := range eng.meshes {
		combinedMeshObject.tris = append(combinedMeshObject.tris, meshOjbect.tris...)
	}

	eng.meshDrawer = newMeshDrawer(combinedMeshObject, eng)

	go func() {
		for range time.Tick(time.Duration(float64(time.Second) / eng.fps)) {
			eng.cam.fTheta = float64(eng.lastUpdate.UnixMilli()-eng.start.UnixMilli()) / 1000
			eng.update()
		}
	}()

	eng.window.ShowAndRun()
}

func (eng *engine) update() {
	eng.lastUpdate = time.Now()
	eng.window.SetContent(eng.meshDrawer)
}
