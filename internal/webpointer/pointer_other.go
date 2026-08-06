//go:build !js

package webpointer

// Offset はブラウザ以外では常に 0。デスクトップ版はポインタ位置がそのまま
// 論理座標で返ってくる。
func Offset(logicalWidth, logicalHeight float64) (x, y int) {
	return 0, 0
}
