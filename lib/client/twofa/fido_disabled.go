//go:build android || !cgo || nou2f || nowebauthn
// +build android !cgo nou2f nowebauthn

package twofa

import (
	"net/http"

	"github.com/Cloud-Foundations/golib/pkg/log"
)

// tryFidoMFA is stubbed for android/non-cgo builds where hardware FIDO/U2F is not available
func tryFidoMFA(
	baseURL string,
	client *http.Client,
	userAgentString string,
	logger log.DebugLogger,
) (bool, error) {
	logger.Debugf(1, "FIDO/U2F authentication is not supported on this platform/build")
	return false, nil
}
