package quitter

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// fireGET sends a get request and returns byte, error.
func (a Account) fireGET(path string) ([]byte, error) {

	if path == "" {
		return nil, errors.New("No path")
	}
	apipath := a.Scheme + a.Node + path
	req, err := http.NewRequest("GET", apipath, nil)
	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("User-Agent", UserAgent)
	if err != nil {
		return nil, err
	}

	var resp *http.Response

	resp, err = apigun.Do(req)


	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) == "" {
		return nil, errors.New("node response: " + resp.Status)
	}

	/*
		Here, I tried a couple ways of telling if the response was an error.
		This one ended up working the most reliably. TODO:

	*/
	var apiresponse Badrequest
	_ = json.Unmarshal(body, &apiresponse)
	if apiresponse.Error != "" {
		return nil, errors.New(apiresponse.Error)
	}

	return body, nil
}

// firePOST uses the account details to send a POST request to the node.
// HTTP_REFERER is added to keep nodes happy.
func (a Account) firePOST(path string, v url.Values) ([]byte, error) {
	if path == "" {
		return nil, errors.New("No path")
	}
	if v.Encode() == "" && !strings.Contains(path, "update") { // update can use a blank post request..
		return nil, errors.New("No values to post")
	}
	apipath := a.Scheme + a.Node + path
	b := bytes.NewBufferString(v.Encode())
	req, err := http.NewRequest("POST", apipath, b)
	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("HTTP_REFERER", a.Scheme+a.Node+"/")
	// req.Header.Add("Content-Type", "[application/json; charset=utf-8") // is this a typo ?
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", UserAgent)
	var resp *http.Response
	resp, err = apigun.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) == "" {
		return nil, errors.New("node response: " + resp.Status)
	}

	var apiresponse Badrequest
	_ = json.Unmarshal(body, &apiresponse)
	if apiresponse.Error != "" {
		return nil, errors.New(apiresponse.Error)
	}
	return body, nil
}
