package quitter

// Social is credentials needed for logging in. Set it with NewSocial()
type Social struct {
	Username string
	Password string
	Node     string
	Scheme   string
}

// Sets the authentication method and choose node.
//Use like this:
/*

 q := qw.NewSocial()
 q.Username = "john"
 q.Password = "pass123"
 q.Node = "gnusocial.de"
 q.GetHome(false)

*/
func NewSocial() *Social {
	return &Social{
		Username: "gopher",
		Password: "password",
		Node:     "localhost",
		Scheme:   "https://",
	}
}

// User is a GNU Social User, gets returned by GS API
type User struct {
	Name       string `json:"name"`
	Screenname string `json:"screen_name"`
}

// Quip is a GNU Social Quip, gets returned by GS API.
type Quip struct {
	Id                   int64    `json:"id"`
	IdStr                string   `json:"id_str"`
	InReplyToScreenName  string   `json:"in_reply_to_screen_name"`
	InReplyToStatusID    int64    `json:"in_reply_to_status_id"`
	InReplyToStatusIdStr string   `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64    `json:"in_reply_to_user_id"`
	InReplyToUserIdStr   string   `json:"in_reply_to_user_id_str"`
	Lang                 string   `json:"lang"`
	Place                string   `json:"place"`
	PossiblySensitive    bool     `json:"possibly_sensitive"`
	RetweetCount         int      `json:"retweet_count"`
	Retweeted            bool     `json:"retweeted"`
	RetweetedStatus      *Quip    `json:"retweeted_status"`
	Source               string   `json:"source"`
	Text                 string   `json:"text"`
	Truncated            bool     `json:"truncated"`
	User                 User     `json:"user"`
	WithheldCopyright    bool     `json:"withheld_copyright"`
	WithheldInCountries  []string `json:"withheld_in_countries"`
	WithheldScope        string   `json:"withheld_scope"`
}

// Group is a GNU Social Group, gets returned by GS API.
type Group struct {
	Id          int64  `json:"id"`
	Url         string `json:"url"`
	Nickname    string `json:"nickname"`
	Fullname    string `json:"fullname"`
	Member      bool   `json:"member"`
	Membercount int64  `json:"member_count"`
	Description string `json:"description"`
}

// Badrequest is an error. If the API doesn't respond how we like, it replies using this struct (in json)
type Badrequest struct {
	Error   string `json:"error"`
	Request string `json:"request"`
}
