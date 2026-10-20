/**
 * xgoef frontend bridge (Tauri-like invoke / listen / emit).
 *
 * Transport: CEF MessageRouter `window.cefQuery` (injected by the host
 * renderer shim). Falls back to `__xgoefPost` when present.
 */
(function (global) {
  const pending = new Map();
  let nextId = 1;
  const listeners = new Map();

  function post(msg) {
    const body = JSON.stringify(msg);

    if (typeof global.cefQuery === "function") {
      global.cefQuery({
        request: body,
        persistent: false,
        onSuccess: function (response) {
          try {
            dispatch(JSON.parse(response));
          } catch (e) {
            dispatch({
              type: "response",
              id: msg.id,
              ok: false,
              error: "invalid host response: " + (e && e.message ? e.message : e),
            });
          }
        },
        onFailure: function (code, message) {
          dispatch({
            type: "response",
            id: msg.id,
            ok: false,
            error: message || "cefQuery failed (" + code + ")",
          });
        },
      });
      return;
    }

    if (typeof global.__xgoefPost === "function") {
      global.__xgoefPost(body);
      return;
    }

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
        pending.set(id, { resolve: resolve, reject: reject });
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
