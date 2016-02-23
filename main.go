package main

import (
  "os"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	//"gopkg.in/resty.v0"
	//    "github.com/aerth/anaconda"
)

func main() {

	log.Println("go-quitter v0.0.1")
	log.Println("Copyright 2016 aerth@sdf.org")
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
		Name string `json:"name"`
	}
	type Tweet struct {
		Id int64 `json:"id"`
		/*		IdStr                string `json:"id_str"`
				InReplyToScreenName  string `json:"in_reply_to_screen_name"`
				InReplyToStatusID    int64  `json:"in_reply_to_status_id"`
				InReplyToStatusIdStr string `json:"in_reply_to_status_id_str"`
				InReplyToUserID      int64  `json:"in_reply_to_user_id"`
				InReplyToUserIdStr   string `json:"in_reply_to_user_id_str"`
				Lang                 string `json:"lang"`
				Place                string `json:"place"`
				PossiblySensitive    bool   `json:"possibly_sensitive"`
				RetweetCount         int    `json:"retweet_count"`
				Retweeted            bool   `json:"retweeted"`
				RetweetedStatus      *Tweet `json:"retweeted_status"`
				Source               string `json:"source"`
		*/
		Text string `json:"text"`
		//		Truncated           bool     `json:"truncated"`
		User User `json:"user"`
		/*		WithheldCopyright   bool     `json:"withheld_copyright"`
				WithheldInCountries []string `json:"withheld_in_countries"`
				WithheldScope       string   `json:"withheld_scope"`*/
		//Geo is deprecated
		//Geo                  interface{} `json:"geo"`
	}

	/*
	   type Tweets struct {

	   	Status []*Tweet `json:"status"`


	   }
	*/
	//	v := []Tweets{}
	//	s := Tweet{}

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
	var gnusocialnode string
	if os.Getenv("GNUSOCIALNODE") == "" {
		gnusocialnode = "gs.sdf.org"
	}else{
		gnusocialnode = os.Getenv("GNUSOCIALNODE")
	}

		res, err := http.Get("https://"+ gnusocialnode +"/api/statuses/public_timeline.json")
		if err != nil {
			log.Fatalln(err)
		}


	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var tweets []Tweet
	err = json.Unmarshal(body, &tweets)

	if err != nil {
		fmt.Println("error:", err)
	}

	for i := range tweets {
		fmt.Printf("[" + tweets[i].User.Name + "] " + tweets[i].Text + "\n\n")
		time.Sleep(1000 * time.Millisecond)
	}

}
