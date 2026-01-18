package renderer

import (
	"image"
	"image/draw"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	rgba    *image.RGBA
	texture uint32
)

const (
	vertexShaderSource = `
		#version 410
		in vec3 vp;
		in vec2 vt;
		out vec2 fragTexCoord;
		void main() {
			fragTexCoord = vt;
			gl_Position = vec4(vp, 1.0);
		}` + "\x00"

	fragmentShaderSource = `
		#version 410
		in vec2 fragTexCoord;
		out vec4 outputColor;
		uniform sampler2D tex;
		void main() {
			outputColor = texture(tex, fragTexCoord);
		}` + "\x00"
)

// Start initializes the window and runs the game loop.
func Start(width, height int, nextFrame func() image.Image) {
	runtime.LockOSThread()

	window := initGlfw(width, height)
	defer glfw.Terminate()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	program := initShaders()
	vao := makeVao()

	// 1. Setup Texture once
	rgba = image.NewRGBA(image.Rect(0, 0, width, height))
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)

	for !window.ShouldClose() {
		img := nextFrame()
		updateTexture(img)

		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(program)
		gl.BindVertexArray(vao)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func initGlfw(width, height int) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Go-GL High Performance", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	return window
}

func initShaders() uint32 {
	vShader := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fShader := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vShader)
	gl.AttachShader(prog, fShader)
	gl.LinkProgram(prog)
	return prog
}

func updateTexture(img image.Image) {
	// Copy pixel data from image.Image to our internal RGBA buffer
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexSubImage2D(
		gl.TEXTURE_2D, 0, 0, 0,
		int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y),
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix),
	)
}

func makeVao() uint32 {
	var vertices = []float32{
		-1.0, 1.0, 0.0, 0.0, 0.0,
		-1.0, -1.0, 0.0, 0.0, 1.0,
		1.0, -1.0, 0.0, 1.0, 1.0,
		-1.0, 1.0, 0.0, 0.0, 0.0,
		1.0, -1.0, 0.0, 1.0, 1.0,
		1.0, 1.0, 0.0, 1.0, 0.0,
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0) // Position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1) // Texture Coords
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	return vao
}

func compileShader(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		panic(log)
	}
	return shader
}
