package main

import (
	"bufio"
	"fmt"
	"os"
)

// Receive non-hidden input from user.
func getTypin() string {
	fmt.Printf("\nPress ENTER when you are finished typing.\n\n")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text()
		//	fmt.Println(line)
		return line
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return ""
}


// Ask user to confirm the action.
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	quitResponses := []string{"q", "Q", "exit", "quit"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else if containsString(quitResponses, response) {
		return false
	} else {
		fmt.Println("\nNot valid answer, try again. [y/n] [yes/no]")
		return askForConfirmation()
	}
}

// Does []string contain element?
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

// Find the index of a string in a []string
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}



// returnHomeDir gives us the true home directory for letting the user know where the config file is. Windows, Unix, OS X
func returnHomeDir() (homedir string) {
	homedir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if homedir == "" {
		homedir = os.Getenv("USERPROFILE")
	}
	if homedir == "" {
		homedir = os.Getenv("HOME")
	}
	return homedir
}


func PrintQuips(quips []quitter.Quip, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(quips) == 0 && err == nil {
		fmt.Println("No results.")
		return
	}

	for i, j := 0, len(quips)-1; i < j; i, j = i+1, j-1 {
		quips[i], quips[j] = quips[j], quips[i]
	}
	for i := range quips {
		if quips[i].User.Screenname == quips[i].User.Name {
			fmt.Printf("[@" + quips[i].User.Screenname + "] " + quips[i].Text + "\n\n")
		} else {
			fmt.Printf("@" + quips[i].User.Screenname + " [" + quips[i].User.Name + "] " + quips[i].Text + "\n\n")
		}
	}
}
func PrintQuip(quip quitter.Quip, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	if quip.Text == "" && err == nil {
		fmt.Println("No quip.")
		return
	}
	if quip.User.Screenname == quip.User.Name {
		fmt.Printf("[@" + quip.User.Screenname + "] " + quip.Text + "\n\n")
	} else {
		fmt.Printf("@" + quip.User.Screenname + " [" + quip.User.Name + "] " + quip.Text + "\n\n")
	}

}

func PrintUsers(users []quitter.User, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(users) == 0 && err == nil {
		fmt.Println("No users.")
		return
	}
	for i := range users {
		if users[i].Screenname == users[i].Name {
			fmt.Printf("[@" + users[i].Screenname + "]\n\n")
		} else {
			fmt.Printf("@" + users[i].Screenname + " [" + users[i].Name + "]\n\n")
		}
	}
}

//PrintUser prints a single @user
func PrintUser(user quitter.User, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	if user.Screenname == "" && err == nil {
		fmt.Println("No user.")
		return
	}
	fmt.Printf("[@" + user.Screenname + "]\n\n")

}
func PrintGroup(group quitter.Group, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	if group.Nickname == "" && err == nil {
		fmt.Println("No group.")
		return
	}
	fmt.Printf("!" + group.Nickname + " [" + group.Fullname + "] \n" + group.Description + "\n\n")

}

func PrintGroups(groups []quitter.Group, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(groups) == 0 && err == nil {
		fmt.Println("No groups.")
		return
	}
	for i := range groups {
		fmt.Printf("!" + groups[i].Nickname + " [" + groups[i].Fullname + "] \n" + groups[i].Description + "\n\n")
	}
}
