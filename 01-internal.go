package quitter

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/proxy"
)

// internal functions and vars not exported

var apipath string

//var apipath = "https://localhost/api/statuses/home_timeline.json"
var proxyDialer proxy.Dialer
var socks = os.Getenv("SOCKS")
var err error

var (
	goquitter = "go-quitter v0.0.9"
)

// // Set User Agent
var apigun = &http.Client{
	CheckRedirect: redirectPolicyFunc,
	Transport:     tr,
}

var unsafe = os.Getenv("GNUSOCIALUNSAFE")
var tr = &http.Transport{}

func init() {
	Init()
}

// Init can be called after something like
// os.Setenv("SOCKS", "socks5://localhost:9050"),
// otherwise we don't recognize the env var because we have
// already initialized (before injecting os.Setenv).
func Init() {
	/*

		 This Init() is just for those who want either:
		 	To use a proxy
			To ignore ssl verification

			And:
			Want to set os.Getenv("SOCKS") in a program
			they are writing, instead of typing it on the command line.
	*/
	socks = os.Getenv("SOCKS")

	if socks == "" {
		socks = os.Getenv("PROXY")
	}

	if socks == "" && unsafe == "" {
		return
	}

	/*
		   Socks proxy
		   First check it is valid. Needs to be valid.
			 		SOCKS=socks5://127.0.0.1:9050 program_name


	*/
	if os.Getenv("TOR") != "" || strings.ToUpper(socks) == "TOR" {
		socks = "socks5://127.0.0.1:9050"
	} else if strings.ToUpper(socks) == "TRUE" {
		socks = "socks5://127.0.0.1:1080"
	}
	if socks != "" {
		u, err := url.Parse(socks)
		if err != nil {
			log.Fatal("Error parsing SOCKS proxy URL:", socks, ".", err)
		}
		fmt.Fprintf(os.Stderr, "Using SOCKS proxy: %q\n", u.String())
		proxyDialer, err = proxy.FromURL(u, proxy.Direct)
		if err != nil {
			log.Fatal("Error setting SOCKS proxy.", err)
		}
		tr.Dial = proxyDialer.Dial

	}

	/*
	  Unsafe SSL
	  In the rare event you want to ignore SSL certificate checks
	*/
	if unsafe != "" {
		tr.TLSClientConfig.InsecureSkipVerify = true
	}

}

func redirectPolicyFunc(req *http.Request, reqs []*http.Request) error {
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", goquitter)
	return nil
}
