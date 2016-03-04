// The MIT License (MIT)
//
// Copyright (c) 2016 aerth
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

package quitter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
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

var goquitter = "go-quitter v0.0.7"

var fast bool = false

var hashbar = strings.Repeat("#", 80)
var versionbar = strings.Repeat("#", 10) + "\t" + goquitter + "\t" + strings.Repeat("#", 30)

type Auth struct {
	Username string
	Password string
	Node string
}

func NewAuth() *Auth {
	return &Auth{
		Username:         "gopher",
		Password:         "password",
		Node:         "gs.sdf.org",
	}
}

var apipath string = "https://null/api/statuses/home_timeline.json"


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

func bar() {
	print("\033[H\033[2J")
	fmt.Println(versionbar)
}

// ReadPublic shows 20 new messages. Defaults to a 2 second delay, but can be called with ReadPublic(fast) for a quick dump.
func (a Auth) ReadPublic(fast bool) {
	fmt.Println("node: " + a.Node)
	res, err := http.Get("https://" + a.Node + "/api/statuses/public_timeline.json")
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

// ReadMentions shows 20 newest mentions of your username. Defaults to a 2 second delay, but can be called with ReadPublic(fast) for a quick dump.
func (a *Auth) ReadMentions(fast bool) {
	if a.Username == "" || a.Password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	fmt.Println("node: " + a.Node)
	apipath := "https://" + a.Node + "/api/statuses/mentions.json"
	req, err := http.NewRequest("GET", apipath, nil)
	req.Header.Set("User-Agent", goquitter)
	req.SetBasicAuth(a.Username, a.Password)
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

// ReadHome shows 20 from home timeline. Defaults to a 2 second delay, but can be called with ReadHome(fast) for a quick dump.
func (a Auth) ReadHome(fast bool) {
	if a.Username == "" || a.Password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to view home timeline.")
	}
	fmt.Println("node: " + a.Node)
	apipath := "https://" + a.Node + "/api/statuses/home_timeline.json"
	req, err := http.NewRequest("GET", apipath, nil)
	req.Header.Set("User-Agent", goquitter)
	req.SetBasicAuth(a.Username, a.Password)
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
func (a Auth) DoSearch(searchstr string, fast bool) {
	if searchstr == "" {
		searchstr = getTypin()
	}
	if searchstr == "" {
		log.Fatalln("Blank search detected. Not searching.")
	}
	fmt.Println("searching " + searchstr + " @ " + a.Node)
	v := url.Values{}
	v.Set("q", searchstr)
	searchq := url.Values.Encode(v)

	apipath := "https://" + a.Node + "/api/search.json?" + searchq
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
func (a Auth) DoFollow(followstr string) {
	if a.Username == "" || a.Password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	if followstr == "" {
		fmt.Println("Who to follow?\nExample: someone (without the @)")
		followstr = getTypin()
	}
	if followstr == "" {
		log.Fatalln("Blank search detected. Not going furthur.")
	}
	//fmt.Println("following " + followstr + " @ " + a.Node)
	v := url.Values{}

	v.Set("id", followstr)
	followstr = url.Values.Encode(v)
	b := bytes.NewBufferString(v.Encode())
	apipath := "https://" + a.Node + "/api/friendships/create.json?" + followstr
	req, err := http.NewRequest("POST", apipath, b)
	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("HTTP_REFERER", "https://"+a.Node+"/")
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

func (a Auth) DoUnfollow(followstr string) {
	if a.Username == "" || a.Password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	if followstr == "" {
		fmt.Println("Who to unfollow?\nExample: someone (without the @)")
		followstr = getTypin()
	}
	if followstr == "" {
		log.Fatalln("Blank search detected. Not going furthur.")
	}
	//fmt.Println("following " + followstr + " @ " + a.Node)
	v := url.Values{}

	v.Set("id", followstr)
	followstr = url.Values.Encode(v)
	b := bytes.NewBufferString(v.Encode())
	apipath := "https://" + a.Node + "/api/friendships/destroy.json?" + followstr
	req, err := http.NewRequest("POST", apipath, b)
	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("HTTP_REFERER", "https://"+a.Node+"/")
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

func (a Auth) GetUserTimeline(userlookup string, fast bool) {
	fmt.Println("user " + userlookup + " @ " + a.Node)
	apipath := "https://" + a.Node + "/api/statuses/user_timeline.json?screen_name=" + userlookup
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

func (a Auth) PostNew(content string) {
	if a.Username == "" || a.Password == "" {
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
	fmt.Println("Preview:\n\n[" + a.Username + "] " + content)
	fmt.Println("\nType YES to publish!")
	if AskForConfirmation() == false {
		fmt.Println("Your status was not updated.")
		os.Exit(0)
	}
	fmt.Println("posting on node: " + a.Node)
	v := url.Values{}
	v.Set("status", content)
	content = url.Values.Encode(v)
	apipath := "https://" + a.Node + "/api/statuses/update.json?" + content
	req, err := http.NewRequest("POST", apipath, nil)
	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("HTTP_REFERER", "https://"+a.Node+"/")
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

func ContainsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

// Unexpected newline
func AskForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	quitResponses := []string{"q", "Q", "exit", "quit"}
	if ContainsString(okayResponses, response) {
		return true
	} else if ContainsString(nokayResponses, response) {
		return false
	} else if ContainsString(quitResponses, response) {
		return false
	} else {
		fmt.Println("\nNot valid answer, try again. [y/n] [yes/no]")
		return AskForConfirmation()
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

// command: go-quitter groups
func (a Auth) ListAllGroups(speed bool) {
	if a.Username == "" || a.Password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	initwin()
	apipath := "https://" + a.Node + "/api/statusnet/groups/list_all.json"

	req, err := http.NewRequest("GET", apipath, nil)
	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("HTTP_REFERER", "https://"+a.Node+"/")
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

func (a Auth) ListMyGroups(speed bool) {
	if a.Username == "" || a.Password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	initwin()
	apipath := "https://" + a.Node + "/api/statusnet/groups/list.json"
	req, err := http.NewRequest("GET", apipath, nil)
	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("HTTP_REFERER", "https://"+a.Node+"/")
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
func (a Auth) JoinGroup(groupstr string) {
	if a.Username == "" || a.Password == "" {
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
	apipath := "https://" + a.Node + "/api/statusnet/groups/join.json?" + groupstr
	req, err := http.NewRequest("POST", apipath, b)
	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("HTTP_REFERER", "https://"+a.Node+"/")
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
func (a Auth) PartGroup(groupstr string) {
	if a.Username == "" || a.Password == "" {
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
	if AskForConfirmation() == false {
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
	apipath := "https://" + a.Node + "/api/statusnet/groups/leave.json?" + groupstr
	req, err := http.NewRequest("POST", apipath, b)
	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("HTTP_REFERER", "https://"+a.Node+"/")
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

// This will change with the real UI. Ugly on windows.
func initwin() {
	print("\033[H\033[2J")
	fmt.Println(versionbar)
}

// ReturnHome gives us the true home directory for letting the user know where the config file is.
func ReturnHome() (homedir string) {
	homedir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if homedir == "" {
		homedir = os.Getenv("USERPROFILE")
	}
	if homedir == "" {
		homedir = os.Getenv("HOME")
	}
	return
}

func init() {
//	if a.Node == "" {
//		a.Node = "gs.sdf.org"
//	}

}
