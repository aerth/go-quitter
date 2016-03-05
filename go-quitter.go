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
//	"strconv"
	"strings"

)

const keySize = 32
const nonceSize = 24
var goquitter = "go-quitter v0.0.7"
var fast bool = false
var hashbar = strings.Repeat("#", 80)
var versionbar = strings.Repeat("#", 10) + "\t" + goquitter + "\t" + strings.Repeat("#", 30)

// Basic https authentication struct. Set it with NewAuth()
type Auth struct {
	Username string
	Password string
	Node     string
}

// Sets the Authentication method and choose node.
//Use like this:
/*

 q := qw.NewAuth()
 q.Username = "john"
 q.Password = "pass123"
 q.Node = "gnusocial.de"
 q.GetHome(false)

*/
func NewAuth() *Auth {
	return &Auth{
		Username: "gopher",
		Password: "password",
		Node:     "gs.sdf.org",
	}
}

var apipath string = "https://null/api/statuses/home_timeline.json"

// GNU Social User, gets returned by GS API
type User struct {
	Name       string `json:"name"`
	Screenname string `json:"screen_name"`
}

// GNU Social Quip, gets returned by GS API.
type Quip struct {
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
	RetweetedStatus      *Quip   `json:"retweeted_status"`
	Source               string   `json:"source"`
	Text                 string   `json:"text"`
	Truncated            bool     `json:"truncated"`
	User                 User     `json:"user"`
	WithheldCopyright    bool     `json:"withheld_copyright"`
	WithheldInCountries  []string `json:"withheld_in_countries"`
	WithheldScope        string   `json:"withheld_scope"`
}

// GNU Social Group, gets returned by GS API.
type Group struct {
	Id          int64  `json:"id"`
	Url         string `json:"url"`
	Nickname    string `json:"nickname"`
	Fullname    string `json:"fullname"`
	Member      bool   `json:"member"`
	Membercount int64  `json:"member_count"`
	Description string `json:"description"`
}

// If the API doesn't respond how we like, it replies using this struct (in json)
type Badrequest struct {
	Error   string `json:"error"`
	Request string `json:"request"`
}



// GetPublic shows 20 new messages. Defaults to a 2 second delay, but can be called with GetPublic(fast) for a quick dump. This and DoSearch() and GetUserTimeline() are some of the only functions that don't require auth.Username + auth.Password
func (a Auth) GetPublic(fast bool) ([]Quip, error) {
	resp, err := http.Get("https://" + a.Node + "/api/statuses/public_timeline.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var quips []Quip

	var apres Badrequest
	_ = json.Unmarshal(body, &apres)
	if apres.Error != "" {
		return nil, NewError(apres.Error)
		os.Exit(1)
	}

	_ = json.Unmarshal(body, &quips)

	return quips, err
}

// GetMentions shows 20 newest mentions of your username. Defaults to a 2 second delay, but can be called with GetPublic(fast) for a quick dump.
func (a Auth) GetMentions(fast bool) ([]Quip, error) {
	if a.Username == "" || a.Password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	//fmt.Println("node: " + a.Node)
	apipath := "https://" + a.Node + "/api/statuses/mentions.json"
	req, err := http.NewRequest("GET", apipath, nil)
	if req.Header == nil {
		req.Header = http.Header{}
	}
	if req.Header.Get("User-Agent") == "" {
		if req.Header == nil {
  req.Header = http.Header{}
}
if req.Header.Get("User-Agent") == "" {
  req.Header.Set("User-Agent", goquitter)
}
	}
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

	var quips []Quip

	var apres Badrequest
	_ = json.Unmarshal(body, &apres)
	if apres.Error != "" {
		fmt.Println(apres.Error)
		os.Exit(1)
	}

	_ = json.Unmarshal(body, &quips)

	return quips, err
}

// GetHome shows 20 from home timeline. Defaults to a 2 second delay, but can be called with GetHome(fast) for a quick dump.
func (a Auth) GetHome(fast bool) ([]Quip, error) {
	if a.Username == "" || a.Password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to view home timeline.")
	}
	//fmt.Println("node: " + a.Node)
	apipath := "https://" + a.Node + "/api/statuses/home_timeline.json"
	req, err := http.NewRequest("GET", apipath, nil)
	if req.Header == nil {
  req.Header = http.Header{}
}
if req.Header.Get("User-Agent") == "" {
  req.Header.Set("User-Agent", goquitter)
}
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
	var quips []Quip
	_ = json.Unmarshal(body, &quips)


	return quips, err

}

// command: go-quitter search
func (a Auth) DoSearch(searchstr string, fast bool) ([]Quip, error) {
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
	if req.Header == nil {
  req.Header = http.Header{}
}
if req.Header.Get("User-Agent") == "" {
  req.Header.Set("User-Agent", goquitter)
}
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

	var quips []Quip
	_ = json.Unmarshal(body, &quips)
	if len(quips) == 0 {
		fmt.Println("No results for \"" + searchstr + "\"")
	}

 return quips, err

}

