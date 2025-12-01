//go:build cgo

package ntgcalls

type MediaDevices struct {
	Microphone, Speaker, Camera, Screen []DeviceInfo
}
