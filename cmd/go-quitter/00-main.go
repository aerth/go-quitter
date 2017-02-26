// go-quitter is a console GNU Social client.
package main

/*
The MIT License (MIT)

Copyright (c) 2016 aerth

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
*/

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aerth/go-quitter"
)

var (
	goquitter     = "go-quitter v0.0.9"
	username      = os.Getenv("GNUSOCIALUSER")
	password      = os.Getenv("GNUSOCIALPASS")
	gnusocialnode = os.Getenv("GNUSOCIALNODE")
	apipath       = "https://" + gnusocialnode + "/api/statuses/home_timeline.json"
	gnusocialpath = "go-quitter"
	configuser    string
	configpass    string
	confignode    string
	configlock    string
	configstrings string
	hashbar       string
)

var versionbar = strings.Repeat("#", 10) + "\t" + goquitter + "\t" + strings.Repeat("#", 30)

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
GNUSOCIALPATH - path to config file (default ~/.go-quitter)

Set your environmental variable to change nodes, use a different config,
	or change user or password for a one-time session.

Want to use a SOCKS proxy? Set the SOCKS environmental variable. Here is an example:

	SOCKS=socks5://127.0.0.1:1080 ./go-quitter
`

var q *quitter.Account
var allCommands = []string{"help", "config",
	"read", "user", "search", "home", "follow", "unfollow",
	"post", "mentions", "groups", "mygroups", "join", "leave",
	"replies", "gui-test", "upload"}

func main() {

	q = quitter.NewAccount()
	q.Proxy = os.Getenv("SOCKS")
	if containsString(os.Args, "-debug") {
		q.Scheme = "http://"
	}

	if os.Args[1] == "config" {
		makeConfig()
		os.Exit(0)
	}

	// command: go-quitter help
	helpArg := []string{"help", "halp", "usage", "-help", "-h"}
	if containsString(helpArg, os.Args[1]) {
		fmt.Println(usage)
		fmt.Println(hashbar)
		os.Exit(0)
	}

	// command: go-quitter version (or -v)
	versionArg := []string{"version", "-v"}
	if containsString(versionArg, os.Args[1]) {
		fmt.Println(goquitter)
		os.Exit(0)
	}

	// command requires login credentials
	needLogin := []string{"gui-test", "home", "follow", "unfollow",
		"post", "mentions", "mygroups", "groups", "search", "join", "leave", "mention",
		"replies", "direct", "inbox", "sent", "upload"}
	if containsString(needLogin, os.Args[1]) {
		needConfig()
	} else { // command doesn't need login
		if configExists() {
			dontNeedConfig()
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

	switch os.Args[1] {
	// command: go-quitter read
	case "test":
		//		runtests()
	case "gui-test":

		initgui()
		os.Exit(0)

	case "read":
		PrintQuips(q.GetPublic())
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
		PrintQuips(q.Search(searchstr))
		os.Exit(0)

		// command: go-quitter user aerth
	case "user":
		if len(os.Args) > 2 && os.Args[2] != "" {
			userlookup := os.Args[2]
			PrintQuips(q.GetUserTimeline(userlookup))

			os.Exit(0)
		}
		fmt.Println("Need user to search for")
		os.Exit(1)

		// command: go-quitter mentions
	case "mentions", "replies", "mention":
		PrintQuips(q.GetMentions())
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
		PrintUser(q.Follow(followstr))
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
		PrintUser(q.UnFollow(followstr))
		os.Exit(0)

	// command: go-quitter home
	case "home":
		PrintQuips(q.GetHome())
		os.Exit(0)

	// command: go-quitter groups
	case "groups":
		PrintGroups(q.ListAllGroups())
		os.Exit(0)

		// command: go-quitter mygroups
	case "mygroups":
		PrintGroups(q.ListMyGroups())
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

		// command: go-quitter post
	case "post":
		var content string
		if len(os.Args) > 1 {
			content = strings.Join(os.Args[2:], " ") // go-quitter post wow this is a post\!
		}
		if content == "" {
			content = getTypin()
		}
		// go-quitter post -y hello world
		if !strings.HasPrefix(content, "-y ") {
			fmt.Println("Preview:\n\n[" + q.Username + "] " + content)
			fmt.Println("\nType YES to publish!")
			if askForConfirmation() == false {
				fmt.Println("Your status was not updated.")
				os.Exit(0)
			}
		} else {
			content = strings.TrimPrefix(content, "-y ")
		}

		PrintQuip(q.PostNew(content))

		os.Exit(0)
		// command: go-quitter post
	case "upload":
		var path, content string // go-quitter upload cat.gif lol
		if len(os.Args) > 1 {
			path = os.Args[2] // cat.gif
		}
		if path == "" {
			path = getTypin()
		}
		if len(os.Args) > 2 {
			content = strings.Join(os.Args[3:], " ") // lol
		}
		if content == "" {
			content = getTypin()
		}
		fmt.Printf("Uploading %q", path)
		if content != "" {
			fmt.Printf(" with caption %q", content)
		}
		fmt.Println()
		time.Sleep(time.Second)
		PrintQuip(q.Upload(path, content))
	default:
		// this happens if we invoke with somehing like "go-quitter test"
		fmt.Println("Command not found, try ", os.Args[0]+" help")
		os.Exit(1)
	}
}

var initgui = func() { fmt.Println("go-quitter not built with cui support.") }

func init() {

	if len(os.Args) < 2 {

		fmt.Println("\n\n" + versionbar)
		fmt.Println("Current list of commands:")
		fmt.Println(allCommands)
		fmt.Printf("\nRun '%s -help' for more information.\n\n", os.Args[0])

		os.Exit(0)
	}
	if os.Getenv("GNUSOCIALPATH") != "" {
		gnusocialpath = os.Getenv("GNUSOCIALPATH")
	}
	if gnusocialnode == "" {
		gnusocialnode = "gnusocial.de"
	}
	if !containsString(allCommands, os.Args[1]) {
		fmt.Println(usage)
		fmt.Println("Current list of commands:")
		fmt.Println(allCommands)
		fmt.Println(hashbar)
		os.Exit(1)
	}
}
