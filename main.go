package main // import "hello-filepath"

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
)

func getDir(path string) {
	absPath, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		log.Fatal(err)
	}

	log.Println(absPath)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		log.Println(file.Name(), file.IsDir())
	}
}

func main() {
	getDir("../")
}
