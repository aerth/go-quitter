# go-quitter
gnusocial client in golang. work in progress. things will break.

## Usage

```

$ go-quitter read

or

$ go-quitter read fast

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
