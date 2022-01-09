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

func getDir(path string, sortby, direction int) ([]fs.FileInfo, error) {
	absPath, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	fmt.Println("pwd:", absPath)
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
	// files, err = getDir("../", NOTSORT, ASC)
	files, err = getDir("../", NAME, ASC)
	// files, err = getDir("../", NAME, DESC)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
}
