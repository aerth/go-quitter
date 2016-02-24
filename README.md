# go-quitter
gnusocial client in golang. work in progress. things will break.

## Usage

```

$ go-quitter read // ticker style public timeline
$ go-quitter read fast // reads public timeline
$ go-quitter user aerth // looks up a user timeline
$ go-quitter home fast // reads your home timeline
$ go-quitter post "posting doesn't work yet"

```

Default node is gs.sdf.org!

```
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
* get simple posting to work
* write tests
* save account information in encoded config file


### Contributing

* Pull requests are welcome.
