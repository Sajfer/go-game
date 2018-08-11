package shaders

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/sajfer/go-game/utils"
)

// Shader a shader object
type Shader struct {
	ID uint32
}

var errShaderCompilationError = errors.New("Could not compile shader")
var errShaderLinkError = errors.New("Could not link shader")

// NewShader return a shader object
func NewShader(vertexPath, fragmentPath string) (*Shader, error) {
	s := new(Shader)

	vertexPath, _ = filepath.Abs(vertexPath)
	fragmentPath, _ = filepath.Abs(fragmentPath)

	vertexSource, err := ioutil.ReadFile(vertexPath)
	utils.Check(err)

	fragmentSource, err := ioutil.ReadFile(fragmentPath)
	utils.Check(err)

	vertexSourceStr := string(vertexSource) + "\x00"
	fragmentSourceStr := string(fragmentSource) + "\x00"

	var status int32

	cvertexSource, free := gl.Strs(vertexSourceStr)
	vertex := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertex, 1, cvertexSource, nil)
	free()
	gl.CompileShader(vertex)
	gl.GetShaderiv(vertex, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertex, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vertex, logLength, nil, gl.Str(log))
		println("failed to compile %v: %v", cvertexSource, log)
		return nil, errShaderCompilationError
	}

	cfragmentSource, free := gl.Strs(fragmentSourceStr)
	fragment := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragment, 1, cfragmentSource, nil)
	free()
	gl.CompileShader(fragment)
	gl.GetShaderiv(fragment, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fragment, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(fragment, logLength, nil, gl.Str(log))
		println("failed to compile %v: %v", cfragmentSource, log)
		return nil, errShaderCompilationError
	}

	s.ID = gl.CreateProgram()
	gl.AttachShader(s.ID, vertex)
	gl.AttachShader(s.ID, fragment)
	gl.LinkProgram(s.ID)

	gl.GetShaderiv(fragment, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fragment, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(fragment, logLength, nil, gl.Str(log))
		println("failed to link: %v", log)
		return nil, errShaderLinkError
	}

	gl.DeleteShader(vertex)
	gl.DeleteShader(fragment)

	return s, nil
}

// Use run UseProgram on shader program
func (s *Shader) Use() {
	gl.UseProgram(s.ID)
}

// SetBool Set boolean Uniform value
func (s *Shader) SetBool(name string, value bool) {
	var valueInt int32
	if value {
		valueInt = 1
	}
	gl.Uniform1i(gl.GetUniformLocation(uint32(s.ID), gl.Str(name+"\x00")), valueInt)
}

// SetInt Set int Uniform value
func (s *Shader) SetInt(name string, value int) {
	gl.Uniform1i(gl.GetUniformLocation(uint32(s.ID), gl.Str(name+"\x00")), int32(value))
}

// SetFloat Set float Uniform value
func (s *Shader) SetFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(uint32(s.ID), gl.Str(name+"\x00")), value)
}

// SetVec2 Set Vec2 uniform value
func (s *Shader) SetVec2(name string, value mgl32.Vec2) {
	gl.Uniform2fv(gl.GetUniformLocation(uint32(s.ID), gl.Str(name+"\x00")), 1, &value[0])
}

// SetVec2f Set vec2f uniform value
func (s *Shader) SetVec2f(name string, x, y float32) {
	gl.Uniform2f(gl.GetUniformLocation(uint32(s.ID), gl.Str(name+"\x00")), x, y)
}

// SetVec3 Set vec3 uniform value
func (s *Shader) SetVec3(name string, value mgl32.Vec3) {
	gl.Uniform3fv(gl.GetUniformLocation(uint32(s.ID), gl.Str(name+"\x00")), 1, &value[0])
}

// SetVec3f Set vec3f uniform value
func (s *Shader) SetVec3f(name string, x, y, z float32) {
	gl.Uniform3f(gl.GetUniformLocation(uint32(s.ID), gl.Str(name+"\x00")), x, y, z)
}

// SetVec4 Set vec4 uniform value
func (s *Shader) SetVec4(name string, value mgl32.Vec4) {
	gl.Uniform4fv(gl.GetUniformLocation(uint32(s.ID), gl.Str(name+"\x00")), 1, &value[0])
}

// SetVec4f Set vec4f uniform value
func (s *Shader) SetVec4f(name string, x, y, z, w float32) {
	gl.Uniform4f(gl.GetUniformLocation(uint32(s.ID), gl.Str(name+"\x00")), x, y, z, w)
}

// SetMat2 Set Mat2 unform value
func (s *Shader) SetMat2(name string, mat mgl32.Mat2) {
	gl.UniformMatrix2fv(gl.GetUniformLocation(uint32(s.ID), gl.Str(name+"\x00")), 1, false, &mat[0])
}

// SetMat3 Set Mat3 unform value
func (s *Shader) SetMat3(name string, mat mgl32.Mat3) {
	gl.UniformMatrix3fv(gl.GetUniformLocation(uint32(s.ID), gl.Str(name+"\x00")), 1, false, &mat[0])
}

// SetMat4 Set Mat4 unform value
func (s *Shader) SetMat4(name string, mat mgl32.Mat4) {
	gl.UniformMatrix4fv(gl.GetUniformLocation(uint32(s.ID), gl.Str(name+"\x00")), 1, false, &mat[0])
}
