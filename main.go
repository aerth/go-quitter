package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gcmurphy/getpass"
	"golang.org/x/crypto/nacl/secretbox"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const keySize = 32
const nonceSize = 24

var goquitter = "go-quitter v0.0.6-develop"
var username = os.Getenv("GNUSOCIALUSER")
var password = os.Getenv("GNUSOCIALPASS")
var gnusocialnode = os.Getenv("GNUSOCIALNODE")
var fast bool = false
var apipath string = "https://" + gnusocialnode + "/api/statuses/home_timeline.json"
var configuser = ""
var configpass = ""
var confignode = ""
var configlock = ""
var configstrings = ""
var hashbar = strings.Repeat("#", 80)
var versionbar = strings.Repeat("#", 10) + "\t" + goquitter + "\t" + strings.Repeat("#", 30)

type User struct {
	Name       string `json:"name"`
	Screenname string `json:"screen_name"`
}

//var usage string

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
type Group struct {
	Id          int64  `json:"id"`
	Url         string `json:"url"`
	Nickname    string `json:"nickname"`
	Fullname    string `json:"fullname"`
	Member      bool   `json:"member"`
	Membercount int64  `json:"member_count"`
	Description string `json:"description"`
}
type Badrequest struct {
	Error   string `json:"error"`
	Request string `json:"request"`
}

var usage = "\n\t" + goquitter + "\t" + `Copyright 2016 aerth@sdf.org

go-quitter config		Creates config file
go-quitter read			Reads 20 new posts
go-quitter read fast		Reads 20 new posts (no delay)
go-quitter home			Your home timeline.
go-quitter user username	Looks up "username" timeline
go-quitter post ____ 		Posts to your node.
go-quitter post 		Post mode.
go-quitter mentions		Mentions your @name
go-quitter search ___		Searches for ____
go-quitter search		Search mode.
go-quitter follow		Follow a user
go-quitter unfollow		Unfollow a user
go-quitter groups		List all groups on current node
go-quitter mygroups		List only groups you are member of
go-quitter join ___		Join a !group
go-quitter leave ___		Part a !group (can also use part)


Set your GNUSOCIALNODE environmental variable to change nodes.
For example: "export GNUSOCIALNODE=gs.sdf.org" in your ~/.shrc or ~/.profile
`

func bar() {
	print("\033[H\033[2J")
	fmt.Println(versionbar)
}

