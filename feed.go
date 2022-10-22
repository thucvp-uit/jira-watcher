package main

import "encoding/xml"

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Id      string   `xml:"id"`
	Title   string   `xml:"title"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	XMLName   xml.Name `xml:"entry"`
	Id        string   `xml:"id"`
	Title     string   `xml:"title"`
	Summary   string   `xml:"summary"`
	Content   string   `xml:"content"`
	Published string   `xml:"published"`
	Updated   string   `xml:"updated"`
	Object    Object   `xml:"object"`
	Target    Target   `xml:"target"`
	Author    Author   `xml:"author"`
}

type Target struct {
	XMLName xml.Name `xml:"target"`
	Id      string   `xml:"id"`
	Title   string   `xml:"title"`
}

type Object struct {
	XMLName xml.Name `xml:"object"`
	Id      string   `xml:"id"`
	Title   string   `xml:"title"`
}

type Author struct {
	XMLName  xml.Name `xml:"author"`
	Name     string   `xml:"name"`
	UserName string   `xml:"username"`
}
