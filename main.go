package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/js10x/okta-token-vendor/vendor"
)

func main() {

	var username, password, cid, iss, callback, out string
	var validConfig bool = false

	flag.StringVar(&username, "user", "The username associated with your Okta application.", "abc")
	flag.StringVar(&password, "pw", "The password associated with your Okta application.", "abc")
	flag.StringVar(&cid, "cid", "", "The client ID configured for your Okta application.")
	flag.StringVar(&iss, "iss", "", "The ISSUER configured for your Okta application.")
	flag.StringVar(&callback, "callback", "", "One of the configured REDIRECT URIs configured in your Okta application.")
	flag.StringVar(&out, "o", "", "Print the access token to the provided file.")
	flag.Parse()

	ops := []vendor.Option{
		vendor.ClientID(cid),
		vendor.Issuer(iss),
		vendor.RedirectURI(callback),
		vendor.OnTokenReceived(func(accessToken string) {
			if len(strings.TrimSpace(out)) <= 0 {
				return
			}
			// Write the access token to the provided file, if the user asked for it.
			file, err := os.Create(out)
			if err != nil {
				fmt.Printf("Error occurred when creating the output file provided: %v\n", err)
			} else {
				file.WriteString(accessToken)
			}
			file.Close()
		}),
	}
	oktv := vendor.NewTokenVendor(ops)

	switch {

	// Validate User ID and PW
	case len(strings.TrimSpace(username)) <= 0, len(strings.TrimSpace(password)) <= 0:
		fmt.Fprintf(os.Stderr, "You must specify both your username and password\n")

	// Validate Client ID
	case len(strings.TrimSpace(oktv.Ops.ClientID)) <= 0:
		fmt.Fprintf(os.Stderr, "You must specify a CLIENT ID\n")

	// Validate Issuer
	case len(strings.TrimSpace(oktv.Ops.Issuer)) <= 0:
		fmt.Fprintf(os.Stderr, "You must specify an ISSUER\n")

	// Validate Redirect URI
	case len(strings.TrimSpace(oktv.Ops.RedirectURI)) <= 0:
		fmt.Fprintf(os.Stderr, "You must specify a Redirect URI\n")

	default:
		validConfig = true
	}

	if !validConfig {
		os.Exit(0)
	}
	fmt.Fprintf(os.Stdout, "Configuration Accepted => Let's go get you a token.\n")

	// 1.) Get the session token
	sessionToken, err := oktv.GetSessionToken(username, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred when fetching the SESSION TOKEN: %v\n", err)
		os.Exit(0)
	}

	// 2.) Get the authorization code using the session token
	authCode, err := oktv.GetAuthorizationCode(sessionToken.Token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred when fetching the AUTHORIZATION TOKEN: %v\n", err)
		os.Exit(0)
	}

	// 3.) Get the access token using the authorization code
	accessToken, err := oktv.GetAccessToken(authCode.CodeVerifier, authCode.Code)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred when fetching the ACCESS TOKEN: %v\n", err)
		os.Exit(0)
	}
	fmt.Println(accessToken.ToString())
}
