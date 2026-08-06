// Package webpointer は、ブラウザで Ebitengine が返すポインタ位置に乗って
// しまう canvas の位置を打ち消すためのもの。
//
// Ebitengine の js バックエンドは、マウスもタッチも clientX / clientY を
// そのままポインタ位置として使う（internal/ui/input_js.go の
// setMouseCursorFromEvent と updateTouchesFromEvent）。canvas がビューポートの
// 左上にある前提の計算なので、web/<game>/index.html のように canvas を中央へ
// 寄せているページでは、canvas の左上の位置だけ座標がずれる。タップした場所で
// 当たり判定をするゲームは、Offset が返す値を引く必要がある。
//
// canvas の中の余白（論理画面がぴったり収まらないときにできる letterbox）は
// Ebitengine 自身が差し引いているので、ここで面倒を見るのは canvas の外側の
// 位置だけでいい。
package webpointer
