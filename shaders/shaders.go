package shaders

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-gl/mathgl/mgl32"

	"go-game/utils"

	"github.com/go-gl/gl/v4.3-core/gl"
)

// Shader a shader object
type Shader struct {
	ID uint32
}

// Program a program object
type Program struct {
	ID      uint32
	shaders []*Shader
}

func readFile(path string) string {
	path, _ = filepath.Abs(path)
	source, err := ioutil.ReadFile(path)
	utils.Check(err)

	sourceStr := string(source) + "\x00"

	return sourceStr
}

func getGlError(glHandle uint32, checkTrueParam uint32, failMsg string) error {

	var status int32
	gl.GetShaderiv(glHandle, checkTrueParam, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(glHandle, gl.INFO_LOG_LENGTH, &logLength)

		logLength = 512

		log := gl.Str(strings.Repeat("\x00", int(logLength+1)))
		gl.GetShaderInfoLog(glHandle, logLength, nil, log)

		return fmt.Errorf("%s%s", failMsg, gl.GoStr(log))
	}
	return nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	//err := getGlError(shader, gl.COMPILE_STATUS, "SHADER::COMPILE_FAILURE::")
	//if err != nil {
	//		return 0, err
	//	}

	return shader, nil
}

// NewShader return a shader object
func NewShader(vertexPath, fragmentPath string) (*Shader, error) {
	s := new(Shader)

	vertexSourceStr := readFile(vertexPath)
	fragmentSourceStr := readFile(fragmentPath)

	vertex, err := compileShader(vertexSourceStr, gl.VERTEX_SHADER)
	if err != nil {
		println("Vertex: ", err.Error())
		panic(err)
	}
	fragment, err := compileShader(fragmentSourceStr, gl.FRAGMENT_SHADER)
	if err != nil {
		println("fragment: ", err.Error())
		panic(err)
	}
	s.ID = gl.CreateProgram()
	gl.AttachShader(s.ID, vertex)
	gl.AttachShader(s.ID, fragment)
	gl.LinkProgram(s.ID)

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
