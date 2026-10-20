package cef

// QueryHandler is invoked on the CEF UI thread for each cefQuery request.
// windowLabel is the browser's bound label. Implementations must eventually
// call Respond (any thread).
type QueryHandler func(queryID int64, windowLabel, requestJSON string)

// WindowSpec describes a window to create at runtime.
type WindowSpec struct {
	Label  string
	Title  string
	URL    string
	Width  int
	Height int
}

// Bridge wires the C MessageRouter shim to Go IPC and window control.
type Bridge interface {
	SetQueryHandler(h QueryHandler)
	Respond(queryID int64, ok bool, resultUTF8 string)

	// Events (host → renderer). jsonUTF8 is a full envelope.
	EmitTo(windowLabel, jsonUTF8 string) error
	EmitAll(jsonUTF8 string) error

	// Window control (P1).
	SetPendingLabel(label string)
	CreateWindow(spec WindowSpec) error
	CloseWindow(label string) error
	FocusWindow(label string) error
}
