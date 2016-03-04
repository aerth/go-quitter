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
	"fmt"
	q "github.com/aerth/go-quitter"
	"github.com/aerth/seconf"
	"os"
	"strings"
)

var goquitter = "go-quitter v0.0.7"
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
var hashbar = strings.Repeat("#", 80)
var versionbar = strings.Repeat("#", 10) + "\t" + goquitter + "\t" + strings.Repeat("#", 30)

var usage = "\n" + "\t" + `  Copyright 2016 aerth@sdf.org

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

func init() {
	if gnusocialnode == "" {
		gnusocialnode = "gs.sdf.org"
	}

func bar() {
	print("\033[H\033[2J")
	fmt.Println(versionbar)
}

func main() {
	// list all commands here
	if os.Getenv("GNUSOCIALPATH") != "" {
		gnusocialpath = os.Getenv("GNUSOCIALPATH")
	}
	allCommands := []string{"help", "config", "read", "user", "search", "home", "follow", "unfollow", "post", "mentions", "groups", "mygroups", "join", "leave", "part", "mention", "replies"}

	// command: go-quitter
	if len(os.Args) < 2 {
		bar()
		fmt.Println("Current list of commands:")
		fmt.Println(allCommands)
		fmt.Println(hashbar)
		os.Exit(1)
	}

	if !ContainsString(allCommands, os.Args[1]) {
		bar()
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
			fmt.Println("Your password will NOT echo.")
			seconf.Create(gnusocialpath, "GNU Social", "GNU Social username", "Which GNU Social node? Example: gnusocial.de", "password: will not echo")
		} else {
			bar()
			fmt.Println("Config file already exists.\nIf you want to create a new config file, move or delete the existing one.\nYou can also set the GNUSOCIALPATH env to use multiple config files. \nExample: export GNUSOCIALPATH=gnusocial.de")
			fmt.Println("Config exists:", ReturnHome()+"/."+gnusocialpath)
			os.Exit(1)
		}
	}

	// command: go-quitter help
	helpArg := []string{"help", "halp", "usage", "-help", "-h"}
	if ContainsString(helpArg, os.Args[1]) {
		bar()
		fmt.Println(usage)
		fmt.Println(hashbar)
		os.Exit(1)
	}

	// command: go-quitter version (or -v)
	versionArg := []string{"version", "-v"}
	if ContainsString(versionArg, os.Args[1]) {
		fmt.Println(goquitter)
		os.Exit(1)
	}
	bar()

	// command requires login credentials
	needLogin := []string{"home", "follow", "unfollow", "post", "mentions", "mygroups", "join", "leave", "mention", "replies", "direct", "inbox", "sent"}
	if ContainsString(needLogin, os.Args[1]) {
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
			password = string(configarray[2])
			fmt.Println("Hello, " + username)
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
		ReadPublic(speed)
		os.Exit(0)
	}
	// command: go-quitter search _____
	if os.Args[1] == "search" {
		searchstr := ""
		if len(os.Args) > 1 {
			searchstr = strings.Join(os.Args[2:], " ")
		}
		DoSearch(searchstr, speed)
		os.Exit(0)
	}

	// command: go-quitter user aerth
	if os.Args[1] == "user" && os.Args[2] != "" {
		userlookup := os.Args[2]
		GetUserTimeline(userlookup, speed)
		os.Exit(0)
	}

	// command: go-quitter mentions
	if os.Args[1] == "mentions" || os.Args[1] == "replies" || os.Args[1] == "mention" {
		ReadMentions(speed)
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
		ReadHome(speed)
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
		PostNew(content)
		os.Exit(0)
	}

	// this happens if we invoke with somehing like "go-quitter test"
	fmt.Println(os.Args[0] + " -h")
	os.Exit(1)

}

}
