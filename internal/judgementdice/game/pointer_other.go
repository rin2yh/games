//go:build !js

package game

// pointerOffset はブラウザ以外では常に 0。デスクトップ版はポインタ位置が
// そのまま論理座標で返ってくる。
func pointerOffset() (x, y int) {
	return 0, 0
}
