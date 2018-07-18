package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

// Tree struct to contain information on a tree species
type Tree struct {
	Genus   string `json:"genus"`
	Species string `json:"species"`
}

var trees = readFromFile("global_tree_search_trees_1_2.csv")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFromFile(path string) []Tree {
	file, err := os.Open(path)
	check(err)
	reader := csv.NewReader(bufio.NewReader(file))
	var coll []Tree
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		}
		check(error)
		str := strings.Fields(line[0])
		coll = append(coll, Tree{
			Genus:   str[0],
			Species: str[1],
		})
	}
	file.Close()
	return coll
}

func getShortName(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	w.Header().Set("Content-Type", "application/json")
	t := trees[rand.Intn(len(trees))]
	str, err := json.Marshal(t)
	check(err)
	fmt.Fprintf(w, string(str))
}

func main() {
	http.HandleFunc("/", getShortName)     // set router
	err := http.ListenAndServe(":80", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
