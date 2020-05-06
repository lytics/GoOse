package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/willf/bloom"
)

var (
	comment     = []byte(`#`)
	outtemplate = `// This file was created by goose/precompile
package goose

import (
	"fmt"
	"encoding/json"
	"github.com/willf/bloom"
)

var (
	filters map[string]*bloom.BloomFilter
)

func init() {
	filters = make(map[string]*bloom.BloomFilter)

	rawFilters := map[string]string{ {{ range $key, $value := . }}
	"{{ $key }}": ` + "`{{ $value }}`" + `,{{ end }}
	}

	for lang, rawFilter := range rawFilters {
		filter := &bloom.BloomFilter{}
		err := json.Unmarshal([]byte(rawFilter), filter)
		if err != nil {
			panic(fmt.Sprintf("Invalid bloom filter value; something is terrible: %+v", err))
		}
		filters[lang] = filter
	}
}
`
)

func main() {

	args := os.Args[1:]
	if len(args) != 2 {
		log.Println("Usage: compilefilters <INPUT_PATH> <OUTPUT_FILE>")
		os.Exit(1)
	}

	filters := loadFilters(args[0])
	tmpl, err := template.New("filters").Parse(outtemplate)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	log.Printf("Creating output file: %s", args[1])
	outf, err := os.Create(args[1])
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outf.Close()

	log.Println("Writing template")
	err = tmpl.Execute(outf, filters)
	if err != nil {
		log.Fatalf("Error formatting template: %v", err)
	}
}

func loadFilters(path string) map[string]string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("Error reading input path: %v", err)
	}

	loadFilter := func(filename string) string {
		file, err := os.Open(path + "/" + filename)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}
		defer file.Close()

		filter := bloom.NewWithEstimates(2000, 0.01) // TODO: default params??
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Bytes()
			if bytes.HasPrefix(line, comment) {
				continue
			}
			filter.Add(bytes.TrimSpace(line))
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error scanning file: %v", err)
		}

		filterBytes, err := json.Marshal(filter)
		if err != nil {
			log.Fatalf("Error marshaling filter: %v", err)
		}

		return string(filterBytes)
	}

	filters := make(map[string]string)
	for _, file := range files {
		filename := file.Name()
		fileparts := strings.Split(filename, "_")
		if fileparts[len(fileparts)-1] != "stopwords.txt" {
			continue
		}
		filters[fileparts[0]] = loadFilter(filename)
	}
	return filters
}
