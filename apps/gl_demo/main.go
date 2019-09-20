package main

import (
	"flag"
	"os"
	"path"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/ldeng7/go-fc/core"
)

type conf struct {
	romPath  string
	patchTyp uint64
	isPal    bool
}

func parseArgs() *conf {
	c := &conf{}
	flag.StringVar(&c.romPath, "rom", "", "rom path")
	flag.Uint64Var(&c.patchTyp, "patch", 0, "patch type")
	flag.BoolVar(&c.isPal, "pal", false, "run in PAL mode instead of NTSC")
	flag.Parse()
	if len(c.romPath) == 0 {
		flag.PrintDefaults()
		return nil
	}
	return c
}

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

func newApp(c *conf) (*App, error) {
	a := &App{}
	var err error
	defer func() {
		if err != nil {
			a.deInit()
		}
	}()

	f, err := os.Open(c.romPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		f.Close()
	}()
	ac := &core.Conf{
		PatchTyp:  c.patchTyp,
		IsPal:     c.isPal,
		AllSprite: true,
	}
	if a.sys, err = core.NewSys(f, ac); err != nil {
		return nil, err
	}

	_, filename := path.Split(c.romPath)
	if a.graphic, err = newGraphic(filename); err != nil {
		return nil, err
	}

	a.graphic.window.SetKeyCallback(a.onKey)
	return a, nil
}

func (a *App) deInit() {
	if a.graphic != nil {
		a.graphic.deInit()
	}
}

func (a *App) onKey(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyEscape:
			a.sys.Reset()
		}
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
		sys.SetFrameBuffer(a.graphic.fb)
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
	c := parseArgs()
	if nil == c {
		return
	}
	app, err := newApp(c)
	if err != nil {
		println(err.Error())
		return
	}
	app.mainLoop()
	app.deInit()
}
