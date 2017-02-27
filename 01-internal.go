package quitter

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/proxy"
)

var apipath string
var ProxyDialer proxy.Dialer
var ProxyString string
var err error
var apigun = &http.Client{
	CheckRedirect: redirectPolicyFunc,
	Transport:     Transport,
}

var Transport = &http.Transport{}

func init() {

}

var enableSOCKS bool
var enableInvalidTLS bool

func EnableSOCKS(socksurl string) {
	if socksurl == "" {
		return
	}

	// tor users can use SOCKS=tor
	if strings.ToUpper(socksurl) == "TOR" {
		socksurl = "socks5://127.0.0.1:9050"
		// i use SOCKS=true
	} else if strings.ToUpper(socksurl) == "TRUE" {
		socksurl = "socks5://127.0.0.1:1080"
	}

	u, err := url.Parse(socksurl)
	if err != nil {
		log.Fatal("Error parsing SOCKS proxy URL:", socksurl, ".", err)
	}

	ProxyDialer, err = proxy.FromURL(u, proxy.Direct)
	if err != nil {
		log.Fatal("Error setting SOCKS proxy.", err)
	}

	Transport.Dial = ProxyDialer.Dial
	apigun.Transport = Transport
	ProxyString = u.String()
	fmt.Fprintf(os.Stderr, "Using SOCKS proxy: %q\n", ProxyString)
}

/*
EnableInvalidTLS In the rare event you want to ignore SSL certificate checks
To ignore ssl verification
*/
func EnableInvalidTLS() {
	if Transport == nil {
		Transport = new(http.Transport)
	}
	if Transport.TLSClientConfig == nil {

		Transport.TLSClientConfig = new(tls.Config)
	}

	Transport.TLSClientConfig.InsecureSkipVerify = true
	fmt.Fprint(os.Stderr, "Skipping TLS Verification (unsafe)\n")
	return
}

// UserAgent to send
var UserAgent = "go-quitter/0.9"

func redirectPolicyFunc(req *http.Request, reqs []*http.Request) error {
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", UserAgent)
	return nil
}
