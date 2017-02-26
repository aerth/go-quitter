/*
Package quitter is a Go library to interact with GNU Social instances.

		The MIT License (MIT)

		Copyright (c) 2016-2017 aerth <aerth@riseup.net>

		Permission is hereby granted, free of charge, to any person obtaining a copy
		of this software and associated documentation files (the "Software"), to deal
		in the Software without restriction, including without limitation the rights
		to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
		copies of the Software, and to permit persons to whom the Software is
		furnished to do so, subject to the following conditions:

		The above copyright notice and this permission notice shall be included in all
		copies or substantial portions of the Software.
*/
package quitter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// GetPublic shows 20 new messages.
func (a Account) GetPublic() ([]Quip, error) {
	resp, err := apigun.Get(a.Scheme + a.Node + "/api/statuses/public_timeline.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var quips []Quip

	var apiresponse Badrequest
	_ = json.Unmarshal(body, &apiresponse)
	if apiresponse.Error != "" {
		return nil, errors.New(apiresponse.Error)
	}

	_ = json.Unmarshal(body, &quips)

	return quips, err
}

// GetMentions shows 20 newest mentions of your username.
func (a Account) GetMentions() ([]Quip, error) {
	if a.Username == "" || a.Password == "" {
		return nil, errors.New("Invalid Credentials")
	}
	path := "/api/statuses/mentions.json"
	body, err := a.fireGET(path)
	if err != nil {
		return nil, err
	}
	var quips []Quip
	_ = json.Unmarshal(body, &quips)

	return quips, err
}

// GetHome shows 20 from home timeline.
func (a Account) GetHome() ([]Quip, error) {
	if a.Username == "" || a.Password == "" {
		return nil, errors.New("Invalid Credentials")
	}

	path := "/api/statuses/home_timeline.json"

	body, err := a.fireGET(path)
	if err != nil {
		return nil, err
	}
	var quips []Quip
	_ = json.Unmarshal(body, &quips)

	return quips, err

}

// Search returns results for query searchstr. Does send auth info.
func (a Account) Search(searchstr string) ([]Quip, error) {
	if a.Username == "" || a.Password == "" {
		return nil, errors.New("Invalid Credentials")
	}
	if searchstr == "" {
		return nil, errors.New("No query")
	}

	v := url.Values{}
	v.Set("q", searchstr)
	searchq := url.Values.Encode(v)

	path := "/api/search.json?" + searchq
	body, err := a.fireGET(path)
	if err != nil {
		return nil, err
	}
	var quips []Quip
	_ = json.Unmarshal(body, &quips)

	return quips, err

}

// PublicSearch returns results for query searchstr. Does not send auth info.
func (a Account) PublicSearch(searchstr string) ([]Quip, error) {
	if searchstr == "" {
		return nil, errors.New("No query")
	}

	v := url.Values{}
	v.Set("q", searchstr)
	searchq := url.Values.Encode(v)
	resp, err := apigun.Get(a.Scheme + a.Node + "/api/search.json?" + searchq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var quips []Quip
	var apiresponse Badrequest
	_ = json.Unmarshal(body, &apiresponse)
	if apiresponse.Error != "" {
		return nil, errors.New(apiresponse.Error)
	}

	_ = json.Unmarshal(body, &quips)

	return quips, err

}

// Follow sends a request to follow a user
func (a Account) Follow(followstr string) (user User, err error) {
	if a.Username == "" || a.Password == "" {
		return user, errors.New("Invalid Credentials")
	}

	if followstr == "" {
		return user, errors.New("no query")
	}

	v := url.Values{}

	v.Set("id", followstr)
	followstr = url.Values.Encode(v)

	path := "/api/friendships/create.json?" + followstr
	body, err := a.firePOST(path, v)
	if err != nil {
		return user, err
	}

	// Return one user
	_ = json.Unmarshal(body, &user)

	return user, err

}

// UnFollow sends a request to unfollow a user
func (a Account) UnFollow(followstr string) (user User, err error) {
	if a.Username == "" || a.Password == "" {
		return user, errors.New("Invalid Credentials")
	}
	if followstr == "" {
		return user, errors.New("No query")
	}
	v := url.Values{}
	v.Set("id", followstr)
	followstr = url.Values.Encode(v)
	path := "/api/friendships/destroy.json?" + followstr
	body, err := a.firePOST(path, v)
	if err != nil {
		return user, err
	}
	_ = json.Unmarshal(body, &user)

	return user, err

}

// GetUserTimeline returns a userlookup's timeline
func (a Account) GetUserTimeline(userlookup string) ([]Quip, error) {

	path := "/api/statuses/user_timeline.json?screen_name=" + userlookup
	body, err := a.fireGET(path)
	if err != nil {
		return nil, err
	}

	var quips []Quip
	_ = json.Unmarshal(body, &quips)
	return quips, err
}

// PostNew publishes content on a.Node, returns the new quip or an error.
func (a Account) PostNew(content string) (q Quip, err error) {
	if a.Username == "" || a.Password == "" {
		return q, errors.New("Invalid Credentials")
	}

	if content == "" {
		return q, errors.New("No query")
	}
	//content = url.QueryEscape(content)

	v := url.Values{}
	v.Set("status", content)
	content = url.Values.Encode(v)
	path := "/api/statuses/update.json?" + content
	body, err := a.firePOST(path, nil)
	if err != nil {
		return q, err
	}
	_ = json.Unmarshal(body, &q)
	return q, err

}

// ListAllGroups lists each group in a.Node
// Some nodes don't return anything.
func (a Account) ListAllGroups() ([]Group, error) {
	if a.Username == "" || a.Password == "" {
		return nil, errors.New("Invalid Credentials")
	}

	path := "/api/statusnet/groups/list_all.json"
	body, err := a.fireGET(path)
	if err != nil {
		return nil, err
	}

	var groups []Group
	err = json.Unmarshal(body, &groups)

	return groups, err
}

// ListMyGroups lists each group a a.Username is a member of
func (a Account) ListMyGroups() ([]Group, error) {
	if a.Username == "" || a.Password == "" {
		return nil, errors.New("Invalid Credentials")
	}

	path := "/api/statusnet/groups/list.json"
	body, err := a.fireGET(path)
	if err != nil {
		return nil, err
	}

	var groups []Group
	_ = json.Unmarshal(body, &groups)

	return groups, err

}

// JoinGroup sends a request to join group grp.
func (a Account) JoinGroup(grp string) (g Group, err error) {
	if a.Username == "" || a.Password == "" {
		return g, errors.New("Invalid Credentials")
	}

	if grp == "" {
		return g, errors.New("no group")
	}
	v := url.Values{}

	v.Set("group_name", grp)
	v.Set("group_id", grp)
	v.Set("id", grp)
	v.Set("nickname", grp)
	grp = url.Values.Encode(v)

	path := "/api/statusnet/groups/join.json?" + grp
	body, err := a.firePOST(path, v)
	if err != nil {
		return g, err
	}

	_ = json.Unmarshal(body, &g)

	return g, err
}

// PartGroup sends a request to part group grp.
func (a Account) PartGroup(grp string) (g Group, err error) {
	if a.Username == "" || a.Password == "" {
		return g, errors.New("Invalid Credentials")
	}

	if grp == "" {
		return g, errors.New("no group")
	}

	v := url.Values{}
	v.Set("group_name", grp)
	v.Set("group_id", grp)
	v.Set("id", grp)
	v.Set("nickname", grp)
	grp = url.Values.Encode(v)
	path := "/api/statusnet/groups/leave.json?" + grp
	body, err := a.firePOST(path, v)
	if err != nil {
		return g, err
	}

	_ = json.Unmarshal(body, &g)

	return g, err
}

// Upload an image and optional caption
func (a Account) Upload(fpath string, content ...string) (quip Quip, err error) {

	// Optional Caption
	var caption string
	if content != nil {
		caption = strings.Join(content, "\n")
	} else {
		caption = " "
	}

	// Multipart Form
	var b = new(bytes.Buffer)
	w := multipart.NewWriter(b)
	f, err := os.Open(fpath)
	if err != nil {
		return Quip{}, err
	}
	defer f.Close()
	fw, err := w.CreateFormFile("media", filepath.Base(fpath))
	if err != nil {
		return Quip{}, err
	}
	fb, err := ioutil.ReadAll(f)
	if err != nil {
		return Quip{}, err
	}
	n, err := fw.Write(fb)
	if err != nil {
		return Quip{}, err
	}
	if fw, err = w.CreateFormField("status"); err != nil {
		return Quip{}, err
	}
	if _, err = fw.Write([]byte(caption)); err != nil {
		return
	}
	w.Close()

	// Request
	uploadURL := a.Scheme + a.Node + "/api/statuses/update.json"
	req, err := http.NewRequest("POST", uploadURL, b)

	if err != nil {
		return Quip{}, err
	}
	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("User-Agent", goquitter)
	req.Header.Set("Content-Type", w.FormDataContentType())
	fmt.Printf("%v bytes uploading\n", n)
	res, err := apigun.Do(req)
	if err != nil {
		return Quip{}, err
	}


	// Response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return Quip{}, err
	}

	var bb []byte
	bb, err = ioutil.ReadAll(res.Body)
	bstr := string(bb)
	if strings.Contains("Page not found", bstr) {
		return Quip{Text: "404 Page not found"}, err
	}
	var tweet Quip
	err = json.Unmarshal(bb, &tweet)

	return tweet, err
}
