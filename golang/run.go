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
	"time"
)

// Tree struct to contain information on a tree species
type Tree struct {
	Genus      string `json:"genus"`
	Species    string `json:"species"`
	CommonName string `json:"commonName"`
}

var trees = readFromFile("global_tree_search_trees_1_2.csv")
var keys = make(map[Tree]int)
var names = make(map[int]string)

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
			Genus:      str[0],
			Species:    str[1],
			CommonName: "",
		})
	}
	file.Close()
	return coll
}

func compileMap(t []Tree) map[string][]Tree {
	treeMap := make(map[string][]Tree)
	for _, v := range t {
		treeMap[v.Genus] = append(treeMap[v.Genus], v)
	}
	return treeMap
}

var client = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	fmt.Println("getting: " + url)
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

type Response struct {
	Offset       int  `json:"offset"`
	Limit        int  `json:"limit"`
	EndOfRecords bool `json:"endOfRecords"`
	Results      []struct {
		Key int `json:"key"`
	} `json:"results"`
}

type NameResponse struct {
	Offset       int  `json:"offset"`
	Limit        int  `json:"limit"`
	EndOfRecords bool `json:"endOfRecords"`
	Results      []struct {
		VernacularName string `json:"vernacularName"`
		Language       string `json:"language"`
	} `json:"results"`
}

func getKey(tree Tree) int {
	val, ok := keys[tree]
	if !ok {
		res := new(Response)
		getJson(fmt.Sprintf("http://api.gbif.org/v1/species?name=%s%s%s", tree.Genus, "%20", tree.Species), res)
		if len(res.Results) > 0 {
			val = res.Results[0].Key
		} else {
			val = -1
		}
		keys[tree] = val
	}
	return val
}

func getNames(key int) string {
	val, ok := names[key]
	if !ok {
		res := new(NameResponse)
		getJson(fmt.Sprintf("http://api.gbif.org/v1/species/%v/vernacularNames", key), res)
		val = ""
		for _, v := range res.Results {
			fmt.Print(v)
			if v.Language == "eng" {
				val = v.VernacularName
				break
			}
		}
		names[key] = val
	}
	return val
}

func getShortName(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("path", r.URL.Path)
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	w.Header().Set("Content-Type", "application/json")
	t := trees[rand.Intn(len(trees))]
	fmt.Println(getKey(t))
	fmt.Println(getNames(getKey(t)))
	t.CommonName = getNames(getKey(t))
	str, err := json.Marshal(t)
	check(err)
	fmt.Fprintf(w, string(str))
}

func main() {
	http.HandleFunc("/tree", getShortName) // set router
	err := http.ListenAndServe(":80", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
