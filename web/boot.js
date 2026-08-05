// 各ゲーム共通の WASM 起動処理。index.html からは
// <script src="../wasm_exec.js"></script><script src="../boot.js"></script>
// の順に読む。カレントディレクトリの game.wasm を読むので、置き場所は
// web/<game>/index.html が想定。
(() => {
  // Ebitengine の js バックエンドは keydown / keyup を canvas に登録する
  // (internal/ui/ui_js.go)。canvas は tabindex 付きで生成されるものの、
  // 誰もフォーカスを当てないので、フォーカスが canvas に無い間はキー入力が
  // どこにも届かない。canvas をクリックすれば Ebitengine 側の mousedown
  // ハンドラが focus() を呼ぶが、それまでは何分待っても効かないままで、
  // キー入力しか受け付けないゲームは操作する手段が無い。
  const focusCanvas = () => {
    const canvas = document.querySelector("canvas");
    if (!canvas) {
      return false;
    }
    canvas.focus();
    return true;
  };

  // canvas は go.run() の中で body に追加されるので、出てくるのを待つ。
  const observer = new MutationObserver(() => {
    if (focusCanvas()) {
      observer.disconnect();
    }
  });
  observer.observe(document.body, { childList: true });

  // canvas の外（本文の余白）をクリックするとフォーカスが外れ、以降キーが
  // 効かなくなる。戻ってきたときに拾い直す。リンクは押した先に用があるので
  // そのままにする。
  document.addEventListener("pointerdown", (event) => {
    if (event.target instanceof Element && event.target.closest("a")) {
      return;
    }
    focusCanvas();
  });

  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("game.wasm"), go.importObject).then(
    (result) => {
      go.run(result.instance);
    },
  );
})();
