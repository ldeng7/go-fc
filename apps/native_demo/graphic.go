package main

import (
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/ldeng7/go-fc/core"
)

const (
	screenWidth  = core.ScreenWidth * 2
	screenHeight = core.ScreenHeight * 2
)

type Graphic struct {
	glfwInited bool
	window     *glfw.Window
	texture    uint32
	fb         *core.FrameBuffer
	fbp        unsafe.Pointer
}

func newGraphic(title string) (*Graphic, error) {
	g := &Graphic{}
	var err error
	defer func() {
		if err != nil {
			g.deInit()
		}
	}()

	if err = glfw.Init(); err != nil {
		return nil, err
	}
	g.glfwInited = true
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	g.window, err = glfw.CreateWindow(screenWidth, screenHeight, title, nil, nil)
	if err != nil {
		return nil, err
	}
	g.window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		return nil, err
	}
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &g.texture)
	gl.BindTexture(gl.TEXTURE_2D, g.texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	g.fb = &core.FrameBuffer{}
	g.fbp = unsafe.Pointer(g.fb)
	return g, nil
}

func (g *Graphic) deInit() {
	if g.glfwInited {
		glfw.Terminate()
	}
}

func (g *Graphic) runFrame() {
	gl.BindTexture(gl.TEXTURE_2D, g.texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, core.ScreenWidth, core.ScreenHeight,
		0, gl.RGBA, gl.UNSIGNED_BYTE, g.fbp)

	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-1, -1)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(1, -1)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(1, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-1, 1)
	gl.End()

	gl.BindTexture(gl.TEXTURE_2D, 0)
	g.window.SwapBuffers()
}
