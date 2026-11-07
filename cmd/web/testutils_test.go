package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	return &application{
		logger: slog.New(slog.DiscardHandler),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// we add cookie jar to the test server client, any response cookies
	// will now be stored and sent with subsequent requests when using
	// this client
	ts.Client().Jar = jar

	//	Disabling redirect-following by setting custom check redirect func.
	//	It's called whenever a 3xx response is received by the client, and by always
	//	returning a http.ErrUseLastResponse error it forces the client to immediatly
	//	return the recieved response
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}
