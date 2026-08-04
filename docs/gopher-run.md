# Gopher Run

Go と Ebitengine で作られた 2D エンドレスランナーゲーム。Web ブラウザのみ対応。

公開先: https://rin2yh.github.io/games/gopher-run/

## 操作方法

- Space / クリック / タッチ: ジャンプ
- 長押し: より高くジャンプ
- ↓ / 2本指タッチ: もぐる

## 構成

```
cmd/gopher-run/         エントリポイントとアセットの埋め込み
  main.go               スプライトの切り出しとゲームの起動
  embed.go              assets/ の PNG を埋め込む
  assets/               スプライト
internal/gopherrun/
  game/                 ゲーム本体とシーン遷移
  game/scene/           タイトル・プレイ中・ゲームオーバー・スコア
  entity/player/        プレイヤー
  entity/enemy/         鷲とヘビ
  entity/particle/      パーティクル
  entity/popup/         ポップアップ
  input/                入力
  world/                地形
```

## 実行

```sh
GAME=gopher-run mise run serve-game
```
