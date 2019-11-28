package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Link struct {
	Link  string
	Title string
}

var tmpl = template.Must(template.ParseFiles("page.html"))
var linkMap = make(map[int][]Link)

const domain = "http://localhost:8080"

func pageHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		// id is not an int. Create error page
		_, err := fmt.Fprintf(w, "Invalid id")
		if err != nil {
			panic("Failed to write body")
		}
	} else {
		// Create page
		links := linkMap[id]
		err = tmpl.Execute(w, links)
		if err != nil {
			panic("Failed to execute template")
		}
	}
}

func loadGraph(filename string) error {
	// Open the file
	csvfile, err := os.Open(filename)
	if err != nil {
		return err
	}
	reader := csv.NewReader(bufio.NewReader(csvfile))

	edges := 0
	for {
		// Read the values of one line the graph
		values, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Parse the values
		from, to := values[0], values[1]
		fromId, err := strconv.Atoi(from)
		if err != nil {
			return err
		}
		_, err = strconv.Atoi(to)
		if err != nil {
			return err
		}

		// Construct link
		link := Link{
			Title: to,
			Link:  domain + "/?id=" + to,
		}
		pageLinks := linkMap[fromId]
		pageLinks = append(pageLinks, link)
		linkMap[fromId] = pageLinks
		edges++
	}

	fmt.Printf("%s was loaded succesfully (%d edges)\n", filename, edges)
	return nil
}

func main() {
	// A graph generated using http://pywebgraph.sourceforge.net/
	err := loadGraph("graph.csv")
	if err != nil {
		panic(err)
	}

	// Starting server simulating the graph
	fmt.Println("Starting server..")
	http.HandleFunc("/", pageHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
