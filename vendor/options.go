package vendor

import (
	"net/http"
	"os"
	"strings"
)

type TokenReceivedHandler func(string)

type Option func(*Options)

type Options struct {
	ClientID        string
	Issuer          string
	RedirectURI     string
	Client          HttpClient
	OnTokenReceived TokenReceivedHandler
}

func GetDefaultOptions() Options {
	return Options{
		ClientID:    os.Getenv("CLIENT_ID"),
		Issuer:      os.Getenv("ISSUER"),
		RedirectURI: os.Getenv("REDIRECT_URI"),
		Client: &http.Client{
			// Instructs the client not to follow a redirect, allowing us to
			// grab the token from the URL before the redirect occurs.
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			// Hook up a custom transport so that we can log each request.
			Transport: &LoggingRoundTripper{
				DefaultRoundTripper: http.DefaultTransport,
				Logger:              os.Stdout,
			},
		},
	}
}

func ClientID(cid string) Option {
	return func(o *Options) {
		if len(strings.TrimSpace(cid)) > 0 {
			o.ClientID = cid
		}
	}
}

func Issuer(iss string) Option {
	return func(o *Options) {
		if len(strings.TrimSpace(iss)) > 0 {
			o.Issuer = iss
		}
	}
}

func RedirectURI(uri string) Option {
	return func(o *Options) {
		if len(strings.TrimSpace(uri)) > 0 {
			o.RedirectURI = uri
		}
	}
}

func Client(c HttpClient) Option {
	return func(o *Options) { o.Client = c }
}

func OnTokenReceived(c TokenReceivedHandler) Option {
	return func(o *Options) { o.OnTokenReceived = c }
}
