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

	"github.com/aerth/go-quitter" // libgoquitter
)

var (
	release       = "v0.0.9 (go get)"
	buildinfo     string
	goquitter     = "go-quitter " + release
	username      = os.Getenv("GNUSOCIALUSER")
	password      = os.Getenv("GNUSOCIALPASS")
	gnusocialnode = os.Getenv("GNUSOCIALNODE")
	gnusocialpath = os.Getenv("GNUSOCIALPATH")
	apipath       = "https://" + gnusocialnode + "/api/statuses/home_timeline.json"
	builtWithCUI  = false
	initgui       func()
	configuser    string
	configpass    string
	confignode    string
	configlock    string
	configstrings string
	hashbar       string
)

func init() {
	if gnusocialpath == "" {
		gnusocialpath = "go-quitter"
	}
}

var versionbar = goquitter + " built " + buildinfo

var usage = `
config		Creates config file	*do this first*
read			Reads 20 new posts
read fast		Reads 20 new posts (no delay)
home			Your home timeline.
user username	Looks up "username" timeline
post ____ 		Posts to your node.
post 		Post mode.
mentions		Mentions your @name
search ___		Searches for ____
search		Search mode.
follow		Follow a user
unfollow		Unfollow a user
groups		List all groups on current node
mygroups		List only groups you are member of
join ___		Join a !group
leave ___		Part a !group (can also use part)


Using environmental variables will override the config:

GNUSOCIALPATH - path to config file (default ~/.go-quitter)
GNUSOCIALNODE, GNUSOCIALPASS, GNUSOCIALUSER - account info

Want to use a SOCKS proxy?
Set the SOCKS environmental variable. Here are a few examples:

	SOCKS=true go-quitter -socks # short for 127.0.0.1:1080
	SOCKS=tor go-quitter -socks # short for 127.0.0.1:9050
	SOCKS=socks5://127.0.0.1:22000 go-quitter -socks

[FLAGS] can be placed before a command. Here are the available flags:

	-socks Don't connect without proxy
	-http Don't use https
	-unsafe Don't validate TLS cert
`

var q *quitter.Account

var allCommands = []string{"version", "help", "config",
	"read", "user", "search", "home", "follow", "unfollow",
	"post", "mentions", "groups", "mygroups", "join", "leave",
	"replies", "upload"}
var needLogin = []string{"home", "follow", "unfollow",
	"post", "mentions", "mygroups", "groups", "search", "join", "leave", "mention",
	"replies", "direct", "inbox", "sent", "upload", "cui"}

func print(s string) {
	fmt.Fprint(os.Stderr, s)
}

func printf(f string, i ...interface{}) {
	fmt.Fprintf(os.Stderr, f, i...)
}

// flagy can transcend space and time
func flagy(a []string) []string {

	//	-unsafe flag does not validate TLS certs
	if containsString(a, "-unsafe") {
		// -unsafe: remove -unsafe arg
		a = func(old []string) (new []string) {
			for i := range old {
				if old[i] == "-unsafe" {
					continue
				}
				new = append(new, old[i])
			}
			return new
		}(a)
		// -unsafe: warn user on stderr
		q.Scheme = "http://"
		printf("*Using %q scheme*\n", q.Scheme)
	}

	// -http flag uses http instead of https
	if containsString(a, "-http") {
		// -http: remove -http arg
		a = func(old []string) (new []string) {
			for i := range old {
				if old[i] == "-http" {
					continue
				}
				new = append(new, old[i])
			}
			return new
		}(a)
		// -unsafe: warn user on stderr
		quitter.EnableInvalidTLS()
		print("*Not validating TLS certificates*\n")
	}
	//	-socks flag MAKE SURE we are using a socks proxy.
	// it must be configured using SOCKS environmental variable or new config
	if containsString(a, "-socks") {
		// -socks: remove -socks arg
		a = func(old []string) (new []string) {
			for i := range old {
				if old[i] == "-socks" {
					continue
				}
				new = append(new, old[i])
			}
			return new
		}(a)
		// -socks: warn user on stderr
		if quitter.ProxyString == "" {
			fmt.Println("No proxy")
			os.Exit(1)
		}
	}
	return a
}

func main() {
	args := os.Args
	q = quitter.NewAccount()

	args = flagy(args)

	if len(args) < 2 || !containsString(allCommands, args[1]) {
		fmt.Println(versionbar)
		fmt.Println("Commands:", allCommands)
		os.Exit(1)
	}

	if args[1] == "config" {
		makeConfig()
		os.Exit(0)
	}

	// command: go-quitter help
	helpArg := []string{"help", "halp", "usage", "-help", "-h"}
	if containsString(helpArg, args[1]) {
		fmt.Println(goquitter, buildinfo)
		fmt.Println(usage)
		fmt.Println(hashbar)
		os.Exit(0)
	}

	// command: go-quitter version (or -v)
	versionArg := []string{"version", "-v"}
	if containsString(versionArg, args[1]) {
		fmt.Println(goquitter, buildinfo)
		os.Exit(0)
	}

	// command requires login credentials

	if containsString(needLogin, args[1]) {
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

	switch args[1] {
	// command: go-quitter read
	case "cui":
		initgui()
		os.Exit(0)

	case "read":
		PrintQuips(q.GetPublic())
		os.Exit(0)

		// command: go-quitter search _____
	case "search":
		searchstr := ""
		if len(args) > 1 {
			searchstr = strings.Join(args[2:], " ")
		}
		if searchstr == "" {
			searchstr = getTypin()
		}
		PrintQuips(q.Search(searchstr))
		os.Exit(0)

		// command: go-quitter user aerth
	case "user":
		if len(args) > 2 && args[2] != "" {
			userlookup := args[2]
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
		if len(args) == 1 {
			followstr = args[2]
		} else if len(args) > 1 {
			followstr = strings.Join(args[2:], " ")
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
		if len(args) == 1 {
			followstr = args[2]
		} else if len(args) > 1 {
			followstr = strings.Join(args[2:], " ")
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
		if len(args) > 1 {
			groupstr = strings.Join(args[2:], " ")
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
		if len(args) > 1 {
			groupstr = strings.Join(args[2:], " ")
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
		if len(args) > 1 {
			content = strings.Join(args[2:], " ")
		}
		PrintGroup(q.PartGroup(content))
		os.Exit(0)

		// command: go-quitter post
	case "post":
		var content string
		if len(args) > 1 {
			content = strings.Join(args[2:], " ") // go-quitter post wow this is a post\!
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

	// command: go-quitter upload
	case "upload":
		var path, content string // go-quitter upload cat.gif lol
		if len(args) > 1 {
			path = args[2] // cat.gif
		}
		if path == "" {
			path = getTypin()
		}
		if len(args) > 2 {
			content = strings.Join(args[3:], " ") // lol
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
		fmt.Println("Command not found, try ", args[0]+" help")
		os.Exit(1)
	}
}

func init() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, versionbar)
		fmt.Fprintln(os.Stderr, "Current list of commands:")
		fmt.Fprintln(os.Stderr, allCommands)
		fmt.Fprintf(os.Stderr, "\nRun '%s -help' for more information.\n\n", os.Args[0])
		os.Exit(1)
	}
	if os.Getenv("SOCKS") != "" {
		quitter.EnableSOCKS(os.Getenv("SOCKS"))
	}
}
