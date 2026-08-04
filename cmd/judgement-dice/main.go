package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/rin2yh/games/internal/judgementdice/game"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("ジャッジメントダイス")
	if err := ebiten.RunGame(game.New()); err != nil {
		log.Fatal(err)
	}
}
