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

// Tree contains info on a species of trees
type Tree struct {
	TaxonIDS    []string `json:"taxonIds"`
	Genus       string   `json:"genus"`
	Species     string   `json:"species"`
	CommonNames []Name   `json:"commonNames"`
}

// Name contains information on a name and its origin
type Name struct {
	Language string `json:"language"`
	Name     string `json:"name"`
}

// Health is a response containing a status on the application
type Health struct {
	Status string `json:"status"`
}

var trees []Tree
var indicies map[string][]int

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

func indexByLang(trees []Tree) map[string][]int {
	arrays := make(map[string][]int)
	for i, t := range trees {
		set := make(map[string]struct{})
		for _, n := range t.CommonNames {
			set[n.Language] = struct{}{}
		}
		for l := range set {
			arrays[l] = append(arrays[l], i)
		}
	}
	return arrays
}

func getRandomTree(trees []Tree, lang string, indicies map[string][]int) Tree {
	i := indicies[lang][rand.Intn(len(indicies[lang]))]
	t := trees[i]
	return t
}

func getShortName(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("language")
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	t := getRandomTree(trees, lang, indicies)
	str, err := json.Marshal(t)
	check(err)
	fmt.Fprintf(w, string(str))
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Requested health")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	t := Health{
		Status: "healthy",
	}
	str, err := json.Marshal(t)
	check(err)
	fmt.Fprintf(w, string(str))
}

func main() {
	datafile := flag.String("data", "/data/data.json", "The file with data in it")
	port := flag.String("port", "80", "The port to listen on")

	flag.Parse()

	fmt.Println("Reading data from: ", *datafile)
	trees = readFromFile(*datafile)
	fmt.Println("Generating indicies...")
	indicies = indexByLang(trees)

	http.HandleFunc("/tree", getShortName) // set router
	http.HandleFunc("/healthz", getHealth)
	fmt.Println("Listening on port " + *port)
	err := http.ListenAndServe(":"+*port, nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
