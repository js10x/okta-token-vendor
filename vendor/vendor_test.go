package vendor_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/js10x/okta-token-vendor/vendor"
)

type mockHttpClient struct {
	doStub func(req *http.Request) (*http.Response, error)
}

func (mc *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return mc.doStub(req)
}

var vendingMachine = func() (*vendor.TokenVendor, *mockHttpClient) {

	mc := &mockHttpClient{}
	ops := []vendor.Option{
		vendor.Client(mc),
		vendor.ClientID("CLIENT_ID"),
		vendor.Issuer("https://host.com/oauth2/randomString"),
		vendor.RedirectURI("http://host/login/callback"),
		vendor.OnTokenReceived(func(accessToken string) {
			log.Printf("received access token [%v]", accessToken)
		}),
	}
	return vendor.NewTokenVendor(ops), mc
}

func Test_GetSessionToken(t *testing.T) {

	scenarios := []struct {
		response interface{}
	}{
		{response: &vendor.OktaError{ErrorCode: "E0000022", ErrorSummary: "The endpoint does not support the provided HTTP method"}},
		{response: &vendor.SessionTokenResponse{Token: "token"}},
		{response: &vendor.SessionTokenResponse{Token: " "}},
		{response: &vendor.SessionTokenResponse{Status: "LOCKED_OUT", Token: ""}},
	}

	var buf bytes.Buffer
	oktv, mockClient := vendingMachine()

	for _, test := range scenarios {

		mockClient.doStub = func(req *http.Request) (*http.Response, error) {
			json.NewEncoder(&buf).Encode(test.response)
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&buf),
			}, nil
		}
		response, err := oktv.GetSessionToken("user", "pw")

		// Testing errors
		if err != nil && response != nil {
			t.Errorf("Failed to return a NIL response when error occured.")
		}

		// Testing valid responses
		if response != nil && err != nil {
			t.Errorf("Failed to return a NIL error when valid response was returned.")
		}
	}
}

func Test_GetAuthorizationCode(t *testing.T) {

	responseHeaders := make(map[string][]string)
	responseHeaders["location"] = []string{"https://test?code=test-code"}

	scenarios := []struct {
		statusCode int
		response   interface{}
		header     map[string][]string
	}{
		{
			statusCode: 200,
			header:     responseHeaders,
			response:   &vendor.AuthorizationCodeResponse{CodeVerifier: "verifier"},
		},
		{
			statusCode: 302,
			header:     responseHeaders,
			response:   &vendor.AuthorizationCodeResponse{Code: "code"},
		},
		{
			statusCode: 404,
			header:     nil,
			response:   &vendor.AuthorizationCodeResponse{Code: " "},
		},
	}

	var buf bytes.Buffer
	oktv, mockClient := vendingMachine()

	for _, test := range scenarios {

		mockClient.doStub = func(req *http.Request) (*http.Response, error) {
			json.NewEncoder(&buf).Encode(test.response)
			return &http.Response{
				StatusCode: test.statusCode,
				Body:       ioutil.NopCloser(&buf),
			}, nil
		}
		response, err := oktv.GetAuthorizationCode("session-token")

		// Testing errors
		if err != nil && response != nil {
			t.Errorf("Failed to return a NIL response when error occured.")
		}

		// Testing valid responses
		if response != nil && err != nil {
			t.Errorf("Failed to return a NIL error when valid response was returned.")
		}
	}
}

func Test_GetAccessToken(t *testing.T) {

	scenarios := []struct {
		response interface{}
	}{
		{response: &vendor.AccessTokenResponse{AccessToken: "token"}},
		{response: &vendor.AccessTokenResponse{AccessToken: " "}},
	}

	var buf bytes.Buffer
	oktv, mockClient := vendingMachine()

	for _, test := range scenarios {

		mockClient.doStub = func(req *http.Request) (*http.Response, error) {
			json.NewEncoder(&buf).Encode(test.response)
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&buf),
			}, nil
		}
		response, err := oktv.GetAccessToken("verifier", "auth-code")

		// Testing errors
		if err != nil && response != nil {
			t.Errorf("Failed to return a NIL response when error occured.")
		}

		// Testing valid responses
		if response != nil && err != nil {
			t.Errorf("Failed to return a NIL error when valid response was returned.")
		}
	}
}
