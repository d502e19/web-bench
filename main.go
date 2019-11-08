package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

type Link struct {
	Link  string
	Title string
}

var idRegex = regexp.MustCompile("^[abcdefABCDEF0123456789]{16}$")
var tmpl = template.Must(template.ParseFiles("page.html"))

func pageHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if idRegex.MatchString(id) {
		bytes, err := hex.DecodeString(id)
		if err != nil {
			panic("failed to decode")
		}

		num := binary.BigEndian.Uint64(bytes)

		links := make([]Link, len(constants))

		for i := 0; i < 64; i += 1 {
			newNum := num ^ (0x01 << i)
			newId := fmt.Sprintf("%08x", newNum)
			links[i] = Link{
				Title: newId,
				Link:  "http://localhost:8080/?id=" + newId,
			}
		}

		tmpl.Execute(w, links)
	} else {
		fmt.Fprintf(w, "Invalid id")
	}
}

func main() {
	fmt.Println("vim-go")
	http.HandleFunc("/", pageHandler)
	http.ListenAndServe(":8080", nil)
}
