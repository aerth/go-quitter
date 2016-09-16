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

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	qw "github.com/aerth/go-quitter"
	"github.com/aerth/seconf"
)

var goquitter = "go-quitter v0.0.8"
var username = os.Getenv("GNUSOCIALUSER")
var password = os.Getenv("GNUSOCIALPASS")
var gnusocialnode = os.Getenv("GNUSOCIALNODE")
var fast bool = false
var apipath string = "https://" + gnusocialnode + "/api/statuses/home_timeline.json"
var gnusocialpath = "go-quitter"
var configuser = ""
var configpass = ""
var confignode = ""
var configlock = ""
var configstrings = ""
var hashbar = ""

//var hashbar = strings.Repeat("#", 80)
//var versionbar = strings.Repeat("#", 10) + "\t" + goquitter + "\t" + strings.Repeat("#", 30)

var usage = "\n" + "\t" + `  Copyright 2016 aerth@sdf.org

go-quitter config		Creates config file	*do this first*
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

Using environmental variables will override the config:

GNUSOCIALNODE
GNUSOCIALPASS
GNUSOCIALUSER
GNUSOCIALPATH

Set your environmental variable to change nodes, use a different config,
	or change user or password for a one-time session.

For example: "export GNUSOCIALNODE=gs.sdf.org" in your ~/.shrc or ~/.profile
`

func init() {
	if gnusocialnode == "" {
		gnusocialnode = "gnusocial.de"
	}
}
func bar() {
	//	print("\033[H\033[2J")
	//	fmt.Println(versionbar)
}

