package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"regexp"
)

type Link struct {
	Link  string
	Title string
}

var idRegex = regexp.MustCompile("^[abcdefABCDEF0123456789]{16}$")
var tmpl = template.Must(template.ParseFiles("page.html"))

const domain = "http://localhost:8080"

func pageHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if idRegex.MatchString(id) {
		bytes, err := hex.DecodeString(id)
		if err != nil {
			panic("failed to decode")
		}

		num := binary.BigEndian.Uint64(bytes)

		links := make([]Link, 64)
		shuffle := rand.Perm(64)

		for i := uint64(0); i < 64; i += 1 {
			newNum := num ^ (uint64(0x01) << i)
			newId := fmt.Sprintf("%08x", newNum)
			links[shuffle[i]] = Link{
				Title: newId,
				Link:  domain + "/?id=" + newId,
			}
		}

		err = tmpl.Execute(w, links)
		if err != nil {
			panic("Failed to execute template")
		}
	} else {
		_, err := fmt.Fprintf(w, "Invalid id")
		if err != nil {
			panic("Failed to write body")
		}
	}
}

func main() {
	fmt.Println("vim-go")
	http.HandleFunc("/", pageHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
