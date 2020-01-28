package middleware_test

import (
	"bytes"
	"github.com/erichnascimento/rocket/middleware"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLogger_EnsureNextIsCalled(t *testing.T) {
	logger := middleware.NewLogger()

	var nextCalled bool
	next := func(rw http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}
	handler := logger.Mount(next)

	body := strings.NewReader("hello")
	req, err := http.NewRequest(http.MethodGet, "/profile", body)
	if err != nil {
		t.Fatal(err)
	}
	rw := httptest.NewRecorder()
	handler(rw, req)
	if !nextCalled {
		t.Fatal("next handler should be invoked")
	}
}

func TestLogger_WithDefaultLogger_ShouldLogUsingDefaultOutputFormat(t *testing.T) {
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())

	var logbuf bytes.Buffer
	log.SetOutput(&logbuf)
	log.SetFlags(0)

	next := func(rw http.ResponseWriter, r *http.Request) {
		time.After(time.Millisecond * 10)
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("hello"))
	}

	logger := middleware.NewLogger()
	handler := logger.Mount(next)

	body := strings.NewReader("hello")
	req, err := http.NewRequest(http.MethodGet, "/profile", body)
	if err != nil {
		t.Fatal(err)
	}
	rw := httptest.NewRecorder()
	handler(rw, req)

	if expected, given := "GET  200 0ms - 5 B\n", logbuf.String(); expected != given {
		t.Errorf(`log output mismatch. Expected:  "%s", given: "%s"`, expected, given)
	}
}

func TestLogger_NewCustomLogger(t *testing.T) {
	var givenRw http.ResponseWriter = nil
	var givenReq *http.Request = nil
	var givenDuration time.Duration = -1

	logger, err := middleware.NewCustomLogger(func(w http.ResponseWriter, r *http.Request, d time.Duration) {
		givenRw, givenReq, givenDuration = w, r, d
	})
	if err != nil {
		t.Fatal(err)
	}
	handler := logger.Mount(createNextHandler(time.Millisecond, http.StatusOK, "pong"))
	req := createRequest(http.MethodGet, "/profile", "ping")
	rw := httptest.NewRecorder()
	handler(rw, req)

	if givenRw == nil {
		t.Errorf(`expected an http response instance, nil given`)
	}

	if givenReq == nil {
		t.Errorf(`expected an http request instance, nil given`)
	}

	if givenDuration == 0 {
		t.Errorf(`expected a request duration, %d given`, givenDuration)
	}
}

func createNextHandler(sleepDuration time.Duration, statusCode int, bodyContent string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		time.After(sleepDuration)
		rw.WriteHeader(statusCode)
		_, _ = rw.Write([]byte(bodyContent))
	}
}

func createRequest(httpMethod, url string, bodyContent string) *http.Request {
	body := strings.NewReader(bodyContent)
	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		panic(err)
	}
	return req
}