func main() {
	q := qw.NewAuth()
	if containsString(os.Args, "-debug") {
		q.Scheme = "http://"
	}

	//	os.Exit(1)

	// list all commands here
	if os.Getenv("GNUSOCIALPATH") != "" {
		gnusocialpath = os.Getenv("GNUSOCIALPATH")
	}

	allCommands := []string{"help", "config", "read", "user", "search", "home", "follow", "unfollow", "post", "mentions", "groups", "mygroups", "join", "leave", "part", "mention", "replies"}

	// command: go-quitter
	if len(os.Args) < 2 {
		fmt.Println(usage)
		fmt.Println("Current list of commands:")
		fmt.Println(allCommands)
		fmt.Println(hashbar)
		os.Exit(1)
	}

	if !containsString(allCommands, os.Args[1]) {
		fmt.Println(usage)
		fmt.Println("Current list of commands:")
		fmt.Println(allCommands)
		fmt.Println(hashbar)
		os.Exit(1)
	}

	// command: go-quitter create
	if os.Args[1] == "config" {
		if seconf.Detect(gnusocialpath) == false {
			bar()
			fmt.Println("Creating config file. You will be asked for your user, node, and password.")
			fmt.Println("Your password will not echo.")
			seconf.Create(gnusocialpath, "GNU Social", "GNU Social username", "Which GNU Social node? Example: gnusocial.de", "password: will not echo")
		} else {
			bar()
			fmt.Println("Config file already exists.\nIf you want to create a new config file, move or delete the existing one.\nYou can also set the GNUSOCIALPATH env to use multiple config files. \nExample: GNUSOCIALPATH=gnusocial.de go-quitter config")
			fmt.Println("\nConfig exists:", returnHomeDir()+"/."+gnusocialpath)
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
		if seconf.Detect(gnusocialpath) == true {
			configdecoded, err := seconf.Read(gnusocialpath)
			if err != nil {
				fmt.Println("error:")
				fmt.Println(err)
				os.Exit(1)
			}
			//configstrings := string(configdecoded)
			//		fmt.Println("config strings:")
			//		fmt.Println(configdecoded)
			configarray := strings.Split(configdecoded, "::::")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if len(configarray) != 3 {
				fmt.Println("Broken config file. Create a new one. :(")
				os.Exit(1)
			}
			username = string(configarray[0])
			gnusocialnode = string(configarray[1])
			gnusocialnode = strings.Replace(gnusocialnode, "http://", "", -1)
			gnusocialnode = strings.Replace(gnusocialnode, "https://", "", -1)
			password = string(configarray[2])

			q.Username = username
			q.Password = password
			q.Node = gnusocialnode
			if os.Getenv("GNUSOCIALUSER") != "" {
				q.Username = os.Getenv("GNUSOCIALUSER")
			}
			if os.Getenv("GNUSOCIALPASS") != "" {
				q.Password = os.Getenv("GNUSOCIALPASS")
			}
			if os.Getenv("GNUSOCIALNODE") != "" {
				q.Node = os.Getenv("GNUSOCIALNODE")
			}
			if q.Username == username {
				fmt.Println("Welcome back, " + q.Username + "@" + q.Node)
			} else {
				fmt.Println("Welcome, " + q.Username + "@" + q.Node)
			}
		} else {
			fmt.Println("No config file detected.")
		}
		// command doesn't need login
	} else {
		if seconf.Detect(gnusocialpath) == true {
			//fmt.Println("Config file detected, but this command doesn't need to login.\nWould you like to select the GNU Social node using the config?\nType YES or NO (y/n)")
			//if AskForConfirmation() == true {
			// only use gnusocial node from config
			configdecoded, err := seconf.Read(gnusocialpath)
			if err != nil {
				fmt.Println("error:")
				fmt.Println(err)
				os.Exit(1)
			}
			configarray := strings.Split(configdecoded, "::::")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if len(configarray) != 3 {
				fmt.Println("Broken config file. Create a new one.")
				os.Exit(1)
			}
			gnusocialnode = string(configarray[1])
			//fmt.Println(gnusocialnode)
			gnusocialnode = strings.Replace(gnusocialnode, "http://", "", -1)
			gnusocialnode = strings.Replace(gnusocialnode, "https://", "", -1)
			//			gnusocialnode = strings.TrimLeft(gnusocialnode, "https://")
			//			gnusocialnode = strings.TrimLeft(gnusocialnode, "http://")
			q.Node = gnusocialnode

			//}
		} else {
			// We are relying on environmental vars or default node.
		}

	}
	// user environmental credentials if they exist
	if os.Getenv("GNUSOCIALUSER") != "" {
		q.Username = os.Getenv("GNUSOCIALUSER")
	}
	if os.Getenv("GNUSOCIALPASS") != "" {
		q.Password = os.Getenv("GNUSOCIALPASS")
	}
	if os.Getenv("GNUSOCIALNODE") != "" {
		q.Node = os.Getenv("GNUSOCIALNODE")
	}

	// Set speed default slow
	speed := false
	lastvar := len(os.Args) - 1

	if os.Args[lastvar] == "fast" || os.Getenv("GNUSOCIALFAST") == "true" {
		speed = true
	}

	switch os.Args[1] {
	// command: go-quitter read
	case "read":
		PrintQuips(q.GetPublic(speed))
		os.Exit(0)

		// command: go-quitter search _____
	case "search":
		searchstr := ""
		if len(os.Args) > 1 {
			searchstr = strings.Join(os.Args[2:], " ")
		}
		if searchstr == "" {
			searchstr = getTypin()
		}
		PrintQuips(q.DoSearch(searchstr, speed))
		os.Exit(0)

		// command: go-quitter user aerth
	case "user":
		if len(os.Args) > 2 && os.Args[2] != "" {
			userlookup := os.Args[2]
			PrintQuips(q.GetUserTimeline(userlookup, speed))

			os.Exit(0)
		}
		fmt.Println("Need user to search for")
		os.Exit(1)

		// command: go-quitter mentions
	case "mentions", "replies", "mention":
		PrintQuips(q.GetMentions(speed))
		os.Exit(0)

		// command: go-quitter follow
	case "follow":
		followstr := ""
		if len(os.Args) == 1 {
			followstr = os.Args[2]
		} else if len(os.Args) > 1 {
			followstr = strings.Join(os.Args[2:], " ")
		}
		if followstr == "" {
			fmt.Println("Who to follow?\nExample: someone (without the @)")
			followstr = getTypin()
		}
		PrintUser(q.DoFollow(followstr))
		os.Exit(0)

	// command: go-quitter unfollow
	case "unfollow":
		followstr := ""
		if len(os.Args) == 1 {
			followstr = os.Args[2]
		} else if len(os.Args) > 1 {
			followstr = strings.Join(os.Args[2:], " ")
		}
		if followstr == "" {
			fmt.Println("Who to unfollow?\nExample: someone (without the @)")
			followstr = getTypin()
		}
		PrintUser(q.DoUnfollow(followstr))
		os.Exit(0)

	// command: go-quitter home
	case "home":
		PrintQuips(q.GetHome(speed))
		os.Exit(0)

	// command: go-quitter groups
	case "groups":
		PrintGroups(q.ListAllGroups(speed))
		os.Exit(0)

		// command: go-quitter mygroups
	case "mygroups":
		PrintGroups(q.ListMyGroups(speed))
		os.Exit(0)

		// command: go-quitter join
	case "join":
		groupstr := ""
		if len(os.Args) > 1 {
			groupstr = strings.Join(os.Args[2:], " ")
		}
		if groupstr == "" {
			fmt.Println("Which group to join?\nExample: groupname (without the !)")
			groupstr = getTypin()
		}
		PrintGroup(q.JoinGroup(groupstr))
		os.Exit(0)

		// command: go-quitter part
	case "part":
		groupstr := ""
		if len(os.Args) > 1 {
			groupstr = strings.Join(os.Args[2:], " ")
		}
		if groupstr == "" {
			fmt.Println("Which group to leave?\nExample: groupname (without the !)")
			groupstr = getTypin()
		}

		fmt.Println("Are you sure you want to leave from group !" + groupstr + "\n Type yes or no [y/n]\n")
		if askForConfirmation() == false {
			fmt.Println("Not leaving group " + groupstr)
			os.Exit(0)
		}

		PrintGroup(q.PartGroup(groupstr))
		os.Exit(0)

		// command: go-quitter leave
	case "leave":
		content := ""
		if len(os.Args) > 1 {
			content = strings.Join(os.Args[2:], " ")
		}
		PrintGroup(q.PartGroup(content))
		os.Exit(0)

		// go-quitter
	case "post":
		content := ""
		if len(os.Args) > 1 {
			content = strings.Join(os.Args[2:], " ") // go-quitter post wow this is a post\!
		}
		if content == "" {
			content = getTypin()
		}

		fmt.Println("Preview:\n\n[" + q.Username + "] " + content)
		fmt.Println("\nType YES to publish!")
		if askForConfirmation() == false {
			fmt.Println("Your status was not updated.")
			os.Exit(0)
		}

		PrintQuip(q.PostNew(content))
		os.Exit(0)

	default:
		// this happens if we invoke with somehing like "go-quitter test"
		fmt.Println("Command not found, try ", os.Args[0]+" help")
		os.Exit(1)
	}
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

// Does []string contain element?
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

// Find the index of a string in a []string
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

func PrintQuips(quips []qw.Quip, err error) {
	if err != nil {
		fmt.Println(err)
		return

	}
	if len(quips) == 0 && err == nil {
		fmt.Println("No results.")
		return
	}
	for i := range quips {
		if quips[i].User.Screenname == quips[i].User.Name {
			fmt.Printf("[@" + quips[i].User.Screenname + "] " + quips[i].Text + "\n\n")
		} else {
			fmt.Printf("@" + quips[i].User.Screenname + " [" + quips[i].User.Name + "] " + quips[i].Text + "\n\n")
		}
	}
}
func PrintQuip(quip qw.Quip, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	if quip.Text == "" && err == nil {
		fmt.Println("No quip.")
		return
	}
	if quip.User.Screenname == quip.User.Name {
		fmt.Printf("[@" + quip.User.Screenname + "] " + quip.Text + "\n\n")
	} else {
		fmt.Printf("@" + quip.User.Screenname + " [" + quip.User.Name + "] " + quip.Text + "\n\n")
	}

}

func PrintUsers(users []qw.User, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(users) == 0 && err == nil {
		fmt.Println("No users.")
		return
	}
	for i := range users {
		if users[i].Screenname == users[i].Name {
			fmt.Printf("[@" + users[i].Screenname + "]\n\n")
		} else {
			fmt.Printf("@" + users[i].Screenname + " [" + users[i].Name + "]\n\n")
		}
	}
}
func PrintUser(user qw.User, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	if user.Screenname == "" && err == nil {
		fmt.Println("No user.")
		return
	}
	fmt.Printf("[@" + user.Screenname + "]\n\n")

}
func PrintGroup(group qw.Group, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	if group.Nickname == "" && err == nil {
		fmt.Println("No group.")
		return
	}
	fmt.Printf("!" + group.Nickname + " [" + group.Fullname + "] \n" + group.Description + "\n\n")

}

func PrintGroups(groups []qw.Group, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(groups) == 0 && err == nil {
		fmt.Println("No groups.")
		return
	}
	for i := range groups {
		fmt.Printf("!" + groups[i].Nickname + " [" + groups[i].Fullname + "] \n" + groups[i].Description + "\n\n")
	}
}

// returnHomeDir gives us the true home directory for letting the user know where the config file is. Windows, Unix, OS X
func returnHomeDir() (homedir string) {
	homedir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if homedir == "" {
		homedir = os.Getenv("USERPROFILE")
	}
	if homedir == "" {
		homedir = os.Getenv("HOME")
	}
	return homedir
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
