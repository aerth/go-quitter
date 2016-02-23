package main

import (
  "os"
  "log"
  "fmt"
  "io/ioutil"
  "net/http"
)

func main() {

  log.Println("go-quitter v0.0.0")

  if os.Getenv("GNUSOCIALUSER") == "" {
  		fmt.Println("Set environmental variable GNUSOCIALUSER before running go-quitter.")
  		os.Exit(1)
  	}
  if os.Getenv("GNUSOCIALPASS") == "" {
  		fmt.Println("Set environmental variable GNUSOCIALPASS before running go-quitter.")
  		os.Exit(1)
  }

  res, err := http.Post("https://gs.sdf.org/api/statuses/update.xml", "", nil)
	if err != nil {
		log.Fatal(err)
	}
	status, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", status)

}
