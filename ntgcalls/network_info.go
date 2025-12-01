//go:build cgo

package ntgcalls

type NetworkInfo struct {
	Kind  ConnectionKind
	State ConnectionState
}
