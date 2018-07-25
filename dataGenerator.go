package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dogenzaka/tsv"
)

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// plantae
// Coniferophyta
func findTrees(path string) map[string]struct{} {
	// func findTrees(path string) []string {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	data := Line{}
	parser, _ := tsv.NewParser(file, &data)

	// var array []string
	set := make(map[string]struct{})
	set2 := make(map[string]struct{})
	for {
		eof, err := parser.Next()
		if eof {
			return set2
			// return array
		}
		check(err)
		// if data.Kingdom == "Plantae" && data.Phylum == "Tracheophyta" {
		if data.Kingdom == "Plantae" {
			// array = append(array, data.TaxonID)
			set[data.TaxonID] = struct{}{}
			set2[data.Order] = struct{}{}
		}
	}
}

func findNames(path string, numbers map[string]struct{}) []string {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	data := Line2{}
	parser, _ := tsv.NewParser(file, &data)

	var array []string

	for {
		eof, err := parser.Next()
		if eof {
			return array
		}
		check(err)
		_, y := numbers[data.TaxonID]
		if data.Language == "en" && y {
			array = append(array, data.VernacularName)
		}
	}
}

func singleNames(array []string) []string {
	var a []string
	set := make(map[string]struct{})
	for i := range array {
		if !strings.ContainsAny(array[i], " '`") {
			set[array[i]] = struct{}{}
			// a = append(a, array[i])
		}
	}
	for i := range set {
		a = append(a, i)
	}
	return a
}

func main() {
	set := findTrees("data/Taxon-copy.tsv")
	fmt.Println(len(set))

	for i := range set {
		fmt.Println(i)
	}
	// set2 := findNames("data/VernacularName-copy.tsv", set)
	// a := singleNames(set2)
	// for i := range a {
	// 	fmt.Println(a[i])
	// }
	// fmt.Println(len(a))
}
