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

var keyMapP1 = map[glfw.Key]byte{
	glfw.KeyW:           core.PadKeyUp,
	glfw.KeyS:           core.PadKeyDown,
	glfw.KeyA:           core.PadKeyLeft,
	glfw.KeyD:           core.PadKeyRight,
	glfw.KeySpace:       core.PadKeyStart,
	glfw.KeyLeftControl: core.PadKeySelect,
	glfw.KeyK:           core.PadKeyA,
	glfw.KeyJ:           core.PadKeyB,
}
var keyMapP2 = map[glfw.Key]byte{
	glfw.KeyUp:        core.PadKeyUp,
	glfw.KeyDown:      core.PadKeyDown,
	glfw.KeyLeft:      core.PadKeyLeft,
	glfw.KeyRight:     core.PadKeyRight,
	glfw.KeyKPEnter:   core.PadKeyStart,
	glfw.KeyKPDecimal: core.PadKeySelect,
	glfw.KeyKP2:       core.PadKeyA,
	glfw.KeyKP1:       core.PadKeyB,
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
	switch key {
	case glfw.KeyEscape:
		if action == glfw.Press {
			a.sys.Reset()
		}
	default:
		switch action {
		case glfw.Press:
			if pk, ok := keyMapP1[key]; ok {
				a.sys.SetPadKey(1, pk, true)
			} else if pk, ok := keyMapP2[key]; ok {
				a.sys.SetPadKey(2, pk, true)
			}
		case glfw.Release:
			if pk, ok := keyMapP1[key]; ok {
				a.sys.SetPadKey(1, pk, false)
			} else if pk, ok := keyMapP2[key]; ok {
				a.sys.SetPadKey(2, pk, false)
			}
		}
	}
}

func (a *App) mainLoop() {
	sys, window := a.sys, a.graphic.window
	period := float64(sys.GetFramePeriod()) / 1000.0
	a.te = glfw.GetTime()
	for !window.ShouldClose() {
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
