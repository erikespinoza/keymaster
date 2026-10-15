//go:build !android && cgo && !nou2f && !nowebauthn
// +build !android,cgo,!nou2f,!nowebauthn

package twofa

import (
	"net/http"
	"os"

	"github.com/Cloud-Foundations/golib/pkg/log"
	"github.com/Cloud-Foundations/keymaster/lib/client/twofa/u2f"
	"github.com/flynn/u2f/u2fhid"
	"github.com/marshallbrekka/go-u2fhost"
)

// tryFidoMFA performs a fido authentication step
// If there are no devices connected it will return false, nil
// if there are fido devices connected it will return
// true, nil on successul MFA and false, error on failure to
// perform the Fido authentication
func tryFidoMFA(
	baseURL string,
	client *http.Client,
	userAgentString string,
	logger log.DebugLogger,
) (bool, error) {
	// Linux support for the new library is not quite correct
	// so for now we keep using the old library (pure u2f)
	// for linux cli as default. Windows 10 and MacOS have been
	// tested successfully.
	// The env variable allows us to swap what library is used by
	// default
	useWebAuthh := true
	if os.Getenv("KEYMASTER_USEALTU2FLIB") != "" {
		useWebAuthh = !useWebAuthh
	}
	var err error
	if !useWebAuthh {
		devices, err := u2fhid.Devices()
		if err != nil {
			logger.Printf("could not open hid devices err=%s", err)
			return false, err
		}
		if len(devices) < 1 {
			logger.Debugf(2, "No Fido devices found")
			return false, nil
		}
		err = u2f.DoU2FAuthenticate(
			client, baseURL, userAgentString, logger)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	devices := u2fhost.Devices()
	if devices == nil || len(devices) < 1 {
		logger.Debugf(2, "No Fido devices found")
		return false, nil
	}
	err = u2f.WithDevicesDoWebAuthnAuthenticate(devices,
		client, baseURL, userAgentString, logger)
	if err != nil {
		logger.Debugf(1, "Error doing hid webathentication err=%s", err)
		return false, err
	}
	return true, nil
}