func main() {
	// list all commands here

	allCommands := []string{"help", "config", "read", "user", "search", "home", "follow", "unfollow", "post", "mentions", "groups", "mygroups", "join", "leave", "part", "mention", "replies", "direct", "inbox", "sent"}

	// command: go-quitter
	if len(os.Args) < 2 {
		bar()
		fmt.Println("Current list of commands:")
		fmt.Println(allCommands)
		fmt.Println(hashbar)
		os.Exit(1)
	}

	if !containsString(allCommands, os.Args[1]) {
		bar()
		fmt.Println("Current list of commands:")
		fmt.Println(allCommands)
		fmt.Println(hashbar)
		os.Exit(1)
	}

	// command: go-quitter create
	if os.Args[1] == "config" {

		if DetectConfig() == false {
			bar()
			fmt.Println("Creating config file. You will be asked for your user, node, and password.")
			fmt.Println("Your password will NOT echo.")
			createConfig()
		} else {
			bar()
			fmt.Println("Config file already exists.\nIf you want to create a new config file, move or delete the existing one.")
			fmt.Println(os.Getenv("HOME") + "/.go-quitter")
			os.Exit(1)
		}
	}

	// command: go-quitter help
	helpArg := []string{"help", "halp", "usage", "-help", "-h"}
	if containsString(helpArg, os.Args[1]) {
		bar()
		fmt.Println(usage)
		fmt.Println(hashbar)
		os.Exit(1)
	}

	// command: go-quitter version (or -v)
	versionArg := []string{"version", "-v"}
	if containsString(versionArg, os.Args[1]) {
		fmt.Println(goquitter)
		os.Exit(1)
	}
	bar()

	// command requires login credentials
	needLogin := []string{"home", "follow", "unfollow", "post", "mentions", "mygroups", "join", "leave", "mention", "replies", "direct", "inbox", "sent"}
	if containsString(needLogin, os.Args[1]) {
		if DetectConfig() == true {
			username, gnusocialnode, password, _ = ReadConfig()
			fmt.Println("Config file detected.")
		} else {
			fmt.Println("No config file detected.")
		}
		// command doesn't need login
	} else {
		if DetectConfig() == true {
			//fmt.Println("Config file detected, but this command doesn't need to login.\nWould you like to select the GNU Social node using the config?\nType YES or NO (y/n)")
			//if askForConfirmation() == true {
			// only use gnusocial node from config
			_, gnusocialnode, _, _ = ReadConfig()
			//}
		} else {
			// We are relying on environmental vars or default node.
		}

	}
	// user environmental credentials if they exist
	if os.Getenv("GNUSOCIALUSER") != "" {
		username = os.Getenv("GNUSOCIALUSER")
	}
	if os.Getenv("GNUSOCIALPASS") != "" {
		password = os.Getenv("GNUSOCIALPASS")
	}
	if os.Getenv("GNUSOCIALNODE") != "" {
		gnusocialnode = os.Getenv("GNUSOCIALNODE")
	}

	// Set speed default slow
	speed := false
	lastvar := len(os.Args)
	lastvar = (lastvar - 1)
	if os.Args[lastvar] == "fast" || os.Getenv("GNUSOCIALFAST") == "true" {
		speed = true
	}
	// command: go-quitter read
	if os.Args[1] == "read" {
		readNew(speed)
		os.Exit(0)
	}
	// command: go-quitter search _____
	if os.Args[1] == "search" {
		searchstr := ""
		if len(os.Args) > 1 {
			searchstr = strings.Join(os.Args[2:], " ")
		}
		readSearch(searchstr, speed)
		os.Exit(0)
	}

	// command: go-quitter user aerth
	if os.Args[1] == "user" && os.Args[2] != "" {
		userlookup := os.Args[2]
		readUserposts(userlookup, speed)
		os.Exit(0)
	}

	// command: go-quitter mentions
	if os.Args[1] == "mentions" || os.Args[1] == "replies" || os.Args[1] == "mention" {
		readMentions(speed)
		os.Exit(0)
	}

	// command: go-quitter follow
	if os.Args[1] == "follow" {
		followstr := ""
		if len(os.Args) == 1 {
			followstr = os.Args[2]
		} else if len(os.Args) > 1 {
			followstr = strings.Join(os.Args[2:], " ")
		}
		DoFollow(followstr)
		os.Exit(0)
	}

	// command: go-quitter unfollow
	if os.Args[1] == "unfollow" {
		followstr := ""
		if len(os.Args) == 1 {
			followstr = os.Args[2]
		} else if len(os.Args) > 1 {
			followstr = strings.Join(os.Args[2:], " ")
		}
		DoUnfollow(followstr)
		os.Exit(0)
	}
	// command: go-quitter home
	if os.Args[1] == "home" {
		readHome(speed)
		os.Exit(0)
	}

	// command: go-quitter groups
	if os.Args[1] == "groups" {
		ListAllGroups(speed)
		os.Exit(0)
	}

	// command: go-quitter mygroups
	if os.Args[1] == "mygroups" {
		ListMyGroups(speed)
		os.Exit(0)
	}
	// command: go-quitter join
	if os.Args[1] == "join" {
		content := ""
		if len(os.Args) > 1 {
			content = strings.Join(os.Args[2:], " ")
		}
		JoinGroup(content)
		os.Exit(0)
	}

	// command: go-quitter part
	if os.Args[1] == "part" {
		content := ""
		if len(os.Args) > 1 {
			content = strings.Join(os.Args[2:], " ")
		}
		PartGroup(content)
		os.Exit(0)
	}
	// command: go-quitter leave
	if os.Args[1] == "leave" {
		content := ""
		if len(os.Args) > 1 {
			content = strings.Join(os.Args[2:], " ")
		}
		PartGroup(content)
		os.Exit(0)
	}

	// go-quitter post Testing from console line using go-quitter
	// Notice how we dont need quotation marks.
	if os.Args[1] == "post" {
		content := ""
		if len(os.Args) > 1 {
			content = strings.Join(os.Args[2:], " ")
		}
		postNew(content)
		os.Exit(0)
	}

	// this happens if we invoke with somehing like "go-quitter test"
	fmt.Println(os.Args[0] + " -h")
	os.Exit(1)

}

