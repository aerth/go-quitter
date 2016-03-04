package main

import (
	qw "github.com/aerth/go-quitter"
	"os"
  "fmt"
)

func main() {

	if len(os.Args) == 2 && os.Args[1] == "home" {
		q := qw.NewAuth()
		q.Username = "john"
		q.Password = "pass123"
		q.Node = "gnusocial.de"
		quips, err := q.GetHome(false)
    if err != nil {
      fmt.Println(err)
    }
    for i := range quips {
  		if quips[i].User.Screenname == quips[i].User.Name {
  			fmt.Printf("[@" + quips[i].User.Screenname + "] " + quips[i].Text + "\n\n")
  		} else {
  			fmt.Printf("@" + quips[i].User.Screenname + " [" + quips[i].User.Name + "] " + quips[i].Text + "\n\n")
  		}
  	}
    os.Exit(1)
		// Return: Could not authenticate you.
	}

	if len(os.Args) == 2 && os.Args[1] == "public" {
		q2 := qw.NewAuth()
		q2.Node = "gnusocial.de"
		quips, err := q2.GetPublic(true)
    if err != nil {
      fmt.Println(err)
    }
    for i := range quips {
  		if quips[i].User.Screenname == quips[i].User.Name {
  			fmt.Printf("[@" + quips[i].User.Screenname + "] " + quips[i].Text + "\n\n")
  		} else {
  			fmt.Printf("@" + quips[i].User.Screenname + " [" + quips[i].User.Name + "] " + quips[i].Text + "\n\n")
  		}
  	}
    os.Exit(1)


	}else {

  // Example usage
  fmt.Println(os.Args[0], "public")
  fmt.Println(os.Args[0], "home")
}
}
