package quitter

import "net/http"

var apipath string = "https://localhost/api/statuses/home_timeline.json"

var (
	goquitter = "go-quitter v0.0.8"
)

// Set User Agent
var tr = &http.Transport{
	DisableCompression: true,
}
var apigun = &http.Client{
	CheckRedirect: redirectPolicyFunc,
	Transport:     tr,
}

func redirectPolicyFunc(req *http.Request, reqs []*http.Request) error {
	req.Header.Add("Content-Type", "[application/json; charset=utf-8")
	req.Header.Set("User-Agent", goquitter)
	return nil
}
