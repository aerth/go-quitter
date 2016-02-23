package main

import (
//  "os"
  "log"
  "fmt"
  "io/ioutil"
  "net/http"
  "encoding/json"

  //"gopkg.in/resty.v0"
//    "github.com/aerth/anaconda"

)


func main() {

  log.Println("go-quitter v0.0.1")
/*
    if os.Getenv("GNUSOCIALKEY") == "" {
    		fmt.Println("Set environmental variable GNUSOCIALKEY before running go-quitter.")
    		fmt.Println("GNUSOCIALKEY before running go-quitter.")
    		os.Exit(1)
    	}
    if os.Getenv("GNUSOCIALSECRET") == "" {
    		fmt.Println("Set environmental variable GNUSOCIALSECRET before running go-quitter.")
    		os.Exit(1)
    }
*/

        type User struct {
            Name      string  `json:"name"`
        }
    type Tweet struct {

    	Id                   int64                  `json:"id"`
    	IdStr                string                 `json:"id_str"`
    	InReplyToScreenName  string                 `json:"in_reply_to_screen_name"`
    	InReplyToStatusID    int64                  `json:"in_reply_to_status_id"`
    	InReplyToStatusIdStr string                 `json:"in_reply_to_status_id_str"`
    	InReplyToUserID      int64                  `json:"in_reply_to_user_id"`
    	InReplyToUserIdStr   string                 `json:"in_reply_to_user_id_str"`
    	Lang                 string                 `json:"lang"`
    	Place                string                  `json:"place"`
    	PossiblySensitive    bool                   `json:"possibly_sensitive"`
    	RetweetCount         int                    `json:"retweet_count"`
    	Retweeted            bool                   `json:"retweeted"`
    	RetweetedStatus      *Tweet                 `json:"retweeted_status"`
    	Source               string                 `json:"source"`

    	Text                 string                 `json:"text"`
    	Truncated            bool                   `json:"truncated"`
    	User                 User                   `json:"user"`
    	WithheldCopyright    bool                   `json:"withheld_copyright"`
    	WithheldInCountries  []string               `json:"withheld_in_countries"`
    	WithheldScope        string                 `json:"withheld_scope"`

    	//Geo is deprecated
    	//Geo                  interface{} `json:"geo"`
    }

        v := Tweet{}

    /*
    if os.Getenv("GNUSOCIALACCESSTOKEN") == "" {
    		fmt.Println("Set environmental variable GNUSOCIALACCESSTOKEN before running go-quitter.")
    		os.Exit(1)
    }
    if os.Getenv("GNUSOCIALTOKENSECRET") == "" {
    		fmt.Println("Set environmental variable GNUSOCIALTOKENSECRET before running go-quitter.")
    		os.Exit(1)
    }
    */

    res, err := http.Get("https://gs.sdf.org/api/statuses/show/1.json")
    body, err := ioutil.ReadAll(res.Body)
    defer res.Body.Close()

    err = json.Unmarshal(body, &v)
    if err != nil {
  		log.Fatal(err)
  	}

  	fmt.Printf("%s says %#v\n", v.User.Name, v.Text)


}
