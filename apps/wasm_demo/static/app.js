const screenWidth = 272
const screenHeight = 240
const audioSampRate = 44100
const audioRenderPeriod = 0.05
let goMemArr
let canvasCtx, canvasData, fbOffset
let audioCtx, audioBuf, abOffset

window.goFuncs = {}
window.copyFromJsArr = (arr, ptr) => {
  goMemArr.set(arr, ptr)
}
window.setFrameBuffer = (ptr) => {
  fbOffset = ptr
}
window.updateScreen = () => {
  canvasData.data.set(goMemArr.slice(fbOffset, fbOffset + screenWidth * screenHeight * 4))
  canvasCtx.putImageData(canvasData, 0, 0)
}
window.setAudioBuffer = (ptr) => {
  abOffset = ptr
}
window.updateAudio = () => {
  return
  audioBuf.getChannelData(0).set(goMemArr.slice(abOffset, abOffset + audioSampRate * audioRenderPeriod))
  audioSrc = audioCtx.createBufferSource()
  audioSrc.buffer = audioBuf
  audioSrc.connect(audioCtx.destination)
  audioSrc.start()
}

let onRomFileOpened = event => {
  let fileArr = new Uint8Array(event.target.result)
  let ret = window.goFuncs.start(fileArr, fileArr.length, audioSampRate)
  if (ret === true) {
    ret = "running."
    document.getElementById("file-input").hidden = true
    onEmuStart()
  }
  document.getElementById("msg").innerText = ret
}

let onEmuStart = () => {
  canvasCtx = document.getElementById("screen").getContext("2d", {alpha: false})
  canvasData = canvasCtx.createImageData(screenWidth, screenHeight)

  audioCtx = new window.AudioContext()
  audioBuf = audioCtx.createBuffer(1, audioSampRate * audioRenderPeriod, audioSampRate)

  document.onkeydown = event => {
    window.goFuncs.onKey(event.code, true)
  }
  document.onkeyup = event => {
    window.goFuncs.onKey(event.code, false)
  }
}

const go = new Go()
WebAssembly.instantiateStreaming(fetch("go_main.wasm"), go.importObject).
  then(res => {
    goMemArr = new Uint8Array(res.instance.exports.mem.buffer)
    go.run(res.instance)
  })
window.onload = () => {
  let fileElem = document.getElementById("file-input")
  fileElem.addEventListener("change", event => {
    let reader = new FileReader()
    reader.onload = onRomFileOpened
    if (event.target.files.length > 0) {
      reader.readAsArrayBuffer(event.target.files[0])
    }
  })
}
