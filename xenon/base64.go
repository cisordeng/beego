package xenon

import (
	"encoding/base64"
	"strings"
)

func DecodeBase64(base64String string) ([]byte, error) {
	base64String = strings.TrimSpace(base64String)

	if strings.ContainsAny(base64String, "-_") {
		data, err := base64.RawURLEncoding.DecodeString(base64String)
		if err == nil {
			return data, nil
		}
		return base64.URLEncoding.DecodeString(base64String)
	}

	if len(base64String) % 4 != 0 {
		data, err := base64.RawStdEncoding.DecodeString(base64String)
		if err == nil {
			return data, nil
		}
		for len(base64String) % 4 != 0 {
			base64String += "="
		}
	}

	return base64.StdEncoding.DecodeString(base64String)
}
