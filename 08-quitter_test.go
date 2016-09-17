/*
Some example functions for godoc
*/

package quitter_test

import (
	"fmt"
	"os"

	quitter "github.com/aerth/go-quitter"
)

func Example() {
	q := quitter.NewSocial()
	q.Username = "username"
	q.Password = "password"
	q.Node = "localhost"
	quips, err := q.GetHome()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, quip := range quips {
		fmt.Printf("%s %s", quip.IdStr, quip.Text)
	}
}

func ExampleNewSocial() {
	q := quitter.NewSocial()
	q.Username = "username"
	q.Password = "password"
	q.Node = "localhost"
	quips, err := q.GetHome()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, quip := range quips {
		fmt.Printf("%s %s", quip.IdStr, quip.Text)
	}
}
func ExampleSocial_PostNew() {
	q := quitter.NewSocial()
	q.Username = "username"
	q.Password = "password"
	q.Node = "localhost"
	content := ` dang this is a " < new > ! quip about to be published>>><><><?><?><?><<?><?><?><?><?><?><?><?><"`
	quip, err := q.PostNew(content)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s %s", quip.IdStr, quip.Text)
}

func ExampleSocial_GetPublic() {
	q := quitter.NewSocial()
	q.Username = "username"
	q.Password = "password"
	q.Node = "localhost"
	quips, err := q.GetPublic()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, quip := range quips {
		fmt.Printf("%s %s", quip.IdStr, quip.Text)
	}

}
