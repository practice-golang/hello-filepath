package main // import "hello-filepath"

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gopkg.in/guregu/null.v4"
)

type FilePath struct {
	Path  null.String `json:"path"`
	Sort  null.String `json:"sort"`
	Order null.String `json:"order"`
}

type FileInfo struct {
	Name     null.String `json:"name"`
	Type     null.String `json:"type"`
	Size     null.Int    `json:"size"`
	DateTime null.String `json:"datetime"`
	DTTM     null.String `json:"dttm"`
	IsDir    null.Bool   `json:"isdir"`
}

type FileList struct {
	Path     null.String `json:"path"`
	FullPath null.String `json:"full-path"`
	Files    []FileInfo  `json:"files"`
}

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

func GetDirectoryList(path string, target, order int) (FileList, error) {
	var err error
	var files []fs.FileInfo = []fs.FileInfo{}

	result := FileList{
		Path:     null.StringFrom(path),
		FullPath: null.StringFrom(""),
		Files:    []FileInfo{},
	}

	// path, _ := os.Getwd()
	// path := "../"

	p, _ := os.Stat(path)
	dir := path
	if p.IsDir() {
		dir = path + "/"
	}

	dir = filepath.Dir(dir)
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return result, err
	}

	result.FullPath = null.StringFrom(absPath)

	files, err = Dir(absPath, target, order)

	if err != nil {
		return result, err
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if len(ext) > 0 {
			ext = ext[1:]
		}
		fileInfo := FileInfo{
			Name:     null.StringFrom(file.Name()),
			Type:     null.StringFrom(ext),
			Size:     null.IntFrom(file.Size()),
			DateTime: null.StringFrom(file.ModTime().Format("2006-01-02 15:04:05")),
			DTTM:     null.StringFrom(file.ModTime().Format("20060102150405")),
			IsDir:    null.BoolFrom(file.IsDir()),
		}

		result.Files = append(result.Files, fileInfo)
	}

	return result, nil
}

func DirectoryList(c echo.Context) error {
	// fList, err := GetDirectoryList("..", NAME, DESC)
	// fList, err := GetDirectoryList("..", TYPE, ASC)
	// fList, err := GetDirectoryList("..", SIZE, ASC)
	fList, err := GetDirectoryList("..", NAME, ASC)
	if err != nil {
		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, fList)
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/dir-list", DirectoryList)
	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}
