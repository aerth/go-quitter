// Minimal quitter library example
package main

/*

Copyright 2016 aerth

This is a simple example program showing
the usage of the quitter (GNU Social) library.

You can import the libary as whatever name you want.
Default is quitter, here it has an alias of qw.

*/
import (
	"fmt"
	"os"

	qw "github.com/aerth/go-quitter"
)

func main() {

	if len(os.Args) == 2 && os.Args[1] == "home" {
		q := qw.NewAccount()
		q.Username = "username"
		q.Password = "nopassword"
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
		q2 := qw.NewAccount()
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

	} else {

		// Example usage
		fmt.Println(os.Args[0], "public")
		fmt.Println(os.Args[0], "home")
	}
}
