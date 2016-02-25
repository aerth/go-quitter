# go-quitter
gnusocial client in golang. work in progress.

## Install
```shell
go get -v -u github.com/aerth/go-quitter

```
## Configure

go-quitter saves a config file at ~/.go-quitter, if it gets messed up just delete it and make a new one.
```
go-quitter config
```
For automation, you may want to use an environmental variable config, something like…

```
// cat ~/.shrc || cat ~/.zshrc || cat ~/.bashrc || cat ~/.whatrc
export GNUSOCIALUSER=yourname
export GNUSOCIALPASS=yourpass

// then run this command so you dont have to log out and back in.
. ~/.shrc

// make sure you chmod your shell rc file if shared machine.
chmod o-r ~/.shrc
chmod g-r ~/.shrc

## Usage

```shell

$ go-quitter read // ticker style public timeline
$ go-quitter read fast // reads public timeline
$ go-quitter user aerth // looks up a user timeline
$ go-quitter home fast // reads your home timeline
$ go-quitter post posting totally works!
$ go-quitter post // this presents a prompt
$ go-quitter post \!group \#hashtag \#EscapeSymbolsWithABackslash

```

Default node is gs.sdf.org!

```shell
#!/bin/sh                                                                       
unset GNUSOCIALNODE                                                             
go-quitter read fast >> tweet.log                                                   
GNUSOCIALNODE=gnusocial.de go-quitter read fast >> treet.log                         
GNUSOCIALNODE=quitter.es go-quitter read fast >> treet.log                           
GNUSOCIALNODE=shitposter.club go-quitter fast read >> treet.log                      
GNUSOCIALNODE=sealion.club go-quitter read fast >> treet.log   

```

### Todo

* include user interface with up/down scrolling
* ~~get simple posting to work~~
* write tests
* ~~save account information in encoded config file~~
* cat filename.txt | go-quitter // may do this just because it would make uploading photos easy.
* **learn go**



### Contributing

* Pull requests are welcome.
