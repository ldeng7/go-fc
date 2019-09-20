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

func start(ctx *Ctx, romFileArr js.Value, romFileLen int) interface{} {
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
		ticker := time.NewTicker(time.Duration(sys.GetFramePeriod()*1000.0) * time.Microsecond)
		fb := &core.FrameBuffer{}
		sys.SetFrameBuffer(fb)
		ctx.setFrameBuffer.Invoke(uintptr(unsafe.Pointer(fb)))
		for _ = range ticker.C {
			sys.RunFrame()
			ctx.updateScreen.Invoke()
		}
	}()
	return true
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
		return start(ctx, args[0], args[1].Int())
	}))

	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
