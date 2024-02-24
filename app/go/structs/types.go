package structs

import "html/template"

type Content struct {
	Title string
	Body  template.HTML
}

type Contact struct {
	Id    string
	First string
	Last  string
	Phone string
	Email string
}
