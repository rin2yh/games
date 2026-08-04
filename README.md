# Games

Go + [Ebitengine](https://ebitengine.org/) で作ったゲーム集。WebAssembly にビルドして GitHub Pages で公開しています。

公開先: https://rin2yh.github.io/games/

| ゲーム | 公開先 | 詳細 |
| --- | --- | --- |
| Dinosaur Game | [/dinosaur-game/](https://rin2yh.github.io/games/dinosaur-game/) | [docs/dinosaur-game.md](docs/dinosaur-game.md) |
| Gopher Run | [/gopher-run/](https://rin2yh.github.io/games/gopher-run/) | [docs/gopher-run.md](docs/gopher-run.md) |
| ジャッジメントダイス | [/judgement-dice/](https://rin2yh.github.io/games/judgement-dice/) | [docs/judgement-dice.md](docs/judgement-dice.md) |

## 構成

ゲームごとに `cmd/<game>` の main パッケージと `internal/<game>` のロジックがあり、モジュールと依存はリポジトリ全体で 1 つに揃えています。

```
cmd/
  dinosaur-game/      Ebitengine フロントエンド（デスクトップと WASM）
  dinosaur-firmware/  koebiten フロントエンド（マイコン実機向け、TinyGo でビルド）
  gopher-run/         エントリポイントとアセットの埋め込み
  judgement-dice/     エントリポイント
internal/
  dinosaur/game/      エンジン非依存のゲームロジックと描画
  dinosaur/beep/      効果音の PCM 生成
  gopherrun/          エンティティ・シーン・ワールド
  judgementdice/game/ Dice / Duel / Jenga / Effects
targets/              TinyGo のカスタムターゲット定義（ボードごとに 1 ファイル）
docs/                 ゲームごとの詳細
web/                  Pages に公開するツリーそのもの
  index.html          ゲーム一覧
  wasm_exec.js        Go が生成する共通のグルー（ビルド時に生成）
  <game>/index.html   ゲームごとのページ
  <game>/game.wasm    ゲームごとの WASM（ビルド時に生成）
```

`web/` の中身がそのまま公開されるツリーなので、URL は `web/` のディレクトリ名で決まります。ゲームを増やすときは `cmd/<name>` と `web/<name>/index.html` を足し、`mise.toml` と 2 つのワークフローの `GAMES` 相当のリストに名前を加えます。

## 開発

ツールバージョンは [mise](https://mise.jdx.dev/) で管理しています。

```sh
mise install
```

```sh
mise run build-web      # 全ゲームを web/<game>/game.wasm にビルド
mise run serve          # build-web した web/ を http://localhost:8000 で配信
GAME=gopher-run mise run serve-game
                        # 1 ゲームだけ wasmserve で配信（リロードで再ビルド）
mise run build          # デスクトップ向けバイナリを bin/ にビルド
mise run test           # go test ./...
mise run lint           # go vet ./...
mise run fmt            # gofmt -l -w .
mise run flash-firmware # zero-kb02 に dinosaur のファームウェアを書き込み
mise run clean          # ビルド成果物を削除
```

タスクの一覧と説明は `mise tasks` で確認できます。

デスクトップ実行には音声デバイスが必要です（Ebitengine は音声の初期化失敗をゲームのエラーとして返すため、鳴らせない環境では起動できません）。Linux では Ebitengine のビルド要件に加えて ALSA のヘッダが必要です。

```sh
sudo apt install libasound2-dev
```

Web (WASM) ビルドとファームウェアビルドには不要です。

## デプロイ

`main` への push で `.github/workflows/deploy-pages.yml` が全ゲームの WASM をビルドし、`web/` を GitHub Pages に公開します（手動実行も可）。初回のみリポジトリの Settings → Pages → Build and deployment の Source を「GitHub Actions」にしてください。

## ライセンス

[MIT License](LICENSE) です。

依存ライブラリもいずれも許容的なライセンスで、再配布時はそれぞれの著作権表示を含めてください。

| 依存 | ライセンス |
| --- | --- |
| [hajimehoshi/ebiten](https://github.com/hajimehoshi/ebiten) | Apache-2.0 |
| [hajimehoshi/bitmapfont](https://github.com/hajimehoshi/bitmapfont) | Apache-2.0 |
| [sago35/koebiten](https://github.com/sago35/koebiten) | MIT |
| [tinygo.org/x/drivers](https://github.com/tinygo-org/drivers) | BSD-3-Clause |
| [golang.org/x/image](https://pkg.go.dev/golang.org/x/image) | BSD-3-Clause |
