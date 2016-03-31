package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// ResponseType structure is what we expect to get for our URL call.
// It contains a request URL, the data chunk and possible error from remote
type ResponseType struct {
	Url   string
	Data  []byte
	Error error
}

// helper function to wrap http request. It returns bare response which caller must close
func httpRequest(method, rurl, args string, verbose bool) ResponseType {
	var response ResponseType
	response.Url = rurl
	response.Data = []byte{}

	var req *http.Request
	if method == "POST" {
		var jsonStr = []byte(args)
		req, _ = http.NewRequest("POST", rurl, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest("GET", rurl, nil)
		req.Header.Set(
			"Accept", "application/vnd.github.v3.text-match+json")
	}

	if verbose {
		dump1, err1 := httputil.DumpRequestOut(req, true)
		log.Println("### HTTP request", string(dump1), err1)
	}
	resp, err := http.DefaultClient.Do(req)
	if verbose {
		dump2, err2 := httputil.DumpResponse(resp, true)
		log.Println("### HTTP response", string(dump2), err2)
	}
	if resp.StatusCode != http.StatusOK {
		response.Error = fmt.Errorf("github request failed: %s", resp.Status)
		return response
	}
	if err != nil {
		response.Error = err
		return response
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.Error = err
		return response
	}
	response.Data = body
	return response
}

// helper function to call given url
func apiCall(rurl string, verbose bool) (*Results, error) {
	resp := httpRequest("GET", rurl, "", verbose)
	if resp.Error != nil {
		return nil, resp.Error
	}
	var result Results
	var items []*Item
	if err := json.Unmarshal(resp.Data, &items); err != nil {
		return nil, err
	}
	result.TotalCount = len(items)
	result.Items = items
	return &result, nil
}

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms string, verbose bool) (*Results, error) {
	q := url.QueryEscape(terms)
	rurl := URL + "/search/issues" + "?q=" + q
	return apiCall(rurl, verbose)
}

// GET /repos/:owner/repos
func Repos(user string, verbose bool) (*Results, error) {
	rurl := URL + "/users/" + user + "/repos"
	return apiCall(rurl, verbose)
}

// GET /repos/:owner/:repo/issues
func Issues(repo string, verbose bool) (*Results, error) {
	rurl := URL + "/repos/" + repo + "/issues"
	return apiCall(rurl, verbose)
}

// GET /repos/:owner/:repo/issues
func IssueDetails(arg string, verbose bool) (*Results, error) {
	args := strings.Split(arg, " ")
	rurl := URL + "/repos/" + args[0] + "/issues/" + args[1]
	resp := httpRequest("GET", rurl, "", verbose)
	if resp.Error != nil {
		return nil, resp.Error
	}
	var result Results
	var item Item
	if err := json.Unmarshal(resp.Data, &item); err != nil {
		return nil, err
	}
	var items []*Item
	items = append(items, &item)
	// get comments on issue
	rurl = URL + "/repos/" + args[0] + "/issues/" + args[1] + "/comments"
	resp = httpRequest("GET", rurl, "", verbose)
	if resp.Error != nil {
		return nil, resp.Error
	}
	var comments []*Item
	if err := json.Unmarshal(resp.Data, &comments); err != nil {
		return nil, err
	}
	for _, c := range comments {
		items = append(items, c)
	}

	result.TotalCount = len(items)
	result.Items = items
	return &result, nil
}
