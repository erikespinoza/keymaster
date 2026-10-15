package u2f

import (
	"net/url"
	"strings"
)

// verifyAppId ensures the hostname matches...at this moment we do NOT check port number
// Port number should also be checked but leaving that out for now.
func verifyAppId(baseURLStr string, AppIdStr string) (bool, error) {
	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		return false, err
	}
	baseURLHost, _, _ := strings.Cut(baseURL.Host, ":")
	if AppIdStr == baseURL.Host || AppIdStr == baseURLHost {
		return true, nil
	}
	// The base ID does not match... so we will now try to parse the appID
	AppId, err := url.Parse(AppIdStr)
	if err != nil {
		return false, err
	}
	appIDHost, _, _ := strings.Cut(AppId.Host, ":")
	if appIDHost == baseURLHost {
		return true, nil
	}
	return false, nil
}
