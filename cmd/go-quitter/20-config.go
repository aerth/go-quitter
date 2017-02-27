package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aerth/seconf"
)

func init(){
	if gnusocialpath == "" {
		gnusocialpath = "go-quitter"
	}
}

func makeConfig() {

	if seconf.Detect(gnusocialpath) == false {
		fmt.Println(versionbar)
		fmt.Printf("New config: %q\nYou will be asked for your user, node, and password.\n", returnHomeDir()+"/."+gnusocialpath)
		fmt.Println("You must have a GNU Social account already. To sign up, find a node!")
		fmt.Println("More info here: https://gnu.io/social/try/servers.html")
		fmt.Println("Your password will not echo.")
		seconf.Create(gnusocialpath, "GNU Social", "GNU Social username", "Which GNU Social node? Example: gnusocial.de", "password: will not echo")
	} else {

		fmt.Println("Config file already exists.\nIf you want to create a new config file, move or delete the existing one.\nYou can also set the GNUSOCIALPATH env to use multiple config files. \nExample to use ~/.gs-2 config:\n\t'GNUSOCIALPATH=gs-2 go-quitter config'")
		fmt.Println("\nConfig exists:", returnHomeDir()+"/."+gnusocialpath)
		os.Exit(1)
	}
}

func dontNeedConfig() {
	//fmt.Println("Config file detected, but this command doesn't need to login.\nWould you like to select the GNU Social node using the config?\nType YES or NO (y/n)")
	//if AskForConfirmation() == true {
	// only use gnusocial node from config
	configdecoded, err := seconf.Read(gnusocialpath)
	if err != nil {
		fmt.Println("error:")
		fmt.Println(err)
		os.Exit(1)
	}
	configarray := strings.Split(configdecoded, "::::")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(configarray) != 3 {
		fmt.Println("Broken config file. Create a new one.")
		os.Exit(1)
	}
	gnusocialnode = string(configarray[1])
	//fmt.Println(gnusocialnode)
	gnusocialnode = strings.Replace(gnusocialnode, "http://", "", -1)
	gnusocialnode = strings.Replace(gnusocialnode, "https://", "", -1)
	//			gnusocialnode = strings.TrimLeft(gnusocialnode, "https://")
	//			gnusocialnode = strings.TrimLeft(gnusocialnode, "http://")
	q.Node = gnusocialnode

}
func configExists() bool {
	return seconf.Detect(gnusocialpath)
}
func needConfig() {

	if seconf.Detect(gnusocialpath) == true {
		configdecoded, err := seconf.Read(gnusocialpath)
		if err != nil {
			fmt.Println("error:")
			fmt.Println(err)
			os.Exit(1)
		}
		//configstrings := string(configdecoded)
		//		fmt.Println("config strings:")
		//		fmt.Println(configdecoded)
		configarray := strings.Split(configdecoded, "::::")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(configarray) != 3 {
			fmt.Println("Broken config file. Create a new one. :(")
			os.Exit(1)
		}
		username = string(configarray[0])
		gnusocialnode = string(configarray[1])
		gnusocialnode = strings.Replace(gnusocialnode, "http://", "", -1)
		gnusocialnode = strings.Replace(gnusocialnode, "https://", "", -1)
		password = string(configarray[2])

		q.Username = username
		q.Password = password
		q.Node = gnusocialnode
		if os.Getenv("GNUSOCIALUSER") != "" {
			q.Username = os.Getenv("GNUSOCIALUSER")
		}
		if os.Getenv("GNUSOCIALPASS") != "" {
			q.Password = os.Getenv("GNUSOCIALPASS")
		}
		if os.Getenv("GNUSOCIALNODE") != "" {
			q.Node = os.Getenv("GNUSOCIALNODE")
		}
		if q.Username == username {
			fmt.Fprintln(os.Stderr, "Welcome back, " + q.Username + "@" + q.Node)
		} else {
			fmt.Fprintln(os.Stderr, "Welcome, " + q.Username + "@" + q.Node)
		}
	} else {
		fmt.Fprintln(os.Stderr, "No config file detected at", gnusocialpath)
	}
}
