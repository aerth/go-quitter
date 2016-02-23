# go-quitter
gnusocial client in golang. work in progress. things will break.

## Usage


Default node is gs.sdf.org!

```
#!/bin/sh                                                                       
unset GNUSOCIALNODE                                                             
/go-quitter read >> tweet.log                                                   
GNUSOCIALNODE=gnusocial.de go-quitter read >> treet.log                         
GNUSOCIALNODE=quitter.es go-quitter read >> treet.log                           
GNUSOCIALNODE=shitposter.club go-quitter read >> treet.log                      
GNUSOCIALNODE=sealion.club go-quitter read >> treet.log   

```
