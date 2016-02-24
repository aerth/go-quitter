package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var goquitter = "go-quitter v0.0.4"
var username = os.Getenv("GNUSOCIALUSER")
var password = os.Getenv("GNUSOCIALPASS")
var gnusocialnode = os.Getenv("GNUSOCIALNODE")
var fast bool = false
var apipath string = "https://" + gnusocialnode + "/api/statuses/home_timeline.json"

type User struct {
	Name       string `json:"name"`
	Screenname string `json:"screen_name"`
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

type Badrequest struct {
	terror  string `json:"error"`
	request string `json:"request"`
}

func main() {
	usage = "\t" + `go-quitter v0.0.4	Copyright 2016 aerth@sdf.org
Usage:

	go-quitter read			Reads 20 new posts
	go-quitter read fast		Reads 20 new posts (no delay)
	go-quitter home			Your home timeline.
	go-quitter user username	Looks up "username" timeline
	go-quitter post ____ 		Posts to your node.
	go-quitter post 		Post mode.
	go-quitter mentions		Mentions your @name
	go-quitter search ___		Searches for ____
	go-quitter search		Search mode.

Set your GNUSOCIALNODE environmental variable to change nodes.
For example: "export GNUSOCIALNODE=gs.sdf.org" in your ~/.shrc or ~/.profile

Most commands require GNUSOCIALUSER and GNUSOCIALPASS variables set.

Did you know?	You can "go-quitter read | more"

`
	if len(os.Args) < 2 {
		log.Fatalln(usage)
	}

	// Set speed
	speed := false
	lastvar := len(os.Args)
	lastvar = (lastvar - 1)
	if os.Args[lastvar] == "fast" || os.Getenv("GNUSOCIALFAST") == "true" {
		speed = true
	}
	// go-quitter read
	if os.Args[1] == "read" {
		readNew(speed)
		os.Exit(0)
	}
	// go-quitter read
	if os.Args[1] == "mentions" {
		readMentions(speed)
		os.Exit(0)
	}
	// go-quitter home
	if os.Args[1] == "home" {
		readHome(speed)
		os.Exit(0)
	}

	// go-quitter post Testing from console line using go-quitter
	// Notice how there is no quotation marks.
	if os.Args[1] == "post" {
		content := ""
		if len(os.Args) > 1 {
			content = strings.Join(os.Args[2:], " ")
		}
		postNew(content)
		os.Exit(0)
	}

	// go-quitter search _____
	if os.Args[1] == "search" {
		searchstr := ""
		if len(os.Args) > 1 {
			searchstr = strings.Join(os.Args[2:], " ")
		}
		readSearch(searchstr, speed)
		os.Exit(0)
	}

	// go-quitter user aerth
	if os.Args[1] == "user" && os.Args[2] != "" {
		userlookup := os.Args[2]
		readUserposts(userlookup, speed)
		os.Exit(0)
	}

	// this happens if we invoke with somehing like "go-quitter test"
	log.Fatalln(usage)

}

// readNew shows 20 new messages. Defaults to a 2 second delay, but can be called with readNew(fast) for a quick dump.
func readNew(fast bool) {
	fmt.Println("node: " + gnusocialnode)
	res, err := http.Get("https://" + gnusocialnode + "/api/statuses/public_timeline.json")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	var tweets []Tweet
	_ = json.Unmarshal(body, &tweets)
	for i := range tweets {
		if tweets[i].User.Screenname == tweets[i].User.Name {
			fmt.Printf("[@" + tweets[i].User.Screenname + "] " + tweets[i].Text + "\n\n")
		} else {
			fmt.Printf("@" + tweets[i].User.Screenname + " [" + tweets[i].User.Name + "] " + tweets[i].Text + "\n\n")
		}
		if fast != true {
			time.Sleep(2000 * time.Millisecond)
		}
	}

}

// readMentions shows 20 newest mentions of your username. Defaults to a 2 second delay, but can be called with readNew(fast) for a quick dump.
func readMentions(fast bool) {
	if username == "" || password == "" {
		log.Fatalln("Please set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	fmt.Println("node: " + gnusocialnode)
	apipath := "https://" + gnusocialnode + "/api/statuses/mentions.json"
	req, err := http.NewRequest("GET", apipath, nil)
	req.Header.Set("User-Agent", goquitter)
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
	fmt.Println(apres.terror)
	fmt.Println(apres.request)
	var tweets []Tweet
	_ = json.Unmarshal(body, &tweets)
	for i := range tweets {
		if tweets[i].User.Screenname == tweets[i].User.Name {
			fmt.Printf("[@" + tweets[i].User.Screenname + "] " + tweets[i].Text + "\n\n")
		} else {
			fmt.Printf("@" + tweets[i].User.Screenname + " [" + tweets[i].User.Name + "] " + tweets[i].Text + "\n\n")
		}
		if fast != true {
			time.Sleep(2000 * time.Millisecond)
		}
	}
}

// readHome shows 20 from home timeline. Defaults to a 2 second delay, but can be called with readHome(fast) for a quick dump.
func readHome(fast bool) {
	if username == "" || password == "" {
		log.Fatalln("Please set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to view home timeline.")
	}
	fmt.Println("node: " + gnusocialnode)
	apipath := "https://" + gnusocialnode + "/api/statuses/home_timeline.json"
	req, err := http.NewRequest("GET", apipath, nil)
	req.Header.Set("User-Agent", goquitter)
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
	fmt.Println(apres.terror)
	fmt.Println(apres.request)
	var tweets []Tweet
	_ = json.Unmarshal(body, &tweets)
	for i := range tweets {
		if tweets[i].User.Screenname == tweets[i].User.Name {
			fmt.Printf("[@" + tweets[i].User.Screenname + "] " + tweets[i].Text + "\n\n")
		} else {
			fmt.Printf("@" + tweets[i].User.Screenname + " [" + tweets[i].User.Name + "] " + tweets[i].Text + "\n\n")
		}
		if fast != true {
			time.Sleep(2000 * time.Millisecond)
		}
	}

}

func readSearch(searchstr string, fast bool) {
	if searchstr == "" {
		searchstr = getTypin()
	}
	if searchstr == "" {
		log.Fatalln("Blank search detected. Not searching.")
	}
	fmt.Println("searching " + searchstr + " @ " + gnusocialnode)
	v := url.Values{}
	v.Set("q", searchstr)
	searchstr = url.Values.Encode(v)
	apipath := "https://" + gnusocialnode + "/api/search.json?" + searchstr
	req, err := http.NewRequest("GET", apipath, nil)
	req.Header.Set("User-Agent", goquitter)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var tweets []Tweet
	_ = json.Unmarshal(body, &tweets)

	for i := range tweets {
		if tweets[i].User.Screenname == tweets[i].User.Name {
			fmt.Printf("[@" + tweets[i].User.Screenname + "] " + tweets[i].Text + "\n\n")
		} else {
			fmt.Printf("@" + tweets[i].User.Screenname + " [" + tweets[i].User.Name + "] " + tweets[i].Text + "\n\n")
		}
		if fast != true {
			time.Sleep(2000 * time.Millisecond)
		}
	}

}

func readUserposts(userlookup string, fast bool) {
	fmt.Println("user " + userlookup + " @ " + gnusocialnode)
	apipath := "https://" + gnusocialnode + "/api/statuses/user_timeline.json?screen_name=" + userlookup
	req, err := http.NewRequest("GET", apipath, nil)
	req.Header.Set("User-Agent", goquitter)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var tweets []Tweet
	_ = json.Unmarshal(body, &tweets)
	for i := range tweets {
		if tweets[i].User.Screenname == tweets[i].User.Name {
			fmt.Printf("[@" + tweets[i].User.Screenname + "] " + tweets[i].Text + "\n\n")
		} else {
			fmt.Printf("@" + tweets[i].User.Screenname + " [" + tweets[i].User.Name + "] " + tweets[i].Text + "\n\n")
		}
		if fast != true {
			time.Sleep(2000 * time.Millisecond)
		}
	}
}

func postNew(content string) {
	if username == "" || password == "" {
		log.Fatalln("Please set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	if content == "" {
		content = getTypin()
	}
	if content == "" {
		log.Fatalln("Blank status detected. Not posting.")
	}
	fmt.Println("Preview:\n\n[" + username + "] " + content)
	fmt.Println("\nType YES to publish!")
	if askForConfirmation() == false {
		os.Exit(0)
	}
	fmt.Println("posting on node: " + gnusocialnode)
	v := url.Values{}
	v.Set("status", content)
	content = url.Values.Encode(v)
	apipath := "https://" + gnusocialnode + "/api/statuses/update.json?" + content
	req, err := http.NewRequest("POST", apipath, nil)
	req.SetBasicAuth(username, password)
	req.Header.Set("HTTP_REFERER", "https://"+gnusocialnode+"/")
	req.Header.Add("Content-Type", "[application/json; charset=utf-8")
	req.Header.Set("User-Agent", goquitter)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	apres := Badrequest{}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("\nnode response:", resp.Status)
	_ = json.Unmarshal(body, &apres)
	if apres.terror != "" {
		fmt.Println(apres.terror)
	}
}

func init() {
	if gnusocialnode == "" {
		gnusocialnode = "gs.sdf.org"
	}

}
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	quitResponses := []string{"q", "Q", "exit", "quit"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else if containsString(quitResponses, response) {
		os.Exit(0)
		return false
	} else {
		fmt.Println("\nType YES to publish!")
		return askForConfirmation()
	}
}
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

func getTypin() string {
	fmt.Println("Press ENTER when you are finished typing.")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		return line
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return ""
}