// command: go-quitter follow
func (a Auth) DoFollow(followstr string) (User, error) {
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
	if req.Header == nil {
  req.Header = http.Header{}
}
if req.Header.Get("User-Agent") == "" {
  req.Header.Set("User-Agent", goquitter)
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

	body, _ = ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	var user User
	_ = json.Unmarshal(body, &user)

	return user, err

}

// go-quitter command: go-quitter unfollow
func (a Auth) DoUnfollow(followstr string) (User, error) {
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
	if req.Header == nil {
  req.Header = http.Header{}
}
if req.Header.Get("User-Agent") == "" {
  req.Header.Set("User-Agent", goquitter)
}
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

	var user User
	_ = json.Unmarshal(body, &user)

	return user, err

}

// go-quitter command: go-quitter user
func (a Auth) GetUserTimeline(userlookup string, fast bool) ([]Quip, error) {
	fmt.Println("user " + userlookup + " @ " + a.Node)
	apipath := "https://" + a.Node + "/api/statuses/user_timeline.json?screen_name=" + userlookup
	req, err := http.NewRequest("GET", apipath, nil)
	if req.Header == nil {
  req.Header = http.Header{}
}
if req.Header.Get("User-Agent") == "" {
  req.Header.Set("User-Agent", goquitter)
}
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

	var quips []Quip
	_ = json.Unmarshal(body, &quips)
	return quips, err
}


// go-quitter command: go-quitter post
func (a Auth) PostNew(content string) (Quip, error) {
	if a.Username == "" || a.Password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	if content == "" {
		content = getTypin()
	}
	if content == "" {
		fmt.Println("Blank status detected. Not posting.")
		os.Exit(1)
	}
	fmt.Println("Preview:\n\n[" + a.Username + "] " + content)
	fmt.Println("\nType YES to publish!")
	if askForConfirmation() == false {
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
	if req.Header == nil {
  req.Header = http.Header{}
}
if req.Header.Get("User-Agent") == "" {
  req.Header.Set("User-Agent", goquitter)
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

	var quip Quip
	_ = json.Unmarshal(body, &quip)
	return quip, err


}

// Does x contain y?
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

// Ask user to confirm the action.
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

// For use only in containsString()
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// Receive non-hidden input from user.
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
func (a Auth) ListAllGroups(speed bool) ([]Group, error) {
	if a.Username == "" || a.Password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	initwin()
	apipath := "https://" + a.Node + "/api/statusnet/groups/list_all.json"

	req, err := http.NewRequest("GET", apipath, nil)
	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("HTTP_REFERER", "https://"+a.Node+"/")
	if req.Header == nil {
  req.Header = http.Header{}
}
if req.Header.Get("User-Agent") == "" {
  req.Header.Set("User-Agent", goquitter)
}
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

	return groups, err
}

// command: go-quitter mygroups
func (a Auth) ListMyGroups(speed bool) ([]Group, error) {
	if a.Username == "" || a.Password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	initwin()
	apipath := "https://" + a.Node + "/api/statusnet/groups/list.json"
	req, err := http.NewRequest("GET", apipath, nil)
	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("HTTP_REFERER", "https://"+a.Node+"/")
	if req.Header == nil {
  req.Header = http.Header{}
}
if req.Header.Get("User-Agent") == "" {
  req.Header.Set("User-Agent", goquitter)
}
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

	return groups, err

}

// command: go-quitter join ____
func (a Auth) JoinGroup(groupstr string) (Group, error) {
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
	if strings.HasPrefix(groupstr, "!"){
		groupstr = strings.StripPrefix(groupstr, "!")
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
	if req.Header == nil {
  req.Header = http.Header{}
}
if req.Header.Get("User-Agent") == "" {
  req.Header.Set("User-Agent", goquitter)
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

	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	var group Group
	_ = json.Unmarshal(body, &group)

	return group, err
}

// command: go-quitter part ____
func (a Auth) PartGroup(groupstr string) (Group, error) {
	if a.Username == "" || a.Password == "" {
		log.Fatalln("Please run \"go-quitter config\" or set the GNUSOCIALUSER and GNUSOCIALPASS environmental variables to post.")
	}
	if groupstr == "" {
		fmt.Println("Which group to leave?\nExample: groupname or !group")
		groupstr = getTypin()
	}
	if groupstr == "" {
		log.Fatalln("Blank group detected. Not going furthur.")
	}
	if strings.HasPrefix(groupstr, "!"){
		groupstr = strings.StripPrefix(groupstr, "!")
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
	apipath := "https://" + a.Node + "/api/statusnet/groups/leave.json?" + groupstr
	req, err := http.NewRequest("POST", apipath, b)
	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("HTTP_REFERER", "https://"+a.Node+"/")
	req.Header.Add("Content-Type", "[application/json; charset=utf-8")
	if req.Header == nil {
  req.Header = http.Header{}
}
if req.Header.Get("User-Agent") == "" {
  req.Header.Set("User-Agent", goquitter)
}
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
	var group Group
	_ = json.Unmarshal(body, &group)

	return group, err
}

// This will change with the real UI. Ugly on windows.
func initwin() {
	print("\033[H\033[2J")
	fmt.Println(versionbar)
}

// ReturnHome gives us the true home directory for letting the user know where the config file is. Windows, Unix, OS X
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
func NewError(text string) error {
    return &errorString{text}
}
// errorString is a trivial implementation of error.
type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return e.s
}
