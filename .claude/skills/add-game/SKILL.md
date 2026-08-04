---
name: add-game
description: games リポジトリ (Go + Ebitengine を WebAssembly で GitHub Pages に公開) に新しいゲームを追加する手順。cmd/<name> と internal/<name> の置き方、web/<name>/index.html の書き方、一覧ページへの登録、ビルドと動作の検証までを扱う。「ゲームを追加したい」「新しいゲームを作って」「◯◯というゲームを実装して」のように、このリポジトリにゲームを増やす話が出たら必ず使うこと。ゲーム名を出さずに「テトリスみたいなやつ作りたい」「ぷよぷよ的なもの」のような言い方でも同じ。既存ゲームの移動やリネーム、web/ 配下の構成をいじる作業でも、どこが何を参照しているかを確かめるために読むとよい。
---

# ゲームを追加する

## このリポジトリの決まりごと

先にここを押さえると、後の手順が「なぜそう書くのか」まで分かる。

**`web/` が公開ツリーそのもの。** GitHub Pages には `web/` がまるごと上がるので、URL はそこのディレクトリ名で決まる。`web/tetris/` を作れば `https://rin2yh.github.io/games/tetris/` になる。

**ビルド対象は `web/*/` を列挙して決まる。** `mise.toml` にも 2 つのワークフローにもゲーム名は書かれていない。`web/<name>/` を作れば自動でビルド対象に入り、作らなければ入らない。逆に言うと **`web/<name>/index.html` を用意しないとそのゲームは存在しないのと同じ**になる。`cmd/` を列挙しない理由は、WASM にならない `cmd/dinosaur-firmware`（TinyGo 専用）まで拾ってしまうから。

**モジュールは全体で 1 つ。** `github.com/rin2yh/games` の下に全ゲームが同居する。ゲームのロジックは既定では `internal/<name>/...` に置き、外から import させたいものだけリポジトリ直下に出す（今は `dinosaur/` だけ）。

**WASM の出力名は `game.wasm` 固定、`wasm_exec.js` は `web/` に 1 つだけ。** どのゲームの index.html も `../wasm_exec.js` を参照する。ゲームごとに置くと同じファイルが複製されるだけなので置かない。

**起動処理は `web/boot.js` に 1 つだけ。** `game.wasm` の読み込みと、生成された canvas へのフォーカス当てが入っている。Ebitengine の js バックエンドはキーイベントを canvas に登録するので、フォーカスが無いとページを開いた直後はキーが一切効かない（issue #1）。index.html 側に起動コードを書き直すとこのフォーカス処理が抜けるので、`../boot.js` を読むだけにする。

## 手順

### 1. ゲームロジックを `internal/<name>/` に置く

パッケージ構成は中で自由に切ってよい。既存も揃っていない（`internal/gopherrun/` は entity/scene/world に分かれ、`internal/judgementdice/game/` は 1 パッケージ）。

ディレクトリ名はハイフンを使わない。Go のパッケージ名に合わせるため、`web/`・`cmd/` 側が `judgement-dice` でも `internal/` 側は `judgementdice` になる。

ロジックを Ebitengine から切り離せるなら切り離す価値はある。`dinosaur/game` は描画を `Display` インターフェース越しにしているので、Ebitengine とマイコン (koebiten) の 2 つのフロントエンドから同じロジックを動かせているし、GUI 無しでテストできる。ただしこれは必須ではない。

`dinosaur/` だけが `internal/` の外にいるのは、他所から import できる形にしてバージョンを名乗らせているから（`.tagpr` の `versionFile` がここを指す）。**新しいゲームは `internal/` に置く。** 外に出すのは、他のリポジトリから使わせたい、あるいは API を安定させて公開したいと決めたときだけで、そう決めた時点で互換性を保つ責任が発生する。

### 2. `cmd/<name>/main.go` を書く

`ebiten.RunGame` を呼ぶ薄い main にする。画像などのアセットは `cmd/<name>/assets/` に置いて `//go:embed` で埋め込む（`cmd/gopher-run/embed.go` が例）。WASM ではファイルシステムから読めないので、埋め込みが前提になる。

### 3. `web/<name>/index.html` を作る

`assets/index.html.template` をコピーして、`__TITLE__`・`__DESCRIPTION__`・`__WIDTH__`・`__HEIGHT__` を置き換える。幅と高さはゲームの `Layout` が返す論理サイズを入れる。

### 4. 一覧ページ `web/index.html` にカードを足す

`<ul>` の中に既存と同じ形で 1 つ追加する。ここに足さないとトップからたどり着けない。

```html
<li>
  <a class="card" href="<name>/">
    <div class="title">表示名</div>
    <p class="desc">1 行の説明。</p>
  </a>
</li>
```

### 5. `cmd/<name>/README.md` を書く

ゲームの説明・操作方法・構成を書く。ソースの隣にあることで、GitHub でディレクトリを開けばそのまま読める。公開先 URL は書かない（一覧ページとルート README が持っている）。

ルートの `README.md` にゲーム一覧は無いので、そちらを触る必要はない。

### 6. 新しい依存を足したなら `go mod tidy`

## 検証

ここまでやったら必ず動かして確かめる。手順を追っただけでは、`web/<name>/` の作り忘れやテンプレートの置換漏れは見つからない。

```sh
mise run build-web   # 自分のゲームが web/<name>/game.wasm を出力しているか
mise run lint
mise run test
gofmt -l .           # 何も出力されないこと
```

`build-web` の出力に自分のゲームが含まれていなければ、`web/<name>/` が無いか名前が `cmd/<name>` と食い違っている。

そのうえでブラウザで実際に遊べるか見る。

```sh
mise run serve   # http://localhost:8000
```

トップから自分のゲームに遷移できるか、操作できるか、遊べるかまで確認する。CI は WASM がビルドできることしか見ないので、実際に動くかどうかは人が触るしかない。

## よくある詰まり方

| 症状 | 原因 |
| --- | --- |
| ビルドされない | `web/<name>/` が無い。列挙元はここ |
| 404 になる | `web/<name>/index.html` が無い、または `cmd/<name>` と名前が違う |
| 真っ白なまま | index.html が `../boot.js` を読み込んでいない |
| `Go is not defined` | `../wasm_exec.js` を読み込んでいない。パスが `wasm_exec.js` になっていないか |
| 開いた直後だけキーが効かない | `../boot.js` を読み込んでいない。canvas にフォーカスが当たっていない |
| クリック / タップの位置がずれる | Ebitengine の js バックエンドはポインタ位置に `clientX` / `clientY` をそのまま使うので、index.html のように canvas を画面中央へ寄せていると canvas の位置だけずれる。`internal/webpointer` の `Offset` を引く（`internal/judgementdice/game/input.go` が例）|