// readNew shows 20 new messages. Defaults to a 2 second delay, but can be called with readNew(fast) for a quick dump.
func readNew(fast bool) {
	fmt.Println("node: " + gnusocialnode)
	res, err := http.Get("https://" + gnusocialnode + "/api/statuses/public_timeline.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
			time.Sleep(500 * time.Millisecond)
		}
	}

}

// readMentions shows 20 newest mentions of your username. Defaults to a 2 second delay, but can be called with readNew(fast) for a quick dump.
func readMentions(fast bool) {
	if username == "" || password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	fmt.Println("node: " + gnusocialnode)
	apipath := "https://" + gnusocialnode + "/api/statuses/mentions.json"
	req, err := http.NewRequest("GET", apipath, nil)
	req.Header.Set("User-Agent", goquitter)
	req.SetBasicAuth(username, password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var tweets []Tweet

	var apres Badrequest
	_ = json.Unmarshal(body, &apres)
	if apres.Error != "" {
		fmt.Println(apres.Error)
		os.Exit(1)
	}

	_ = json.Unmarshal(body, &tweets)

	for i := range tweets {
		if tweets[i].User.Screenname == tweets[i].User.Name {
			fmt.Printf("[@" + tweets[i].User.Screenname + "] " + tweets[i].Text + "\n\n")
		} else {
			fmt.Printf("@" + tweets[i].User.Screenname + " [" + tweets[i].User.Name + "] " + tweets[i].Text + "\n\n")
		}
		if fast != true {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// readHome shows 20 from home timeline. Defaults to a 2 second delay, but can be called with readHome(fast) for a quick dump.
func readHome(fast bool) {
	if username == "" || password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to view home timeline.")
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
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var apres Badrequest
	_ = json.Unmarshal(body, &apres)
	if apres.Error != "" {
		fmt.Println(apres.Error)
		os.Exit(1)
	}
	var tweets []Tweet
	_ = json.Unmarshal(body, &tweets)
	for i := range tweets {
		if tweets[i].User.Screenname == tweets[i].User.Name {
			fmt.Printf("[@" + tweets[i].User.Screenname + "] " + tweets[i].Text + "\n\n")
		} else {
			fmt.Printf("@" + tweets[i].User.Screenname + " [" + tweets[i].User.Name + "] " + tweets[i].Text + "\n\n")
		}
		if fast != true {
			time.Sleep(500 * time.Millisecond)
		}
	}

}

// command: go-quitter search
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
	searchq := url.Values.Encode(v)

	apipath := "https://" + gnusocialnode + "/api/search.json?" + searchq
	req, err := http.NewRequest("GET", apipath, nil)
	req.Header.Set("User-Agent", goquitter)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var apres Badrequest
	_ = json.Unmarshal(body, &apres)
	if apres.Error != "" {
		fmt.Println(apres.Error)
		os.Exit(1)
	}

	var tweets []Tweet
	_ = json.Unmarshal(body, &tweets)
	if len(tweets) == 0 {
		fmt.Println("No results for \"" + searchstr + "\"")
	}

	for i := range tweets {
		if tweets[i].User.Screenname == tweets[i].User.Name {
			fmt.Printf("[@" + tweets[i].User.Screenname + "] " + tweets[i].Text + "\n\n")
		} else {
			fmt.Printf("@" + tweets[i].User.Screenname + " [" + tweets[i].User.Name + "] " + tweets[i].Text + "\n\n")
		}
		if fast != true {
			time.Sleep(500 * time.Millisecond)
		}
	}

}

// command: go-quitter follow
func DoFollow(followstr string) {
	if username == "" || password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	if followstr == "" {
		fmt.Println("Who to follow?\nExample: someone (without the @)")
		followstr = getTypin()
	}
	if followstr == "" {
		log.Fatalln("Blank search detected. Not going furthur.")
	}
	//fmt.Println("following " + followstr + " @ " + gnusocialnode)
	v := url.Values{}

	v.Set("id", followstr)
	followstr = url.Values.Encode(v)
	b := bytes.NewBufferString(v.Encode())
	apipath := "https://" + gnusocialnode + "/api/friendships/create.json?" + followstr
	req, err := http.NewRequest("POST", apipath, b)
	req.SetBasicAuth(username, password)
	req.Header.Set("HTTP_REFERER", "https://"+gnusocialnode+"/")
	req.Header.Add("Content-Type", "[application/json; charset=utf-8")
	req.Header.Set("User-Agent", goquitter)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var apres Badrequest
	_ = json.Unmarshal(body, &apres)
	if apres.Error != "" {
		fmt.Println(apres.Error)
		os.Exit(1)
	}

	body, _ = ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	var user []User
	_ = json.Unmarshal(body, &user)

	for i := range user {
		fmt.Printf("[@" + user[i].Screenname + "]\n\n")
	}

}

func DoUnfollow(followstr string) {
	if username == "" || password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	if followstr == "" {
		fmt.Println("Who to unfollow?\nExample: someone (without the @)")
		followstr = getTypin()
	}
	if followstr == "" {
		log.Fatalln("Blank search detected. Not going furthur.")
	}
	//fmt.Println("following " + followstr + " @ " + gnusocialnode)
	v := url.Values{}

	v.Set("id", followstr)
	followstr = url.Values.Encode(v)
	b := bytes.NewBufferString(v.Encode())
	apipath := "https://" + gnusocialnode + "/api/friendships/destroy.json?" + followstr
	req, err := http.NewRequest("POST", apipath, b)
	req.SetBasicAuth(username, password)
	req.Header.Set("HTTP_REFERER", "https://"+gnusocialnode+"/")
	req.Header.Add("Content-Type", "[application/json; charset=utf-8")
	req.Header.Set("User-Agent", goquitter)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var apres Badrequest
	_ = json.Unmarshal(body, &apres)
	if apres.Error != "" {
		fmt.Println(apres.Error)
		os.Exit(1)
	}

	var user []User
	_ = json.Unmarshal(body, &user)

	for i := range user {
		fmt.Printf("[@" + user[i].Screenname + "]\n\n")
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
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var apres Badrequest
	_ = json.Unmarshal(body, &apres)
	if apres.Error != "" {
		fmt.Println(apres.Error)
		os.Exit(1)
	}

	var tweets []Tweet
	_ = json.Unmarshal(body, &tweets)
	for i := range tweets {
		if tweets[i].User.Screenname == tweets[i].User.Name {
			fmt.Printf("[@" + tweets[i].User.Screenname + "] " + tweets[i].Text + "\n\n")
		} else {
			fmt.Printf("@" + tweets[i].User.Screenname + " [" + tweets[i].User.Name + "] " + tweets[i].Text + "\n\n")
		}
		if fast != true {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func postNew(content string) {
	if username == "" || password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	// command: go-quitter post
	if content == "" {
		content = getTypin()
	}
	if content == "" {
		fmt.Println("Blank status detected. Not posting.")
		os.Exit(1)
	}
	fmt.Println("Preview:\n\n[" + username + "] " + content)
	fmt.Println("\nType YES to publish!")
	if askForConfirmation() == false {
		fmt.Println("Your status was not updated.")
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
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) == "" {
		fmt.Println("\nnode response:", resp.Status)
	}
	var apres Badrequest
	_ = json.Unmarshal(body, &apres)
	if apres.Error != "" {
		fmt.Println(apres.Error)
		os.Exit(1)
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

// Unexpected newline
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	quitResponses := []string{"q", "Q", "exit", "quit"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else if containsString(quitResponses, response) {
		return false
	} else {
		fmt.Println("\nNot valid answer, try again. [y/n] [yes/no]")
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
	fmt.Printf("\nPress ENTER when you are finished typing.\n\n")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text()
		//	fmt.Println(line)
		return line
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return ""
}

func createConfig() {
	bar()
	fmt.Printf("\nWhat username? Example: aerth")
	username = getTypin()
	if username == "" {
			fmt.Printf("\nWhat username? Example: aerth")
			username = getTypin()
	} // try 2
	if username == "" {
			fmt.Printf("\nWhat username? Example: aerth")
			username = getTypin()
	} // try 3
	if username == "" {
		// we tried.
		fmt.Println("Need real username.")
		os.Exit(1)
	}
	bar()
	fmt.Printf("\nWhich GNU Social node? Example: gnusocial.de\nPress ENTER to use gs.sdf.org")
	gnusocialnode = getTypin()
	if gnusocialnode == "" {
		gnusocialnode = "gs.sdf.org"
	}
	schema := []string{"http://", "http", "https", "://", "/"}
	if containsString(schema, gnusocialnode) {
		fmt.Printf("\nexample: gs.sdf.org")
		gnusocialnode = getTypin()
	}
	bar()
	fmt.Println("What is your GNU Social password for " + gnusocialnode + "?")
	password, _ = getpass.GetPass()
	if password == "" {
			fmt.Println("What is your GNU Social password for " + gnusocialnode + "?")
		password, _ = getpass.GetPass()
	} // try 2
	if password == "" {
			fmt.Println("What is your GNU Social password for " + gnusocialnode + "?")
		password, _ = getpass.GetPass()
	} // try 3
	if password == "" {
		// we tried.
		fmt.Println("Need real password.")
		os.Exit(1)
	}
	bar()
	fmt.Println("Enter a password to use with go-quitter.")
	fmt.Println("It will be used to encrypt your config file.")
	configlock, _ = getpass.GetPass()
	if configlock == "" {
		fmt.Println("Press ENTER again for a blank password.")
		configlock, _ = getpass.GetPass()
	} // confirm empty password
	bar()
	var userKey = configlock
	var pad = []byte("«super jumpy fox jumps all over»")
	var message = []byte(username + "::::" + gnusocialnode + "::::" + password)
	key := []byte(userKey)
	key = append(key, pad...)
	naclKey := new([keySize]byte)
	copy(naclKey[:], key[:keySize])
	nonce := new([nonceSize]byte)
	// Read bytes from random and put them in nonce until it is full.
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		fmt.Println("Could not read from random:", err)
		os.Exit(1)
	}
	out := make([]byte, nonceSize)
	copy(out, nonce[:])
	out = secretbox.Seal(out, message, nonce, naclKey)
	err = ioutil.WriteFile(os.Getenv("HOME")+"/.go-quitter", out, 0600)
	if err != nil {
		fmt.Println("Error while writing config file: ", err)
		os.Exit(1)
	}
	fmt.Printf("Config file saved at "+os.Getenv("HOME")+"/.go-quitter \nTotal size is %d bytes.\n",
		len(out))
	os.Exit(0)
}

func DetectConfig() bool {
	_, err := ioutil.ReadFile(os.Getenv("HOME") + "/.go-quitter")
	if err != nil {
		return false
	}
	return true
}

func ReadConfig() (configuser string, confignode string, configpass string, err error) {
	print("\033[H\033[2J")
	fmt.Println(versionbar)
	fmt.Println("Unlocking config file")
	configlock, err = getpass.GetPass()
	print("\033[H\033[2J")
	fmt.Println(versionbar)
	var userKey = configlock
	var pad = []byte("«super jumpy fox jumps all over»")
	key := []byte(userKey)
	key = append(key, pad...)
	naclKey := new([keySize]byte)
	copy(naclKey[:], key[:keySize])
	nonce := new([nonceSize]byte)
	in, err := ioutil.ReadFile(os.Getenv("HOME") + "/.go-quitter")
	if err != nil {
		log.Fatalln(err)
	}
	copy(nonce[:], in[:nonceSize])
	configbytes, ok := secretbox.Open(nil, in[nonceSize:], nonce, naclKey)
	if ok {
		fmt.Println("Logged in. Welcome back to go-quitter!")
	} else {
		log.Fatalln("Could not decrypt the config file. Wrong password?")
	}
	configstrings := strings.Split(string(configbytes), "::::")

	username = configstrings[0]
	gnusocialnode = configstrings[1]
	password = configstrings[2]

	return username, gnusocialnode, password, nil

}

// command: go-quitter groups
func ListAllGroups(speed bool) {
	if username == "" || password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	initwin()
	apipath := "https://" + gnusocialnode + "/api/statusnet/groups/list_all.json"

	req, err := http.NewRequest("GET", apipath, nil)
	req.SetBasicAuth(username, password)
	req.Header.Set("HTTP_REFERER", "https://"+gnusocialnode+"/")
	req.Header.Set("User-Agent", goquitter)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if string(body) == "" {
		fmt.Println("\nnode response:", resp.Status)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var apres Badrequest
	_ = json.Unmarshal(body, &apres)
	if apres.Error != "" {
		fmt.Println(apres.Error)
		os.Exit(1)
	}

	var groups []Group
	_ = json.Unmarshal(body, &groups)
	var member string
	var members string
	var id string
	for i := range groups {
		if groups[i].Member == true {
			member = `*member*`
		} else {
			member = ""
		}

		id = strconv.FormatInt(groups[i].Id, 10)
		members = strconv.FormatInt(groups[i].Membercount, 10)
		fmt.Printf("!" + groups[i].Nickname + " (#" + id + ") [" + groups[i].Fullname + "] " + member + "\n" + groups[i].Description + "\n" + groups[i].Nickname + " has " + members + " members total.\n\n")
		if speed != true {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// command: go-quitter mygroups

func ListMyGroups(speed bool) {
	if username == "" || password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	initwin()
	apipath := "https://" + gnusocialnode + "/api/statusnet/groups/list.json"
	req, err := http.NewRequest("GET", apipath, nil)
	req.SetBasicAuth(username, password)
	req.Header.Set("HTTP_REFERER", "https://"+gnusocialnode+"/")
	req.Header.Set("User-Agent", goquitter)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) == "" {
		fmt.Println("\nnode response:", resp.Status)
	}
	var apres Badrequest
	_ = json.Unmarshal(body, &apres)
	if apres.Error != "" {
		fmt.Println(apres.Error)
		os.Exit(1)
	}
	//	fmt.Println(string(body))
	var groups []Group
	_ = json.Unmarshal(body, &groups)

	for i := range groups {

		fmt.Printf("!" + groups[i].Nickname + " [" + groups[i].Fullname + "] \n" + groups[i].Description + "\n\n")
		if speed != true {
			time.Sleep(500 * time.Millisecond)
		}
	}

}

// command: go-quitter join ____
func JoinGroup(groupstr string) {
	if username == "" || password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	if groupstr == "" {
		fmt.Println("Which group to join?\nExample: groupname (without the !)")
		groupstr = getTypin()
	}
	if groupstr == "" {
		log.Fatalln("Blank group detected. Not going furthur.")
	}
	v := url.Values{}

	v.Set("group_name", groupstr)
	v.Set("group_id", groupstr)
	v.Set("id", groupstr)
	v.Set("nickname", groupstr)
	groupstr = url.Values.Encode(v)
	b := bytes.NewBufferString(v.Encode())
	apipath := "https://" + gnusocialnode + "/api/statusnet/groups/join.json?" + groupstr
	req, err := http.NewRequest("POST", apipath, b)
	req.SetBasicAuth(username, password)
	req.Header.Set("HTTP_REFERER", "https://"+gnusocialnode+"/")
	req.Header.Add("Content-Type", "[application/json; charset=utf-8")
	req.Header.Set("User-Agent", goquitter)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if string(body) == "" {
		fmt.Println("\nnode response:", resp.Status)
	}
	var apres Badrequest
	_ = json.Unmarshal(body, &apres)
	if apres.Error != "" {
		fmt.Println(apres.Error)
		os.Exit(1)
	}

	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	var user []User
	_ = json.Unmarshal(body, &user)

	for i := range user {
		fmt.Printf("[@" + user[i].Screenname + "]\n\n")
	}
}

// command: go-quitter part ____
func PartGroup(groupstr string) {
	if username == "" || password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	if groupstr == "" {
		fmt.Println("Which group to leave?\nExample: groupname (without the !)")
		groupstr = getTypin()
	}
	if groupstr == "" {
		log.Fatalln("Blank group detected. Not going furthur.")
	}
	fmt.Println("Are you sure you want to leave from group !" + groupstr + "\n Type yes or no [y/n]\n")
	if askForConfirmation() == false {

		fmt.Println("Not leaving group " + groupstr)
		os.Exit(0)

	}

	v := url.Values{}

	v.Set("group_name", groupstr)
	v.Set("group_id", groupstr)
	v.Set("id", groupstr)
	v.Set("nickname", groupstr)
	groupstr = url.Values.Encode(v)
	b := bytes.NewBufferString(v.Encode())
	apipath := "https://" + gnusocialnode + "/api/statusnet/groups/leave.json?" + groupstr
	req, err := http.NewRequest("POST", apipath, b)
	req.SetBasicAuth(username, password)
	req.Header.Set("HTTP_REFERER", "https://"+gnusocialnode+"/")
	req.Header.Add("Content-Type", "[application/json; charset=utf-8")
	req.Header.Set("User-Agent", goquitter)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var apres Badrequest
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) == "" {
		fmt.Println("\nnode response:", resp.Status)
	}
	_ = json.Unmarshal(body, &apres)

	fmt.Println(apres.Error)

	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	var user []User
	_ = json.Unmarshal(body, &user)

	for i := range user {
		fmt.Printf("[@" + user[i].Screenname + "]\n\n")
	}
}

// This will change with the real UI
func initwin() {
	print("\033[H\033[2J")
	fmt.Println(versionbar)
}
