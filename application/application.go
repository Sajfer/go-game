package application

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sajfer/go-game/camera"
	"github.com/sajfer/go-game/shaders"
	"github.com/sajfer/go-game/texture"
)

// Application Creates window and camera
type Application struct {
	Window *glfw.Window
	Camera *camera.Camera
}

// NewApplication Return an Application object
func NewApplication(width, height int, title string) (*Application, error) {

	runtime.LockOSThread()

	a := new(Application)

	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	gl.Enable(gl.DEPTH_TEST)

	a.Camera, _ = camera.NewCamera()

	a.Window = window

	return a, nil
}

func (a *Application) Draw(vao uint32, shader *shaders.Shader, texture1 *texture.Texture, texture2 *texture.Texture) {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture1.Handle)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, texture2.Handle)

	shader.Use()

	width, height := a.Window.GetSize()

	model := mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{0.5, 1.0, 0.0})
	view := mgl32.LookAtV(a.Camera.CameraPos, a.Camera.CameraPos.Add(a.Camera.CameraFront), a.Camera.CameraUp)
	projection := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/float32(height), 0.1, 100.0)

	mvp := projection.Mul4(view).Mul4(model)

	shader.SetMat4("mvp", mvp)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)

	a.Window.SwapBuffers()
	glfw.PollEvents()
}

func (a *Application) ProcessInput() {
	if a.Window.GetKey(glfw.KeyW) == glfw.Press {
		a.Camera.CameraPos = a.Camera.CameraPos.Add(a.Camera.CameraFront.Mul(a.Camera.CameraSpeed))
	}
	if a.Window.GetKey(glfw.KeyS) == glfw.Press {
		a.Camera.CameraPos = a.Camera.CameraPos.Sub(a.Camera.CameraFront.Mul(a.Camera.CameraSpeed))
	}
	if a.Window.GetKey(glfw.KeyA) == glfw.Press {
		a.Camera.CameraPos = a.Camera.CameraPos.Sub((a.Camera.CameraFront.Cross(a.Camera.CameraUp).Normalize()).Mul(a.Camera.CameraSpeed))
	}
	if a.Window.GetKey(glfw.KeyD) == glfw.Press {
		a.Camera.CameraPos = a.Camera.CameraPos.Add((a.Camera.CameraFront.Cross(a.Camera.CameraUp).Normalize()).Mul(a.Camera.CameraSpeed))
	}
}
