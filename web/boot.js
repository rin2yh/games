// 各ゲーム共通の WASM 起動処理。index.html からは
// <script src="../wasm_exec.js"></script><script src="../boot.js"></script>
// の順に読む。カレントディレクトリの game.wasm を読むので、置き場所は
// web/<game>/index.html が想定。
(() => {
  // Ebitengine の js バックエンドは keydown / keyup を canvas にだけ登録する
  // (internal/ui/ui_js.go)。起動時に canvas へフォーカスを当てるのも
  // Ebitengine 自身がやる（同 initOnMainThread）ので、開いた直後は効く。
  //
  // 問題はそのあと。index.html はゲームを画面中央に置くので、まわりの余白は
  // どのゲームでも広い。そこをクリックするとフォーカスは body に移り、canvas
  // を狙ってクリックし直すまでキーが一切効かなくなる。押しても無反応なので、
  // 原因がフォーカスだと気づける手がかりが画面に何も出ない。
  //
  // click で拾うのは、mousedown の既定の動作がフォーカスを移すのがこの後
  // だから。pointerdown や mousedown で focus() を呼んでも上書きされる。
  // リンクは押した先に用があるので、フォーカスを横取りしない。
  document.addEventListener("click", (event) => {
    if (event.target instanceof Element && event.target.closest("a")) {
      return;
    }
    document.querySelector("canvas")?.focus();
  });

  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("game.wasm"), go.importObject).then(
    (result) => {
      go.run(result.instance);
    },
  );
})();
