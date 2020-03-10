package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	tests := []struct {
		name                          string
		requestHeaders                string
		wantStatus                    int
		wantAccessControlAllowOrigin  string
		wantAccessControlAllowHeaders string
	}{
		{
			name:                          "正常系（Authorization ヘッダーのリクエストは許可する）",
			requestHeaders:                "Authorization",
			wantStatus:                    http.StatusOK,
			wantAccessControlAllowOrigin:  "*",
			wantAccessControlAllowHeaders: "Authorization",
		},
		{
			name:                          "異常系（許可していないヘッダーのリクエストは拒否する）",
			requestHeaders:                "X-Deny-Header",
			wantStatus:                    http.StatusForbidden,
			wantAccessControlAllowOrigin:  "",
			wantAccessControlAllowHeaders: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			ts := runTestServer()
			defer ts.Close()

			u, _ := url.Parse(ts.URL)
			r, err := http.DefaultClient.Do(&http.Request{
				Method: http.MethodOptions,
				URL:    u,
				Header: map[string][]string{
					"Origin":                         {"http://localhost"},
					"Access-Control-Request-Method":  {"POST"},
					"Access-Control-Request-Headers": {tt.requestHeaders},
				},
			})
			assert.NoError(err)
			assert.Equal(tt.wantStatus, r.StatusCode)
			assert.Equal(tt.wantAccessControlAllowOrigin, r.Header.Get("Access-Control-Allow-Origin"))
			assert.Equal(tt.wantAccessControlAllowHeaders, r.Header.Get("Access-Control-Allow-Headers"))
		})
	}
}

func runTestServer() *httptest.Server {
	router := mux.NewRouter().StrictSlash(true)
	router.
		Methods("POST").
		Path("/").
		Name("test").
		HandlerFunc(func(http.ResponseWriter, *http.Request) { return })

	handler := middleware.CORS(router)
	return httptest.NewServer(handler)
}
