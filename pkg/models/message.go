package models

type Message struct {
	Type       string      `json:"type"`
	Data       interface{} `json:"data"`
	TimeStr    string      `json:"time"`
	ServerName string      `json:"serverName"`
}
