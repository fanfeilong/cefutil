/**
 * xgoef frontend bridge (Tauri-like invoke / listen / emit).
 *
 * Phase 1: posts JSON to the host via a placeholder channel.
 * When CEF message-router bindings land, replace __xgoefPost with
 * the real transport (cefQuery or process message).
 */
(function (global) {
  const pending = new Map();
  let nextId = 1;
  const listeners = new Map();

  function post(msg) {
    if (typeof global.__xgoefPost === "function") {
      global.__xgoefPost(JSON.stringify(msg));
      return;
    }
    // Dev fallback: echo reject so pages can detect missing host.
    queueMicrotask(function () {
      dispatch({
        type: "response",
        id: msg.id,
        ok: false,
        error: "xgoef host bridge not connected",
      });
    });
  }

  function dispatch(msg) {
    if (msg.type === "response" && pending.has(msg.id)) {
      const { resolve, reject } = pending.get(msg.id);
      pending.delete(msg.id);
      if (msg.ok) resolve(msg.payload);
      else reject(new Error(msg.error || "invoke failed"));
      return;
    }
    if (msg.type === "event") {
      const list = listeners.get(msg.event) || [];
      list.forEach(function (fn) {
        try {
          fn(msg.payload);
        } catch (_) {}
      });
    }
  }

  global.__xgoefDispatch = dispatch;

  global.__XGOEF__ = {
    invoke: function (cmd, args) {
      const id = nextId++;
      return new Promise(function (resolve, reject) {
        pending.set(id, { resolve, reject });
        post({ type: "invoke", id: id, cmd: cmd, args: args || {} });
      });
    },
    listen: function (event, handler) {
      if (!listeners.has(event)) listeners.set(event, []);
      listeners.get(event).push(handler);
      return function () {
        const list = listeners.get(event) || [];
        listeners.set(
          event,
          list.filter(function (h) {
            return h !== handler;
          })
        );
      };
    },
    emit: function (event, payload) {
      post({ type: "event", event: event, payload: payload });
    },
  };
})(typeof window !== "undefined" ? window : globalThis);
