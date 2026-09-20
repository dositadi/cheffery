package httphelper

import (
	"net/http"
	"strings"

	"github.com/dositadi/cheffery/sharedkernel"
)

func GetAccessToken(r *http.Request) string {
	if header := r.Header.Get(sharedkernel.Authorization.String()); header != "" {
		if payload := strings.Split(header, " "); len(payload) == 2 && payload[0] == "Bearer" {
			return payload[1]
		}
	}
	return ""
}

func GetRefreshToken(r *http.Request) string {
	if header := r.Header.Get(sharedkernel.RefreshCustom.String()); header != "" {
		if payload := strings.Split(header, " "); len(payload) == 2 && payload[0] == "Bearer" {
			return payload[1]
		}
	}
	return ""
}
