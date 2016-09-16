// The MIT License (MIT)
//
// Copyright (c) 2016 aerth
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

package quitter

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
)

// GetPublic shows 20 new messages. Defaults to a 2 second delay, but can be called with GetPublic(fast) for a quick dump. This and DoSearch() and GetUserTimeline() are some of the only functions that don't require auth.Username + auth.Password
func (a Auth) GetPublic(fast bool) ([]Quip, error) {
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

// GetMentions shows 20 newest mentions of your username. Defaults to a 2 second delay, but can be called with GetPublic(fast) for a quick dump.
func (a Auth) GetMentions(fast bool) ([]Quip, error) {
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

// GetHome shows 20 from home timeline. Defaults to a 2 second delay, but can be called with GetHome(fast) for a quick dump.
func (a Auth) GetHome(fast bool) ([]Quip, error) {
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

// command: go-quitter search
func (a Auth) DoSearch(searchstr string, fast bool) ([]Quip, error) {
	if a.Username == "" || a.Password == "" {
		return nil, errors.New("No user/password")
	}
	if searchstr == "" {
		return nil, errors.New("Blank search detected. Not searching.")
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

// command: go-quitter psearch
func (a Auth) DoPublicSearch(searchstr string, fast bool) ([]Quip, error) {
	if searchstr == "" {
		return nil, errors.New("Blank search detected. Not searching.")
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

// command: go-quitter follow
func (a Auth) DoFollow(followstr string) (user User, err error) {
	if a.Username == "" || a.Password == "" {
		return user, errors.New("No user/password")
	}

	if followstr == "" {
		return user, errors.New("Blank search detected. Not going furthur.")
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

// go-quitter command: go-quitter unfollow
func (a Auth) DoUnfollow(followstr string) (user User, err error) {
	if a.Username == "" || a.Password == "" {
		return user, errors.New("No user/password")
	}
	if followstr == "" {
		return user, errors.New("Blank search detected. Not going furthur.")
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

// go-quitter command: go-quitter user
func (a Auth) GetUserTimeline(userlookup string, fast bool) ([]Quip, error) {

	path := "/api/statuses/user_timeline.json?screen_name=" + userlookup
	body, err := a.FireGET(path)
	if err != nil {
		return nil, err
	}

	var quips []Quip
	_ = json.Unmarshal(body, &quips)
	return quips, err
}

// go-quitter command: go-quitter post
func (a Auth) PostNew(content string) (q Quip, err error) {
	if a.Username == "" || a.Password == "" {
		return q, errors.New("No user/password")
	}

	if content == "" {
		return q, errors.New("Blank status detected. Not posting.")
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

// command: go-quitter groups
func (a Auth) ListAllGroups(speed bool) ([]Group, error) {
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

// command: go-quitter mygroups
func (a Auth) ListMyGroups(speed bool) ([]Group, error) {
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

// command: go-quitter join ____
func (a Auth) JoinGroup(groupstr string) (g Group, err error) {
	if a.Username == "" || a.Password == "" {
		return g, errors.New("No user/password")
	}

	if groupstr == "" {
		return g, errors.New("Blank group detected. Not going furthur.")
	}
	v := url.Values{}

	v.Set("group_name", groupstr)
	v.Set("group_id", groupstr)
	v.Set("id", groupstr)
	v.Set("nickname", groupstr)
	groupstr = url.Values.Encode(v)

	path := "/api/statusnet/groups/join.json?" + groupstr
	body, err := a.FirePOST(path, v)
	if err != nil {
		return g, err
	}

	_ = json.Unmarshal(body, &g)

	return g, err
}

// command: go-quitter part ____
func (a Auth) PartGroup(groupstr string) (g Group, err error) {
	if a.Username == "" || a.Password == "" {
		return g, errors.New("No user/password")
	}

	if groupstr == "" {
		return g, errors.New("Blank group detected. Not going furthur.")
	}

	v := url.Values{}
	v.Set("group_name", groupstr)
	v.Set("group_id", groupstr)
	v.Set("id", groupstr)
	v.Set("nickname", groupstr)
	groupstr = url.Values.Encode(v)
	path := "/api/statusnet/groups/leave.json?" + groupstr
	body, err := a.FirePOST(path, v)
	if err != nil {
		return g, err
	}

	_ = json.Unmarshal(body, &g)

	return g, err
}
