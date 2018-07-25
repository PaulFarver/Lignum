package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

type Tree struct {
	TaxonIDS    []string `json:"taxonIds"`
	Genus       string   `json:"genus"`
	Species     string   `json:"species"`
	CommonNames []Name   `json:"commonNames"`
}

type Name struct {
	Language string `json:"language"`
	Name     string `json:"name"`
}

var trees []Tree

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFromFile(path string) []Tree {
	plan, _ := ioutil.ReadFile(path)
	var data []Tree
	err := json.Unmarshal(plan, &data)
	check(err)
	return data
}

func getRandomTree(lang []string) Tree {
	if len(lang) == 0 {
		lang = append(lang, "en")
	}
	t := trees[rand.Intn(len(trees))]
	return t
}

func getShortName(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("path", r.URL.Path)
	// lang := r.URL.Query().Get("language")
	// allowEmpty := r.URL.Query().Get("allowEmptyLang")
	lang, l := r.URL.Query()["language"]
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	w.Header().Set("Content-Type", "application/json")
	t := getRandomTree(lang)
	// fmt.Println(getKey(t))
	// fmt.Println(getNames(getKey(t)))
	// t.CommonName = getNames(getKey(t))
	str, err := json.Marshal(t)
	check(err)
	fmt.Fprintf(w, string(str))
}

func main() {
	datafile := flag.String("data", "data.json", "The file with data in it")

	flag.Parse()

	fmt.Println("data: ", *datafile)
	trees = readFromFile(*datafile)

	http.HandleFunc("/tree", getShortName)   // set router
	err := http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
