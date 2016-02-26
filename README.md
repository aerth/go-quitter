# go-quitter
gnusocial client in golang. work in progress.

## Install binary for your OS
[Linux](https://github.com/aerth/go-quitter/releases/download/v0.0.4/go-quitter-v0.0.4_linux-amd64.tgz)
```
MD5:    50b4670858c3570a344a4a8a45c10703  go-quitter-v0.0.4_linux-amd64.tgz
SHA1:   a8bfcf73a2ceb10442436f58a3c5dc4852122168  go-quitter-v0.0.4_linux-amd64.tgz
SHA256: 7b51bdd663ce1101da0910da8946d7f46a50a8cf600efd53afe7f05d36a42d8d  go-quitter-v0.0.4_linux-amd64.tgz
```

[NetBSD](https://github.com/aerth/go-quitter/releases/download/v0.0.4/go-quitter-v0.0.4_netbsd-amd64.tar.gz)
```
MD5:    b019020684485805994205e52e525428  go-quitter-v0.0.4_netbsd-amd64.tar.gz
SHA1:   d135d6a448c9978c869d7f6c9b8be26eae4d2030  go-quitter-v0.0.4_netbsd-amd64.tar.gz
SHA256: e2647e110ef2d8fd22f1aedf3f776351384577ff0443a1c403b856ff3a57d91e  go-quitter-v0.0.4_netbsd-amd64.tar.gz
```

## Install from Go source
```shell
go get -v -u github.com/aerth/go-quitter

```
## Configure

go-quitter saves a config file at ~/.go-quitter, if it gets messed up just delete it and make a new one.
```
go-quitter config
```
For automation and cronjobs, you may want to use environmental variables instead. Something like…

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
GNUSOCIALNODE=shitposter.club go-quitter read fast >> treet.log                      
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
