# go-quitter
gnusocial client in golang. work in progress.

## Install binary for your OS
[Linux](https://github.com/aerth/go-quitter/releases/download/v0.0.5/go-quitter-v0.0.5_linux-amd64.tgz)
```
MD5:    e965aec65cfaab367bafd65b80fe91fe  go-quitter-v0.0.5_linux-amd64.tgz
SHA1:   2c452cbcc71706a5ce5969f3e0be4901a36c6e9c  go-quitter-v0.0.5_linux-amd64.tgz
SHA256: cb3b749d396d2ffc91754b0bef4f569d901d8f847132baa908221fc2e67cef15  go-quitter-v0.0.5_linux-amd64.tgz
MD5:    10a1912a3a1949750a12eb49745dbe4b  go-quitter/go-quitter
SHA1:   7d5c019bfd2a0579bd46872daeca431f0b843bb9  go-quitter/go-quitter
SHA256: 0782c3d0451456ab06b3c50c2cd61fd07af4f51faf7a64cce7c1812de1ce8ff8  go-quitter/go-quitter
```

[NetBSD](https://github.com/aerth/go-quitter/releases/download/v0.0.5/go-quitter-v0.0.5_netbsd-amd64.tgz)
```
MD5:    e401e0c85d71bee2efae6309222c337c  go-quitter/go-quitter
SHA1:   ddba90e508873d5ea0297c129d083b3588ddd2d4  go-quitter/go-quitter
SHA256: 64ba3d3caa36f8c9212af10233e4cba096e2adfb0d4230a0ccf09954b5df33be  go-quitter/go-quitter
MD5:    ab4eb1af7c6e2f4cf7b3947f19cae24b  go-quitter-v0.0.5_netbsd-amd64.tgz
SHA1:   ebb0a47458e4f97972fe79f038e7ef7945c61ffd  go-quitter-v0.0.5_netbsd-amd64.tgz
SHA256: 0130cd3a6d9ab6dbc5fe6ac68e85f708ac38e915f6fe1dc3e388eb2ae9354320  go-quitter-v0.0.5_netbsd-amd64.tgz
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
For automation, scripts, and cronjobs, you may want to delete config file and use environmental variables instead. Something like…

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
