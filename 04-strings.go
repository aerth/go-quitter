package quitter

import "fmt"

func (q Quip) String() string {
	if q.User.Screenname == q.User.Name {
		return fmt.Sprintf("[@" + q.User.Screenname + "] " + q.Text + "\n")
	}
	return fmt.Sprintf("@" + q.User.Screenname + " [" + q.User.Name + "] " + q.Text + "\n")
}

func (group Group) String() string {
	return fmt.Sprintf("!" + group.Nickname + " [" + group.Fullname + "] " + group.Description + "\n")
}

func (u User) String() string {
	return fmt.Sprintf("@" + u.Screenname)
}
func (a Account) String() string {
	str := "User: "+a.Username+"\n"
  str += "Node: "+a.Node+"\n"
if a.Proxy != "" {
  str += "Proxy: "+a.Proxy+"\n"
}
  return str
}
