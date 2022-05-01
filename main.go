package main // import "hello-filepath"

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	_       = iota
	NOTSORT // 1 - do not sort
	NAME    // filename
	TYPE    // file type
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
		return strings.ToLower(a.Name()) < strings.ToLower(b.Name())
	}
}

func sortByType(a, b fs.FileInfo) bool {
	switch true {
	case a.IsDir() && !b.IsDir():
		return true
	case !a.IsDir() && b.IsDir():
		return false
	default:
		return strings.ToLower(filepath.Ext(a.Name())) < strings.ToLower(filepath.Ext(b.Name()))
	}
}

func sortBySize(a, b fs.FileInfo) bool {
	switch true {
	case a.IsDir() && !b.IsDir():
		return true
	case !a.IsDir() && b.IsDir():
		return false
	default:
		return a.Size() < b.Size()
	}
}

func sortByTime(a, b fs.FileInfo) bool {
	// switch true {
	// case a.IsDir() && !b.IsDir():
	// 	return true
	// case !a.IsDir() && b.IsDir():
	// 	return false
	// default:
	// 	return a.ModTime().Format("20060102150405") < b.ModTime().Format("20060102150405")
	// }
	return a.ModTime().Format("20060102150405") < b.ModTime().Format("20060102150405")
}

func Dir(path string, sortby, direction int) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if sortby > 1 && sortby < 6 {
		sort.Slice(files, func(a, b int) bool {
			switch sortby {
			case NAME:
				if direction == DESC {
					if (files[a].IsDir() && !files[b].IsDir()) || (!files[a].IsDir() && files[b].IsDir()) {
						return sortByName(files[a], files[b])
					} else {
						return !sortByName(files[a], files[b])
					}
				}
				return sortByName(files[a], files[b])
			case TYPE:
				if direction == DESC {
					if (files[a].IsDir() && !files[b].IsDir()) || (!files[a].IsDir() && files[b].IsDir()) {
						return sortByType(files[a], files[b])
					} else {
						return !sortByType(files[a], files[b])
					}
				}
				return sortByType(files[a], files[b])
			case SIZE:
				if direction == DESC {
					return !sortBySize(files[a], files[b])
				}
				return sortBySize(files[a], files[b])
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

	// path := "../go.mod"
	// path := "../"
	path := "../"
	// path, _ := os.Getwd()

	p, _ := os.Stat(path)
	dir := path
	if p.IsDir() {
		dir = path + "/"
	}

	dir = filepath.Dir(dir)
	absPath, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}

	fmt.Println("pwd:", absPath, path[:len(path)-1], dir)

	// if path != dir {
	// 	panic("path is not dir")
	// }

	// files, err = Dir(absPath, NOTSORT, ASC)
	files, err = Dir(absPath, NAME, ASC)
	// files, err = Dir(absPath, NAME, DESC)
	// files, err = Dir(absPath, TYPE, DESC)
	// files, err = Dir(absPath, SIZE, ASC)
	// files, err = Dir(absPath, SIZE, DESC)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), filepath.Ext(file.Name()), file.ModTime().Format("2006-01-02 15:04:05"), file.IsDir())
	}

	// dlist, err := json.Marshal(files)
	// if err != nil {
	// 	panic(err)
	// }

	// log.Println(string(dlist))
}
