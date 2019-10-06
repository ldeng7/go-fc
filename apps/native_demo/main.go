package main

import (
	"flag"
	"os"
	"path"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/ldeng7/go-fc/core"
)

type conf struct {
	romPath  string
	patchTyp uint64
	tvFormat uint
}

func parseArgs() *conf {
	c := &conf{}
	flag.StringVar(&c.romPath, "rom", "", "rom path")
	flag.Uint64Var(&c.patchTyp, "patch", 0, "patch type")
	flag.UintVar(&c.tvFormat, "tv", 0, "tv format: 0=ntsc, 1=pal, 2=pal-china")
	flag.Parse()
	if len(c.romPath) == 0 {
		flag.PrintDefaults()
		return nil
	}
	if c.tvFormat > 2 {
		println("invalid tv format")
		return nil
	}
	return c
}

type App struct {
	t, te   float64
	sys     *core.Sys
	audio   *Audio
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

	_, filename := path.Split(c.romPath)
	if a.graphic, err = newGraphic(filename); err != nil {
		return nil, err
	}

	if a.audio, err = newAudio(); err != nil {
		return nil, err
	}

	f, err := os.Open(c.romPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	ac := &core.Conf{
		PatchTyp:      c.patchTyp,
		TvFormat:      byte(c.tvFormat),
		AllSprite:     true,
		AudioSampRate: a.audio.sampRate,
	}
	if a.sys, err = core.NewSys(f, ac); err != nil {
		return nil, err
	}
	a.audio.source = a.sys.GetAudioDataQueue()

	a.graphic.window.SetKeyCallback(a.onKey)
	return a, nil
}

func (a *App) deInit() {
	if a.audio != nil {
		a.audio.deInit()
	}
	if a.graphic != nil {
		a.graphic.deInit()
	}
}

func (a *App) run() error {
	if err := a.audio.stream.Start(); err != nil {
		return err
	}

	sys, win := a.sys, a.graphic.window
	p := float64(sys.GetFramePeriod()) / 1000.0
	a.te = glfw.GetTime()
	for !win.ShouldClose() {
		sys.SetFrameBuffer(a.graphic.fb)
		a.t = glfw.GetTime()
		for ; a.te < a.t; a.te += p {
			sys.RunFrame()
		}
		a.graphic.runFrame()
		glfw.PollEvents()
	}
	return nil
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
	if err = app.run(); err != nil {
		println(err.Error())
		return
	}
	app.deInit()
}
