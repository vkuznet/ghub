package main

import "time"

const URL = "https://api.github.com"

type Results struct {
	TotalCount int `json:"total_count"`
	Items      []*Item
}

type Item struct {
	Number      int
	HTMLURL     string `json:"html_url"`
	Title       string
	State       string
	FullName    string `json:"full_name"`
	Description string
	Language    string
	Fork        bool
	Homepage    string
	User        *User
	CreatedAt   time.Time `json:"created_at"`
	Body        string    // in Markdown format
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}
