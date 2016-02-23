package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"golang.org/x/crypto/ssh/terminal"
	//"github.com/gillesdemey/go-oauth"
	//"gopkg.in/resty.v0"
	//    "github.com/aerth/anaconda"
)

type User struct {
	Name string `json:"name"`
}
type Tweet struct {
	Id                   int64  `json:"id"`
	IdStr                string `json:"id_str"`
	InReplyToScreenName  string `json:"in_reply_to_screen_name"`
	InReplyToStatusID    int64  `json:"in_reply_to_status_id"`
	InReplyToStatusIdStr string `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64  `json:"in_reply_to_user_id"`
	InReplyToUserIdStr   string `json:"in_reply_to_user_id_str"`
	Lang                 string `json:"lang"`
	Place                string `json:"place"`
	PossiblySensitive    bool   `json:"possibly_sensitive"`
	RetweetCount         int    `json:"retweet_count"`
	Retweeted            bool   `json:"retweeted"`
	RetweetedStatus      *Tweet `json:"retweeted_status"`
	Source               string `json:"source"`
	Text                string   `json:"text"`
	Truncated           bool     `json:"truncated"`
	User                User     `json:"user"`
	WithheldCopyright   bool     `json:"withheld_copyright"`
	WithheldInCountries []string `json:"withheld_in_countries"`
	WithheldScope       string   `json:"withheld_scope"`
}

var fast bool = false

func main() {

oldState, err := terminal.MakeRaw(0)
if err != nil {
        panic(err)
}
defer terminal.Restore(0, oldState)

	log.Println("go-quitter v0.0.1")
	log.Println("Copyright 2016 aerth@sdf.org")
	if len(os.Args) < 2 {
		log.Fatalln("Usage:\n\n\tgo-quitter read\t\t\tReads 20 new posts\n\tgo-quitter read fast\t\tReads 20 new posts (no delay)\n\nYou may set your GNUSOCIALNODE environmental variable to change nodes.\nFor example: `export GNUSOCIALNODE=gs.sdf.org` in your ~/.shrc or ~/.profile\n\nExplore!\n\n\tGNUSOCIALNODE=gnusocial.de go-quitter read fast\n\tGNUSOCIALNODE=quitter.es go-quitter read fast\n\tGNUSOCIALNODE=shitposter.club go-quitter read fast\n\tGNUSOCIALNODE=sealion.club go-quitter read fast\n\t(defaults node is gs.sdf.org)\n")
	}

	// go-quitter read
	if os.Args[1] == "read" && len(os.Args) == 2 {
		readNew(false)
		os.Exit(0)
	}

	// go quitter read fast
	if os.Args[1] == "read" && os.Args[2] == "fast" {
		readNew(true)
		os.Exit(0)
	}

	log.Fatalln("Usage:\n\n\tgo-quitter read\t\t\tReads 20 new posts\n\tgo-quitter read fast\t\tReads 20 new posts (no delay)\n\nYou may set your GNUSOCIALNODE environmental variable to change nodes.\nFor example: `export GNUSOCIALNODE=gs.sdf.org` in your ~/.shrc or ~/.profile\n\nExplore!\n\n\tGNUSOCIALNODE=gnusocial.de go-quitter read\n\tGNUSOCIALNODE=quitter.es go-quitter read\n\tGNUSOCIALNODE=shitposter.club go-quitter read\n\tGNUSOCIALNODE=sealion.club go-quitter read\n\t(defaults node is gs.sdf.org)\n")

}

// readNew shows 20 new messages. Defaults to a 2 second delay, but can be called with readNew(fast) for a quick dump.
func readNew(fast bool) {
	//var Tweet Tweet
	var gnusocialnode string
	if os.Getenv("GNUSOCIALNODE") == "" {
		gnusocialnode = "gs.sdf.org"
	} else {
		gnusocialnode = os.Getenv("GNUSOCIALNODE")
	}
	log.Println("node: " + gnusocialnode)
	res, err := http.Get("https://" + gnusocialnode + "/api/statuses/public_timeline.json")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	var tweets []Tweet
	_ = json.Unmarshal(body, &tweets)

	for i := range tweets {
		fmt.Printf("[" + tweets[i].User.Name + "] " + tweets[i].Text + "\n\n")
		if fast != true {
			time.Sleep(2000 * time.Millisecond)
		}
	}

}

func init() {


}
