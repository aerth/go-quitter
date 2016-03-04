package main

import (
  qw "github.com/aerth/go-quitter"
  "os"
)
func main() {

  if len(os.Args) < 2 {
  q := qw.NewAuth()
  q.Username = "john"
  q.Password = "pass123"
  q.Node = "gnusocial.de"
  q.ReadHome(false)
  // Return: Could not authenticate you.
  }

  if os.Args[1] == "public" {
  q2 := qw.NewAuth()
  q2.Node = "gnusocial.de"
  q2.ReadPublic(true)
  // Return: Public timeline on gnusocial.de
  }

}
