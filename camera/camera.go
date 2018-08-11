package camera

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	CameraPos   mgl32.Vec3
	CameraFront mgl32.Vec3
	CameraUp    mgl32.Vec3
	CameraSpeed float32
	View        mgl32.Mat4
}

func NewCamera() (*Camera, error) {

	c := new(Camera)

	c.CameraPos = mgl32.Vec3{0, 0, 3}
	c.CameraFront = mgl32.Vec3{0, 0, -1}
	c.CameraUp = mgl32.Vec3{0, 1, 0}
	c.CameraSpeed = 0.1

	return c, nil
}
