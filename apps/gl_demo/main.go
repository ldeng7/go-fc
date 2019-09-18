package main

import (
	"os"
	"path"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/ldeng7/go-fc/core"
)

const (
	screenWidth  = core.ScreenWidth * 2
	screenHeight = core.ScreenHeight * 2
)

var keyMapP1 = map[byte]glfw.Key{
	core.PadKeyUp:     glfw.KeyW,
	core.PadKeyDown:   glfw.KeyS,
	core.PadKeyLeft:   glfw.KeyA,
	core.PadKeyRight:  glfw.KeyD,
	core.PadKeyStart:  glfw.KeySpace,
	core.PadKeySelect: glfw.KeyLeftAlt,
	core.PadKeyA:      glfw.KeyK,
	core.PadKeyB:      glfw.KeyJ,
}
var keyMapP2 = map[byte]glfw.Key{
	core.PadKeyUp:     glfw.KeyUp,
	core.PadKeyDown:   glfw.KeyDown,
	core.PadKeyLeft:   glfw.KeyLeft,
	core.PadKeyRight:  glfw.KeyRight,
	core.PadKeyStart:  glfw.KeyKPEnter,
	core.PadKeySelect: glfw.KeyKPDecimal,
	core.PadKeyA:      glfw.KeyKP2,
	core.PadKeyB:      glfw.KeyKP1,
}

type App struct {
	t, te   float64
	sys     *core.Sys
	graphic *Graphic
}

func newApp() (*App, error) {
	a := &App{}
	var err error
	defer func() {
		if err != nil {
			a.deInit()
		}
	}()

	f, err := os.Open(os.Args[1])
	if err != nil {
		return nil, err
	}
	defer func() {
		f.Close()
	}()
	sysConf := &core.Conf{
		IsPal: false,
	}
	if a.sys, err = core.NewSys(f, sysConf); err != nil {
		return nil, err
	}

	_, filename := path.Split(os.Args[1])
	if a.graphic, err = newGraphic(filename); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) deInit() {
	if a.graphic != nil {
		a.graphic.deInit()
	}
}

func (a *App) mainLoop() {
	sys, window := a.sys, a.graphic.window
	period := float64(sys.GetFramePeriod()) / 1000.0
	a.te = glfw.GetTime()
	for !window.ShouldClose() {
		var pk1, pk2 byte
		for pk, k := range keyMapP1 {
			if window.GetKey(k) == glfw.Press {
				pk1 |= pk
			}
		}
		sys.SetPadKey(1, pk1)
		for pk, k := range keyMapP2 {
			if window.GetKey(k) == glfw.Press {
				pk2 |= pk
			}
		}
		sys.SetPadKey(2, pk2)

		gl.Clear(gl.COLOR_BUFFER_BIT)
		sys.SetFrameBuffer(a.graphic.fbf)
		a.t = glfw.GetTime()
		for ; a.te < a.t; a.te += period {
			sys.RunFrame()
		}
		a.graphic.glStepPost()
		glfw.PollEvents()
	}
}

func init() {
	runtime.LockOSThread()
}

func main() {
	app, err := newApp()
	if err != nil {
		println(err.Error())
		return
	}
	app.mainLoop()
	app.deInit()
}
