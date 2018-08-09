package shaders

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/sajfer/graph3d/utils"
)

// Shader a shader object
type Shader struct {
	ID           uint32
	vertexPath   string
	fragmentPath string
}

// NewShader return a shader object
func NewShader(vertexPath, fragmentPath string) *Shader {
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
		return nil
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
		return nil
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
		return nil
	}

	gl.DeleteShader(vertex)
	gl.DeleteShader(fragment)

	return s
}

func (s *Shader) use() {
	gl.UseProgram(s.ID)
}

func (s *Shader) setBool(name string, value bool) {
	var valueVar int32
	if value {
		valueVar = 1
	}
	cName, free := gl.Strs(name)
	gl.Uniform1i(gl.GetUniformLocation(uint32(s.ID), *cName), valueVar)
	free()
}

func (s *Shader) setInt(name string, value int) {
	cName, free := gl.Strs(name)
	gl.Uniform1i(gl.GetUniformLocation(uint32(s.ID), *cName), int32(value))
	free()
}

func (s *Shader) setFloat(name string, value float32) {
	cName, free := gl.Strs(name)
	gl.Uniform1f(gl.GetUniformLocation(uint32(s.ID), *cName), value)
	free()
}
