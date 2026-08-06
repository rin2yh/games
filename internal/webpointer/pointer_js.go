//go:build js

package webpointer

import "syscall/js"

// Offset はブラウザで canvas の位置のぶんを返す。logicalWidth と
// logicalHeight にはゲームの Layout が返す論理サイズを渡す。
func Offset(logicalWidth, logicalHeight float64) (x, y int) {
	canvas := js.Global().Get("document").Call("querySelector", "canvas")
	if !canvas.Truthy() {
		return 0, 0
	}
	rect := canvas.Call("getBoundingClientRect")
	w := rect.Get("width").Float()
	h := rect.Get("height").Float()
	if w <= 0 || h <= 0 {
		return 0, 0
	}

	// Ebitengine は論理画面を canvas に等倍率で収める（余った側は余白になる）。
	// 同じ倍率で CSS ピクセルを論理ピクセルに戻す。
	scale := min(w/logicalWidth, h/logicalHeight)
	return int(rect.Get("left").Float() / scale), int(rect.Get("top").Float() / scale)
}
