/*
Package quitter is a Go library to interact with GNU Social instances.

		The MIT License (MIT)

		Copyright (c) 2016 aerth

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
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
)

// GetPublic shows 20 new messages.
func (a Social) GetPublic() ([]Quip, error) {
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
func (a Social) GetMentions() ([]Quip, error) {
	if a.Username == "" || a.Password == "" {
		return nil, errors.New("No user/password")
	}
	path := a.Scheme + a.Node + "/api/statuses/mentions.json"
	body, err := a.FireGET(path)
	if err != nil {
		return nil, err
	}
	var quips []Quip
	_ = json.Unmarshal(body, &quips)

	return quips, err
}

// GetHome shows 20 from home timeline.
func (a Social) GetHome() ([]Quip, error) {
	if a.Username == "" || a.Password == "" {
		return nil, errors.New("No user/password")
	}

	path := "/api/statuses/home_timeline.json"

	body, err := a.FireGET(path)
	if err != nil {
		return nil, err
	}
	var quips []Quip
	_ = json.Unmarshal(body, &quips)

	return quips, err

}

// DoSearch returns results for query searchstr. Does send auth info.
func (a Social) DoSearch(searchstr string) ([]Quip, error) {
	if a.Username == "" || a.Password == "" {
		return nil, errors.New("No user/password")
	}
	if searchstr == "" {
		return nil, errors.New("No query")
	}

	v := url.Values{}
	v.Set("q", searchstr)
	searchq := url.Values.Encode(v)

	path := "/api/search.json?" + searchq
	body, err := a.FireGET(path)
	if err != nil {
		return nil, err
	}
	var quips []Quip
	_ = json.Unmarshal(body, &quips)

	return quips, err

}

// DoPublicSearch returns results for query searchstr. Does not send auth info.
func (a Social) DoPublicSearch(searchstr string) ([]Quip, error) {
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

// DoFollow sends a request to follow a user
func (a Social) DoFollow(followstr string) (user User, err error) {
	if a.Username == "" || a.Password == "" {
		return user, errors.New("No user/password")
	}

	if followstr == "" {
		return user, errors.New("No query.")
	}

	v := url.Values{}

	v.Set("id", followstr)
	followstr = url.Values.Encode(v)

	path := "/api/friendships/create.json?" + followstr
	body, err := a.FirePOST(path, v)
	if err != nil {
		return user, err
	}

	// Return one user
	_ = json.Unmarshal(body, &user)

	return user, err

}

// DoUnfollow sends a request to unfollow a user
func (a Social) DoUnfollow(followstr string) (user User, err error) {
	if a.Username == "" || a.Password == "" {
		return user, errors.New("No user/password")
	}
	if followstr == "" {
		return user, errors.New("No query")
	}
	v := url.Values{}
	v.Set("id", followstr)
	followstr = url.Values.Encode(v)
	path := "/api/friendships/destroy.json?" + followstr
	body, err := a.FirePOST(path, v)
	if err != nil {
		return user, err
	}
	_ = json.Unmarshal(body, &user)

	return user, err

}

// GetUserTimeline returns a userlookup's timeline
func (a Social) GetUserTimeline(userlookup string) ([]Quip, error) {

	path := "/api/statuses/user_timeline.json?screen_name=" + userlookup
	body, err := a.FireGET(path)
	if err != nil {
		return nil, err
	}

	var quips []Quip
	_ = json.Unmarshal(body, &quips)
	return quips, err
}

// PostNew publishes content on a.Node, returns the new quip or an error.
func (a Social) PostNew(content string) (q Quip, err error) {
	if a.Username == "" || a.Password == "" {
		return q, errors.New("No user/password")
	}

	if content == "" {
		return q, errors.New("No query")
	}
	content = url.QueryEscape(content)

	v := url.Values{}
	v.Set("status", content)
	content = url.Values.Encode(v)
	path := "/api/statuses/update.json?" + content
	body, err := a.FirePOST(path, v)
	if err != nil {
		return q, err
	}
	_ = json.Unmarshal(body, &q)
	return q, err

}

//  ListAllGroups lists each group in a.Node
//  Some nodes don't return anything.
func (a Social) ListAllGroups() ([]Group, error) {
	if a.Username == "" || a.Password == "" {
		return nil, errors.New("No user/password")
	}

	path := a.Scheme + a.Node + "/api/statusnet/groups/list_all.json"
	body, err := a.FireGET(path)
	if err != nil {
		return nil, err
	}

	var groups []Group
	err = json.Unmarshal(body, &groups)

	return groups, err
}

// ListMyGroups lists each group a a.Username is a member of
func (a Social) ListMyGroups() ([]Group, error) {
	if a.Username == "" || a.Password == "" {
		return nil, errors.New("No user/password")
	}

	path := "/api/statusnet/groups/list.json"
	body, err := a.FireGET(path)
	if err != nil {
		return nil, err
	}

	var groups []Group
	_ = json.Unmarshal(body, &groups)

	return groups, err

}

// JoinGroup sends a request to join group grp.
func (a Social) JoinGroup(grp string) (g Group, err error) {
	if a.Username == "" || a.Password == "" {
		return g, errors.New("No user/password")
	}

	if grp == "" {
		return g, errors.New("Blank group detected. Not going furthur.")
	}
	v := url.Values{}

	v.Set("group_name", grp)
	v.Set("group_id", grp)
	v.Set("id", grp)
	v.Set("nickname", grp)
	grp = url.Values.Encode(v)

	path := "/api/statusnet/groups/join.json?" + grp
	body, err := a.FirePOST(path, v)
	if err != nil {
		return g, err
	}

	_ = json.Unmarshal(body, &g)

	return g, err
}

// PartGroup sends a request to part group grp.
func (a Social) PartGroup(grp string) (g Group, err error) {
	if a.Username == "" || a.Password == "" {
		return g, errors.New("No user/password")
	}

	if grp == "" {
		return g, errors.New("Blank group detected. Not going furthur.")
	}

	v := url.Values{}
	v.Set("group_name", grp)
	v.Set("group_id", grp)
	v.Set("id", grp)
	v.Set("nickname", grp)
	grp = url.Values.Encode(v)
	path := "/api/statusnet/groups/leave.json?" + grp
	body, err := a.FirePOST(path, v)
	if err != nil {
		return g, err
	}

	_ = json.Unmarshal(body, &g)

	return g, err
}
