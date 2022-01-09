package main // import "hello-filepath"

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sort"
)

const (
	_       = iota
	NOTSORT // 1 - do not sort
	NAME    // filename
	SIZE    // filesize
	TIME    // filetime
	ASC     // ascending
	DESC    // descending
)

func sortByName(a, b fs.FileInfo) bool {
	switch true {
	case a.IsDir() && !b.IsDir():
		return true
	case !a.IsDir() && b.IsDir():
		return false
	default:
		return a.Name() < b.Name()
	}
}

func sortByTime(a, b fs.FileInfo) bool {
	switch true {
	case a.IsDir() && !b.IsDir():
		return true
	case !a.IsDir() && b.IsDir():
		return false
	default:
		return a.ModTime().Format("20060102150405") < b.ModTime().Format("20060102150405")
	}
}

func Dir(path string, sortby, direction int) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if sortby > 1 && sortby < 5 {
		sort.Slice(files, func(a, b int) bool {
			switch sortby {
			case NAME:
				if direction == DESC {
					return !sortByName(files[a], files[b])
				}
				return sortByName(files[a], files[b])
			case SIZE:
				if direction == DESC {
					return !(files[a].Size() < files[b].Size())
				}
				return files[a].Size() < files[b].Size()
			case TIME:
				if direction == DESC {
					return !sortByTime(files[a], files[b])
				}
				return sortByTime(files[a], files[b])
			default:
				return false
			}
		})
	}

	return files, nil
}

func main() {
	var err error
	var files []fs.FileInfo = []fs.FileInfo{}

	path := "../go.mod"

	dir := filepath.Dir(path)
	absPath, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}

	fmt.Println("pwd:", absPath)

	if path != dir {
		panic("path is not dir")
	}

	// files, err = Dir(path, NOTSORT, ASC)
	files, err = Dir(path, NAME, ASC)
	// files, err = Dir(path, NAME, DESC)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}

	// dlist, err := json.Marshal(files)
	// if err != nil {
	// 	panic(err)
	// }

	// log.Println(string(dlist))
}
