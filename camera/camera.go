package camera

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Camera movements
const (
	FORWARD = iota
	BACKWARD
	LEFT
	RIGHT
)

const (
	cYaw         = -90
	cPitch       = 0
	cSpeed       = 2.5
	cSensitivity = 0.1
	cZoom        = 45
)

// Camera Contains information about the camera
type Camera struct {
	position mgl32.Vec3
	front    mgl32.Vec3
	up       mgl32.Vec3
	right    mgl32.Vec3
	worldUp  mgl32.Vec3
	// Euler angles
	yaw   float64
	pitch float64
	// Options
	moveSpeed   float32
	sensitivity float32
	inverted    bool
}

// NewCamera Returns a new camera object
func NewCamera(inverted bool) (*Camera, error) {

	c := new(Camera)

	c.position = mgl32.Vec3{0, 0, 3}
	c.front = mgl32.Vec3{0, 0, -1}
	c.worldUp = mgl32.Vec3{0, 1, 0}
	c.moveSpeed = cSpeed
	c.yaw = cYaw
	c.pitch = cPitch
	c.sensitivity = cSensitivity
	c.inverted = inverted

	c.updateCameraVectors()

	return c, nil
}

// GetViewMatrix return the View Matrix
func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.position, c.position.Add(c.front), c.up)
}

// ProcessKeyboard Handles keyboard input
func (c *Camera) ProcessKeyboard(direction int, deltaTime float32) {

	velocity := c.moveSpeed * deltaTime
	if direction == FORWARD {
		c.position = c.position.Add(c.front.Mul(velocity))
	}
	if direction == BACKWARD {
		c.position = c.position.Sub(c.front.Mul(velocity))
	}
	if direction == LEFT {
		c.position = c.position.Sub(c.right.Mul(velocity))
	}
	if direction == RIGHT {
		c.position = c.position.Add(c.right.Mul(velocity))
	}
}

// ProcessMouseMovement Handles input of mouse movement data
func (c *Camera) ProcessMouseMovement(xoffset, yoffset float64, constrainPitch bool) {
	xoffset *= float64(c.sensitivity)
	yoffset *= float64(c.sensitivity)

	c.yaw += xoffset
	if c.inverted {
		c.pitch -= yoffset
	} else {
		c.pitch += yoffset
	}

	if constrainPitch {
		if c.pitch > 89.0 {
			c.pitch = 89
		}
		if c.pitch < -89.0 {
			c.pitch = -89
		}
	}
	c.updateCameraVectors()
}

func (c *Camera) updateCameraVectors() {
	var frontVec mgl32.Vec3

	frontVec[0] = float32(math.Cos(float64(mgl32.DegToRad(float32(c.yaw))))) * float32(math.Cos(float64(mgl32.DegToRad(float32(c.pitch)))))
	frontVec[1] = float32(math.Sin(float64(mgl32.DegToRad(float32(c.pitch)))))
	frontVec[2] = float32(math.Sin(float64(mgl32.DegToRad(float32(c.yaw))))) * float32(math.Cos(float64(mgl32.DegToRad(float32(c.pitch)))))

	c.front = frontVec.Normalize()

	c.right = c.front.Cross(c.worldUp).Normalize()
	c.up = c.right.Cross(c.front).Normalize()
}
