package vendor

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type LoggingRoundTripper struct {
	DefaultRoundTripper http.RoundTripper
	Logger              io.Writer
}

func (l LoggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	fmt.Fprintf(l.Logger, "\n\n[%s] %s %s\n\n", time.Now().Format(time.ANSIC), r.Method, r.URL.String())
	return l.DefaultRoundTripper.RoundTrip(r)
}
