# go-quitter

Command line **GNU Social** client and Go library

[![GoDoc](https://godoc.org/github.com/aerth/go-quitter?status.svg)](https://godoc.org/github.com/aerth/go-quitter)

```shell go-quitter help more```

```
Usage: go-quitter [command]
config         Creates config file	*do this first*
read           Reads 20 new posts
home           Your home timeline.
user ____      Looks up "username" timeline
post ____      Posts to your node.
post           Post mode.
mentions       Mentions your @name
search ___     Searches for ____
search         Search mode.
follow         Follow a user
unfollow       Unfollow a user
groups         List all groups on current node
mygroups       List only groups you are member of
join ___       Join a !group
leave ___      Part a !group (can also use part)

* Using environmental variables will override the config:

GNUSOCIALPATH - path to config file (default ~/.go-quitter)
GNUSOCIALNODE, GNUSOCIALPASS, GNUSOCIALUSER - account info

* Want to use a SOCKS proxy?
Set the SOCKS environmental variable. Here are a few examples:

	SOCKS=true go-quitter -socks # short for 127.0.0.1:1080
	SOCKS=tor go-quitter -socks # short for 127.0.0.1:9050
	SOCKS=socks5://127.0.0.1:22000 go-quitter -socks

* -flags can be placed before a [command]. Here are the available flags:

	-socks Don't connect without proxy
	-http Don't use https
	-unsafe Don't validate TLS cert

Check for updates: https://github.com/aerth/go-quitter



################################################################################


```

## Install (outdated) binary for your OS
### [Latest Binary Releases](https://github.com/aerth/go-quitter/releases/latest)

## Install from Go source (sometimes newer)

If you have Go toolchain installed you can build it yourself with:

```shell
GOPATH=/tmp/go go get -v -u -d github.com/aerth/go-quitter/cmd/go-quitter
cd $GOPATH/src/github.com/aerth/go-quitter/cmd/go-quitter
make && su -c 'make install'

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

  * CUI

### Contributing

* Pull requests are welcome.
