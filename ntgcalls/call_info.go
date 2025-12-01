//go:build cgo

package ntgcalls

type CallInfo struct {
	Playback, Capture StreamStatus
}
