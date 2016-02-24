package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var goquitter = "go-quitter v0.0.2"
var username = os.Getenv("GNUSOCIALUSER")
var password = os.Getenv("GNUSOCIALPASS")
var gnusocialnode = os.Getenv("GNUSOCIALNODE")
var fast bool = false
var apipath string = "https://" + gnusocialnode + "/api/statuses/home_timeline.json"

type User struct {
	Name string `json:"name"`
}

var usage string

type Tweet struct {
	Id                   int64    `json:"id"`
	IdStr                string   `json:"id_str"`
	InReplyToScreenName  string   `json:"in_reply_to_screen_name"`
	InReplyToStatusID    int64    `json:"in_reply_to_status_id"`
	InReplyToStatusIdStr string   `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64    `json:"in_reply_to_user_id"`
	InReplyToUserIdStr   string   `json:"in_reply_to_user_id_str"`
	Lang                 string   `json:"lang"`
	Place                string   `json:"place"`
	PossiblySensitive    bool     `json:"possibly_sensitive"`
	RetweetCount         int      `json:"retweet_count"`
	Retweeted            bool     `json:"retweeted"`
	RetweetedStatus      *Tweet   `json:"retweeted_status"`
	Source               string   `json:"source"`
	Text                 string   `json:"text"`
	Truncated            bool     `json:"truncated"`
	User                 User     `json:"user"`
	WithheldCopyright    bool     `json:"withheld_copyright"`
	WithheldInCountries  []string `json:"withheld_in_countries"`
	WithheldScope        string   `json:"withheld_scope"`
}

type AuthSuccess struct {
	/* variables */
}
type AuthError struct {
	/* variables */
}

type Badrequest struct {
	error   string `json:"error"`
	request string `json:"request"`
}

func main() {
	usage = "Usage:\n\n\tgo-quitter read\t\t\tReads 20 new posts\n\tgo-quitter read fast\t\tReads 20 new posts (no delay)\n\tgo-quitter home\t\tReads 20 from your Home timeline.\n\nYou may set your GNUSOCIALNODE environmental variable to change nodes.\nFor example: `export GNUSOCIALNODE=gs.sdf.org` in your ~/.shrc or ~/.profile\n\nExplore!\n\n\tGNUSOCIALNODE=gnusocial.de go-quitter read\n\tGNUSOCIALNODE=quitter.es go-quitter read\n\tGNUSOCIALNODE=shitposter.club go-quitter read\n\tGNUSOCIALNODE=sealion.club go-quitter read\n\t(defaults node is gs.sdf.org)\n"

	if len(os.Args) < 2 {
		log.Println("go-quitter v0.0.2")
		log.Println("Copyright 2016 aerth@sdf.org")
		log.Fatalln(usage)
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

	// go-quitter home
	if os.Args[1] == "home" && len(os.Args) == 2 {
		readHome(false)
		os.Exit(0)
	}
	// go-quitter home fast
	if os.Args[1] == "home" && os.Args[2] == "fast" {
		readHome(true)
		os.Exit(0)
	}

	// go-quitter post "Testing form console line using go-quitter"
	if os.Args[1] == "post" && os.Args[2] != ""{
		postNew(os.Args[2])
		os.Exit(0)
	}

	log.Fatalln(usage)

}

// readNew shows 20 new messages. Defaults to a 2 second delay, but can be called with readNew(fast) for a quick dump.
func readNew(fast bool) {

	log.Println("node: " + gnusocialnode)
	res, err := http.Get("https://" + gnusocialnode + "/api/statuses/public_timeline.json")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	var tweets []Tweet
	_ = json.Unmarshal(body, &tweets)
	//if err != nil { log.Fatalln(err) }
	for i := range tweets {
		fmt.Printf("[" + tweets[i].User.Name + "] " + tweets[i].Text + "\n\n")
		if fast != true {
			time.Sleep(2000 * time.Millisecond)
		}
	}

}

// readNew shows 20 new messages. Defaults to a 2 second delay, but can be called with readNew(fast) for a quick dump.
func readHome(fast bool) {
	if username == "" || password == "" {
		log.Fatalln("Please set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to view home timeline.")
	}
	log.Println("node: " + gnusocialnode)

	apipath := "https://" + gnusocialnode + "/api/statuses/home_timeline.json"
	req, err := http.NewRequest("GET", apipath, nil)
	req.SetBasicAuth(username, password)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var apres Badrequest

	body, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(body, &apres)
	fmt.Println(apres.error)
	fmt.Println(apres.request)

	var tweets []Tweet
	_ = json.Unmarshal(body, &tweets)
	//if err != nil { log.Fatalln(err) } // This fails
	for i := range tweets {
		fmt.Printf("[" + tweets[i].User.Name + "] " + tweets[i].Text + "\n\n")
		if fast != true {
			time.Sleep(2000 * time.Millisecond)
		}
	}

}

func postNew(content string) {
	if username == "" || password == "" {
		log.Fatalln("Please set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}

	log.Println("posting on node: " + gnusocialnode)
	log.Println(content)

	apipath := "https://" + gnusocialnode + "/api/statuses/update.json"

	/*
		resp, err := resty.R().

				SetBasicAuth(username, password).
				SetBody(`{"status":"`+content+`"}`).
				//SetResult(AuthSuccess{}). // or S	etResult(&AuthSuccess{}).
				Post(apipath)

			fmt.Println(resp, err)

	*/
	//var jsonStr = []byte(`{"status": "`+content+`"}`)
	//req, err := http.NewRequest("POST", apipath, bytes.NewBuffer(jsonStr))
	req, err := http.NewRequest("POST", apipath, bytes.NewBuffer([]byte(`{"status": "testing from go-quitter command line.... its not working."}`)))
	req.Header.Add("Authorization", "Basic RmFuY3kgbWVldGluZyB5b3UgaGVyZSEg")
	req.SetBasicAuth(username, password)
	req.Header.Add("Content-Type", "[application/json; charset=utf-8")
	//	req.Header.Set("User-Agent", goquitter)
	log.Println(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var apres Badrequest
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	_ = json.Unmarshal(body, &apres)
	fmt.Println(apres.error)
	fmt.Println(apres.request)

}

func init() {
	if gnusocialnode == "" {
		gnusocialnode = "gs.sdf.org"
	}

}
