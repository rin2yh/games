# Games

Go + [Ebitengine](https://ebitengine.org/) のゲーム集。

# 開発

```sh
mise install
mise run serve   # http://localhost:8000 で全ゲームを配信
```

タスクの一覧は `mise tasks` で確認できます。

ゲームごとに `cmd/<game>` と `internal/<game>` があり、モジュールは全体で 1 つです。
`dinosaur/` だけは他所から import できるよう `internal/` の外にあり、`.tagpr` がそのバージョンを見ています。
ゲームの説明と構成は `cmd/<game>/README.md` にあります。
`web/` が公開されるツリーそのもので、URL はそこのディレクトリ名で決まります。
ビルド対象は `web/` の中身から決まるので、ゲームを増やすときに触るのは `cmd/<name>` と `web/<name>/index.html` だけです。

## ライセンス

[MIT License](LICENSE)
