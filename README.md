# go-quitter

Command line **GNU Social** client and Go library

[![GoDoc](https://godoc.org/github.com/aerth/go-quitter?status.svg)](https://godoc.org/github.com/aerth/go-quitter)

```go 

$ go-quitter help

##########	go-quitter v0.0.8	##############################

	  Copyright 2016 aerth@sdf.org

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



################################################################################


```

## Install binary for your OS (now for every OS)
### [Latest Binary Releases](https://github.com/aerth/go-quitter/releases)
[Windows](https://github.com/aerth/go-quitter/releases/download/v0.0.7/go-quitter-v0.0.7.exe) | 
[Linux](https://github.com/aerth/go-quitter/releases/download/v0.0.6/go-quitter-v0.0.6_linux-amd64.tar.gz) | 
[FreeBSD](https://github.com/aerth/go-quitter/releases/download/v0.0.6/go-quitter-v0.0.6_linux-amd64.tar.gz) | 
[NetBSD](https://github.com/aerth/go-quitter/releases/download/v0.0.6/go-quitter-v0.0.6_linux-amd64.tar.gz) | 
[OS X](https://github.com/aerth/go-quitter/releases/download/v0.0.6/go-quitter-v0.0.6_linux-amd64.tar.gz) | 
[Other OS or Architecture](https://github.com/aerth/go-quitter/releases)

## Install from Go source (sometimes newer)

If you have Go toolchain installed you can build it yourself with:

```shell
GOPATH=/tmp/go go get -v -u github.com/aerth/go-quitter/cmd/go-quitter
cd $GOPATH/src/github.com/aerth/go-quitter/cmd/go-quitter
make && sudo make install
```

## Go Get-able

Or use go get:

```
go get -v -u github.com/aerth/go-quitter/cmd/go-quitter

```


## Configure

To avoid storing the password in plaintext, go-quitter saves an encrypted config file at ~/.go-quitter, if it gets messed up just delete it and make a new one. You can switch config files on the fly using the environmental variable GNUSOCIALPATH.

```
go-quitter config
GNUSOCIALPATH=gnusocial.de go-quitter config
GNUSOCIALPATH=gnusocial.no go-quitter config
GNUSOCIALPATH=gnusocial.se go-quitter config

```

Next time you run it, it will ask for your config password. I like to keep it blank so I just hit ENTER.


## Or.. Use in scripts

For automation, scripts, and cronjobs, you *need* to delete config file and use environmental variables instead. Something like…

```
// cat ~/.shrc || cat ~/.zshrc || cat ~/.bashrc || cat ~/.whatrc
export GNUSOCIALUSER=yourname
export GNUSOCIALPASS=yourpass
export GNUSOCIALNODE=gnusocial.de

// then run this command so you dont have to log out and back in.
. ~/.shrc


// make sure you chmod your shell rc file if shared machine.
chmod o-r ~/.shrc
chmod g-r ~/.shrc

## Usage

When running go-quitter with no arguments, a list of commands is printed.
For more information, run `go-quitter help`

```shell

$ go-quitter read // public timeline
$ go-quitter home // home timeline
$ go-quitter search // enters search mode
$ go-quitter post \!group \#hashtag \#EscapeSymbolsWithABackslash
```

```shell
#!/bin/sh                                                                       
unset GNUSOCIALNODE                                                          
GNUSOCIALNODE=gnusocial.de go-quitter read fast >> treet.log                         
GNUSOCIALNODE=quitter.es go-quitter read fast >> treet.log                           
GNUSOCIALNODE=shitposter.club go-quitter read fast >> treet.log                      
GNUSOCIALNODE=sealion.club go-quitter read fast >> treet.log   

```

### Todo

- [ ] Learn Go
- [ ] Include user interface with up/down scrolling
- [ ] Write tests
- [ ] cat filename.txt | go-quitter // I may introduce this just because it would make uploading photos easy.



### Contributing

* Pull requests are welcome.
