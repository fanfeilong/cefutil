# Events（P1）

Host ↔ renderer events ride the same MessageRouter transport as invoke, plus
`ExecuteJavaScript` for host-initiated delivery.

## Envelope

```json
{ "type": "event", "event": "tick", "payload": { "n": 1 } }
```

## Directions

| Direction | Path |
| --- | --- |
| JS → host | `cefQuery` with `type:"event"` → `ipc.Bus.NotifyLocal` |
| Host → one window | `App.EmitTo(label, name, payload)` → `xgoef_ipc_emit_to` → `__xgoefDispatch` |
| Host → all | `App.Emit(name, payload)` → `xgoef_ipc_emit_all` |

Frontend: `__XGOEF__.listen(event, handler)` / `__XGOEF__.emit(event, payload)`.

## Window labels

Each browser is bound to a label (`--xgoef-label=` for the first window;
`CreateWindow` / pending label for later ones). Invoke ACL and `EmitTo` use that
label.
