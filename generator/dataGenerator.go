package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/dogenzaka/tsv"
)

// Line to store information about a species
type Line struct {
	TaxonID                  string `tsv:"taxonID"`
	DatasetID                string `tsv:"datasetID"`
	ParentNameUsageID        string `tsv:"parentNameUsageID"`
	AcceptedNameUsageID      string `tsv:"acceptedNameUsageID"`
	OriginalNameUsageID      string `tsv:"originalNameUsageID"`
	ScientificName           string `tsv:"scientificName"`
	ScientificNameAuthorship string `tsv:"scientificNameAuthorship"`
	CanonicalName            string `tsv:"canonicalName"`
	GenericName              string `tsv:"genericName"`
	SpecificEpithet          string `tsv:"specificEpithet"`
	InfraspecificEpithet     string `tsv:"infraspecificEpithet"`
	TaxonRank                string `tsv:"taxonRank"`
	NameAccordingTo          string `tsv:"nameAccordingTo"`
	NamePublishedIn          string `tsv:"namePublishedIn"`
	TaxonomicStatus          string `tsv:"taxonomicStatus"`
	NomenclaturalStatus      string `tsv:"nomenclaturalStatus"`
	TaxonRemarks             string `tsv:"taxonRemarks"`
	Kingdom                  string `tsv:"kingdom"`
	Phylum                   string `tsv:"phylum"`
	Class                    string `tsv:"class"`
	Order                    string `tsv:"order"`
	Family                   string `tsv:"family"`
	Genus                    string `tsv:"genus"`
}

// Line2 to store information about vernacular names
type Line2 struct {
	TaxonID        string `tsv:"taxonID"`
	VernacularName string `tsv:"vernacularName"`
	Language       string `tsv:"language"`
	Country        string `tsv:"country"`
	CountryCode    string `tsv:"countryCode"`
	Sex            string `tsv:"sex"`
	LifeStage      string `tsv:"lifeStage"`
	Source         string `tsv:"source"`
}

// Tree contains information about a species of trees
type Tree struct {
	TaxonIDS    []string `json:"taxonIds"`
	Genus       string   `json:"genus"`
	Species     string   `json:"species"`
	CommonNames []Name   `json:"commonNames"`
}

// Name contains information about a name and its origin
type Name struct {
	Language string `json:"language"`
	Name     string `json:"name"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func updateText(f float32) {
	fmt.Print("\r|")
	for i := 0; i <= 50; i++ {
		if f*100 >= float32(i*2) {
			fmt.Print("=")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Print("|")
	fmt.Printf(" % 5.1f%s", f*100, "%")
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func readLines(path string) int {
	file, err := os.Open(path)
	defer file.Close()
	check(err)
	lines, _ := lineCounter(file)
	return lines
}

func readPlants(path string) map[string][]string {
	lines := readLines(path)
	file, err := os.Open(path)
	defer file.Close()
	check(err)
	data := Line{}
	parser, _ := tsv.NewParser(file, &data)
	set := make(map[string][]string)
	reg, _ := regexp.Compile("[^a-zA-Z0-9 ]+")
	var i int
	for {
		i++
		if i%(lines/1000) == 0 {
			updateText(float32(i) / float32(lines))
		}
		eof, err := parser.Next()
		if eof {
			updateText(1)
			fmt.Println()
			return set
		}
		check(err)
		if data.Kingdom == "Plantae" {
			set[strings.ToLower(reg.ReplaceAllString(strings.Replace(data.CanonicalName, "-", " ", -1), ""))] = append(set[strings.ToLower(data.CanonicalName)], data.TaxonID)
		}
	}
}

func readTrees(path string, ids map[string][]string, names map[string][]Name) []Tree {
	lines := readLines(path)
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	var coll []Tree
	var i int
	for {
		i++
		if i%(lines/1000) == 0 {
			updateText(float32(i) / float32(lines))
		}
		line, error := reader.Read()
		if error == io.EOF {
			updateText(1)
			fmt.Println()
			break
		}
		check(error)
		str := strings.Fields(line[0])
		tIds := getIDS(strings.ToLower(str[0]+" "+str[1]), ids)
		coll = append(coll, Tree{
			TaxonIDS:    tIds,
			Genus:       str[0],
			Species:     str[1],
			CommonNames: getNames(tIds, names),
		})
	}
	file.Close()
	return coll
}

func getIDS(str string, ids map[string][]string) []string {
	s, b := ids[str]
	if b {
		return s
	}
	return make([]string, 0)
}

func getNames(ids []string, names map[string][]Name) []Name {
	var result []Name
	for _, i := range ids {
		ns, b := names[i]
		if b {
			for _, n := range ns {
				result = append(result, n)
			}
		}
	}
	return result

}

// Any Check whether any record is true with predicate
func Any(vs []Name, f func(Name) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func readNames(path string) map[string][]Name {
	lines := readLines(path)
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	data := Line2{}
	parser, _ := tsv.NewParser(file, &data)
	reg, _ := regexp.Compile("[^a-øA-Ø0-9 ]+")
	result := make(map[string][]Name)
	var i int
	for {
		i++
		if i%(lines/1000) == 0 {
			updateText(float32(i) / float32(lines))
		}
		eof, err := parser.Next()
		if eof {
			updateText(1)
			fmt.Println()
			return result
		}
		check(err)
		str := strings.ToLower(reg.ReplaceAllString(strings.Replace(data.VernacularName, "-", " ", -1), ""))
		if !Any(result[data.TaxonID], func(v Name) bool { return str == v.Name && data.Language == v.Language }) {
			result[data.TaxonID] = append(result[data.TaxonID], Name{
				Language: data.Language,
				Name:     str,
			})
		}
	}
}

func removeEmpty(trees []Tree) []Tree {
	var result []Tree
	for _, t := range trees {
		if len(t.CommonNames) > 0 {
			result = append(result, t)
		}
	}
	return result
}

func main() {
	datafile := flag.String("taxon", "Taxon.tsv", "Main dataset")
	treefile := flag.String("trees", "global_tree_search_trees_1_2.csv", "Dataset with trees")
	namefile := flag.String("names", "VernacularName.tsv", "Names dataset")
	output := flag.String("output", "output.json", "File with output data")

	flag.Parse()

	fmt.Println("Reading plant records from " + *datafile)
	plants := readPlants(*datafile)
	fmt.Println("")
	fmt.Println("Reading vernacular names from " + *namefile)
	treeNames := readNames(*namefile)
	fmt.Println("")
	fmt.Println("Collecting trees from " + *treefile)
	trees := removeEmpty(readTrees(*treefile, plants, treeNames))
	str, err := json.Marshal(trees)
	check(err)
	fmt.Println("")
	fmt.Println("Writing data to " + *output)
	err = ioutil.WriteFile(*output, str, 0644)
	check(err)
}
