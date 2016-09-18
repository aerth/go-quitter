package quitter

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/proxy"
)

var apipath string = "https://localhost/api/statuses/home_timeline.json"
var proxyDialer proxy.Dialer
var socks = os.Getenv("SOCKS")
var err error

// Set User Agent

var tr = &http.Transport{
	DisableCompression: true,
}

var (
	goquitter = "go-quitter v0.0.9"
)

// // Set User Agent
var apigun = &http.Client{
	CheckRedirect: redirectPolicyFunc,
	Transport:     tr,
}
var proxytr = &http.Transport{}
var proxygun = apigun

func init() {
	if socks != "" {
		urlsocks, err := url.Parse(socks)
		if err != nil {
			log.Fatal("Error parsing SOCKS proxy URL:", socks, ".", err)
		}
		proxyDialer, err = proxy.FromURL(urlsocks, proxy.Direct)
		if err != nil {
			log.Fatal("Error setting SOCKS proxy.", err)
		}

		// Set User Agent
		proxytr = &http.Transport{
			DisableCompression: true,
			Dial:               proxyDialer.Dial,
		}
		proxygun = &http.Client{
			CheckRedirect: redirectPolicyFunc,
			Transport:     proxytr,
		}
		apigun = &http.Client{
			CheckRedirect: redirectPolicyFunc,
			Transport:     proxytr,
		}
	}
}
func redirectPolicyFunc(req *http.Request, reqs []*http.Request) error {
	req.Header.Add("Content-Type", "[application/json; charset=utf-8")
	req.Header.Set("User-Agent", goquitter)
	return nil
}
