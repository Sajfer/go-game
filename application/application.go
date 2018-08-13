package application

import (
	"log"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sajfer/go-game/camera"
	"github.com/sajfer/go-game/shaders"
	"github.com/sajfer/go-game/texture"
)

// Application Creates window and camera
type Application struct {
	Window  *glfw.Window
	Camera  *camera.Camera
	Shaders []*shaders.Shader
}

var (
	firstMouse         = true
	yaw        float64 = -90
	pitch      float64
	lastX      float64 = 400
	lastY      float64 = 300
	fov        float32 = 45
)

// NewApplication Return an Application object
func NewApplication(width, height int, title string) (*Application, error) {
	a := new(Application)

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPosCallback(a.mouseCallback)

	if err := gl.Init(); err != nil {
		return nil, err
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	gl.Enable(gl.DEPTH_TEST)

	a.Camera, _ = camera.NewCamera(true)

	a.Window = window

	return a, nil
}

func (a *Application) mouseCallback(window *glfw.Window, xpos, ypos float64) {

	if firstMouse {
		lastX = xpos
		lastY = ypos
		firstMouse = false
	}
	xoffset := xpos - lastX
	yoffset := ypos - lastY
	lastX = xpos
	lastY = ypos

	a.Camera.ProcessMouseMovement(xoffset, yoffset, true)

}

func (a *Application) Draw(vao uint32, texture1 *texture.Texture, texture2 *texture.Texture) {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture1.Handle)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, texture2.Handle)

	for _, shader := range a.Shaders {
		shader.Use()
	}

	width, height := a.Window.GetSize()

	model := mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{0.5, 1.0, 0.0})
	view := a.Camera.GetViewMatrix()
	projection := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/float32(height), 0.1, 100.0)

	mvp := projection.Mul4(view).Mul4(model)

	for _, shader := range a.Shaders {
		shader.SetMat4("mvp", mvp)
	}

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)

	a.Window.SwapBuffers()
	glfw.PollEvents()
}

// ProcessInput Handles user input
func (a *Application) ProcessInput(deltaTime float32) {

	if a.Window.GetKey(glfw.KeyEscape) == glfw.Press {
		a.Window.SetShouldClose(true)
	}

	if a.Window.GetKey(glfw.KeyW) == glfw.Press {
		a.Camera.ProcessKeyboard(camera.FORWARD, deltaTime)
	}
	if a.Window.GetKey(glfw.KeyS) == glfw.Press {
		a.Camera.ProcessKeyboard(camera.BACKWARD, deltaTime)
	}
	if a.Window.GetKey(glfw.KeyA) == glfw.Press {
		a.Camera.ProcessKeyboard(camera.LEFT, deltaTime)
	}
	if a.Window.GetKey(glfw.KeyD) == glfw.Press {
		a.Camera.ProcessKeyboard(camera.RIGHT, deltaTime)
	}
}

func (a *Application) AddShader(shader *shaders.Shader) {
	a.Shaders = append(a.Shaders, shader)
}

func (a *Application) MainLoop(vao uint32, texture1 *texture.Texture, texture2 *texture.Texture) {

	var deltaTime, lastFrame float32

	if a.Window == nil {
		println("NO WINDOW")
	}

	for !a.Window.ShouldClose() {

		currentFrame := float32(glfw.GetTime())
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		a.ProcessInput(deltaTime)
		a.Draw(vao, texture1, texture2)
	}
}
