package main

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Location    string `json:"location"`
	DateTime    string `json:"dateTime"`
	Source      string `json:"source"`
}

type ErrMsg struct{ Err error }
