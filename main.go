// Package ghub provides a Go APIs for github issue tracker.
// See https://developer.github.com/v3/search/#search-issues.
//
// This work was inspired by
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/
//
// Valentin Kuznetsov (<vkuznet AT gmail dot com>) contributions:
// - provides access to search/users/repos github APIs
// - colorize the output
// - extend templates
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const Usage = "Usage: ghub <repos|search|issues|issue> <user|repo|issue>\n"

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(Usage)
		return
	}
	verbose := false
	args := strings.Join(os.Args[2:], " ")
	if strings.Contains(args, "-verbose") {
		args = strings.Replace(args, "-verbose", "", 1)
		verbose = true
	}
	args = strings.TrimSpace(args)
	request := os.Args[1]
	var result *Results
	var err error
	switch request {
	case "search":
		result, err = SearchIssues(args, verbose)
	case "issues":
		result, err = Issues(args, verbose)
	case "issue":
		result, err = IssueDetails(args, verbose)
	case "repos":
		result, err = Repos(args, verbose)
	default:
		fmt.Printf(Usage)
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	PrintResults(request, *result)
}
