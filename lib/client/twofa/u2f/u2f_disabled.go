//go:build android || !cgo || nou2f || nowebauthn
// +build android !cgo nou2f nowebauthn

package u2f

import (
	"errors"
	"net/http"

	"github.com/Cloud-Foundations/golib/pkg/log"
)

var ErrU2FNotSupported = errors.New("u2f/webauthn is not supported on this platform/build")

// CheckU2FDevices checks the U2F devices (stubbed for non-cgo/android builds)
func CheckU2FDevices(logger log.DebugLogger) error {
	logger.Debugf(1, "CheckU2FDevices called but U2F/WebAuthn is not supported on this platform/build")
	return ErrU2FNotSupported
}

// DoU2FAuthenticate does U2F authentication (stubbed for non-cgo/android builds)
func DoU2FAuthenticate(
	client *http.Client,
	baseURL string,
	userAgentString string,
	logger log.DebugLogger) error {
	return ErrU2FNotSupported
}
