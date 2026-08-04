//go:build js

package game

import "syscall/js"

// pointerOffset は、Ebitengine が返すポインタ位置に乗っている canvas の位置の
// ぶんを、論理座標に直して返す。
//
// Ebitengine の js バックエンドは、マウスもタッチも clientX / clientY を
// そのままポインタ位置として使う（internal/ui/input_js.go の
// setMouseCursorFromEvent と updateTouchesFromEvent）。canvas がビューポートの
// 左上にある前提の計算なので、web/judgement-dice/index.html のように body を
// 中央へ寄せているページでは、canvas の左上の位置だけ座標がずれる。
// タップした場所で当たり判定をするにはこの差を引く必要がある。
func pointerOffset() (x, y int) {
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
	scale := min(w/screenWidth, h/screenHeight)
	return int(rect.Get("left").Float() / scale), int(rect.Get("top").Float() / scale)
}
