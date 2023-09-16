package pkce_test

import (
	"strings"
	"testing"

	"github.com/js10x/okta-token-vendor/pkce"
)

func Test_AuthURL(t *testing.T) {
	scenarios := []struct {
		issuer string
	}{
		{issuer: "nnn"},
		{issuer: ""},
		{issuer: " "},
		{issuer: "htt://eeen.malformed"},
		{issuer: "http://www.google.com/&*)@@($*%"},
	}

	for _, test := range scenarios {
		result := pkce.AuthURL(test.issuer)

		// Did not fail parsing
		if len(result) > 0 {
			lastFive := result[len(result)-5:]
			if lastFive != "authn" {
				t.Errorf("Did not get the expected result. Expected ['%v'] Result ['%v']", test.issuer, result)
			}
		}
	}
}

func Test_OAuth2URL(t *testing.T) {
	scenarios := []struct {
		issuer      string
		endpointUri string
	}{
		{issuer: "nnn", endpointUri: "/m"},
		{issuer: "", endpointUri: "token"},
		{issuer: " ", endpointUri: " "},
		{issuer: "htt://eeen.malformed", endpointUri: "/m"},
		{issuer: "http://www.google.com/&*)@@($*%", endpointUri: "/m"},

		{issuer: "http://host.com/oauth2/issuer1", endpointUri: "token"},
		{issuer: "http://host.com/oauth2/issuer2", endpointUri: "authorize"},
	}

	for _, test := range scenarios {
		result := pkce.OAuth2URL(test.issuer, test.endpointUri)

		// Did not fail parsing
		if len(result) > 0 {
			if !strings.Contains(result, "oauth2") {
				t.Errorf("Did not get the expected result. Expected ['%v'] Result ['%v']", test.issuer, result)
			}
		}
	}
}

func Test_AuthCodeQuery(t *testing.T) {
	scenarios := []struct {
		clientID     string
		redirectUri  string
		sessionToken string
	}{
		{clientID: "", redirectUri: "", sessionToken: ""},
		{clientID: "*&))ine", redirectUri: "*&))ine", sessionToken: "*&))ine"},
		{clientID: "98998989", redirectUri: "0", sessionToken: "_+__#)$)"},
		{clientID: "cid", redirectUri: "callback", sessionToken: "token"},
	}

	for _, test := range scenarios {
		verifier, params := pkce.AuthCodeQuery(test.clientID, test.redirectUri, test.sessionToken)

		if len(verifier) <= 0 || len(params) <= 0 {
			t.Errorf("Failed to build query parameters for the authorization code query")
		}
	}
}

func Test_CodeChallenge_Returns_Valid_PKCE_String(t *testing.T) {

	scenarios := []struct {
		verifier  string
		challenge string
	}{
		{
			verifier:  "iOtLoUBoSr-q5SF-2LOCSF0nD0nK-q8ej_p6cSRAuV4zezS_30vNaA",
			challenge: "lDHeGJUIKkzy8poy_rbis3_noYOGPDWBbwA4fhmsWnU",
		},
		{
			verifier:  "XBQjAyRYjn1OhUU48Nq6usZaYFr_JcPQUGUzzYsYgFzbytRcN4qGwA",
			challenge: "ONkMuadz3XfoM83cEFaznbSxAP2vcBEF9zNi95aeDR8",
		},
	}

	for _, test := range scenarios {
		result := pkce.CodeChallenge(test.verifier)

		if result != test.challenge {
			t.Errorf("Did not get the expected result. Expected ['%v'] Result ['%v']", test.challenge, result)
		}
	}
}
