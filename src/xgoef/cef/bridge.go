package cef

// QueryHandler is invoked on the CEF UI thread for each cefQuery request.
// Implementations must eventually call Respond (any thread).
type QueryHandler func(queryID int64, requestJSON string)

// Bridge wires the C MessageRouter shim to Go IPC.
type Bridge interface {
	SetQueryHandler(h QueryHandler)
	Respond(queryID int64, ok bool, resultUTF8 string)
}
