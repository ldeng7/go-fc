package main

import (
	"bytes"
	"reflect"
	"sync"
	"syscall/js"
	"time"
	"unsafe"

	"github.com/ldeng7/go-fc/core"
)

type Ctx struct {
	sys *core.Sys

	copyFromJsArr  js.Value
	setFrameBuffer js.Value
	updateScreen   js.Value
}

var keyMapP1 = map[string]byte{
	"KeyW":        core.PadKeyUp,
	"KeyS":        core.PadKeyDown,
	"KeyA":        core.PadKeyLeft,
	"KeyD":        core.PadKeyRight,
	"Space":       core.PadKeyStart,
	"ControlLeft": core.PadKeySelect,
	"KeyK":        core.PadKeyA,
	"KeyJ":        core.PadKeyB,
}

func (ctx *Ctx) start(romFileArr js.Value, romFileLen int) interface{} {
	romFile := make([]byte, romFileLen)
	romFilePtr := (*reflect.SliceHeader)(unsafe.Pointer(&romFile)).Data
	ctx.copyFromJsArr.Invoke(romFileArr, romFilePtr)

	sysConf := &core.Conf{
		AllSprite: true,
	}
	sys, err := core.NewSys(bytes.NewReader(romFile), sysConf)
	if err != nil {
		return err.Error()
	}
	ctx.sys = sys

	go func() {
		tSys := time.NewTicker(time.Duration(sys.GetFramePeriod()*1000.0) * time.Microsecond)
		fb := &core.FrameBuffer{}
		sys.SetFrameBuffer(fb)
		ctx.setFrameBuffer.Invoke(uintptr(unsafe.Pointer(fb)))
		for _ = range tSys.C {
			sys.RunFrame()
			ctx.updateScreen.Invoke()
		}
	}()
	return true
}

func (ctx *Ctx) onKey(code string, down bool) {
	switch code {
	case "Escape":
		if down {
			ctx.sys.Reset()
		}
	default:
		if pk, ok := keyMapP1[code]; ok {
			ctx.sys.SetPadKey(1, pk, down)
		}
	}
}

func main() {
	jsGlobal := js.Global()
	ctx := &Ctx{
		copyFromJsArr:  jsGlobal.Get("copyFromJsArr"),
		setFrameBuffer: jsGlobal.Get("setFrameBuffer"),
		updateScreen:   jsGlobal.Get("updateScreen"),
	}
	goFuncs := jsGlobal.Get("goFuncs")
	goFuncs.Set("start", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		return ctx.start(args[0], args[1].Int())
	}))
	goFuncs.Set("onKey", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		ctx.onKey(args[0].String(), args[1].Bool())
		return nil
	}))

	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
