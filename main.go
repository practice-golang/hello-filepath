package main // import "hello-filepath"

import (
	_ "embed"
	"fmt"
	"io/fs"
	"log"

	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/practice-golang/lorca"

	"gopkg.in/guregu/null.v4"
)

//go:embed index.html
var index []byte

type PathRequest struct {
	Path string `json:"path"`
}

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

func Dir(path string, sortby, direction int) ([]fs.DirEntry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if sortby > 1 && sortby < 6 {
		sort.Slice(files, func(a, b int) bool {
			var fa, fb os.FileInfo

			fa, _ = files[a].Info()
			fb, _ = files[b].Info()

			switch sortby {
			case NAME:
				if direction == DESC {
					if (files[a].IsDir() && !files[b].IsDir()) || (!files[a].IsDir() && files[b].IsDir()) {
						return sortByName(fa, fb)
					} else {
						return !sortByName(fa, fb)
					}
				}

				return sortByName(fa, fb)
			case TYPE:
				if direction == DESC {

					if (files[a].IsDir() && !files[b].IsDir()) || (!files[a].IsDir() && files[b].IsDir()) {
						return sortByType(fa, fb)
					} else {
						return !sortByType(fa, fb)
					}
				}
				return sortByType(fa, fb)
			case SIZE:
				if direction == DESC {
					return !sortBySize(fa, fb)
				}
				return sortBySize(fa, fb)
			case TIME:
				if direction == DESC {
					return !sortByTime(fa, fb)
				}
				return sortByTime(fa, fb)
			default:
				return false
			}
		})
	}

	return files, nil
}

func GetDirectoryList(path string, target, order int) (FileList, error) {
	var err error
	var files []fs.DirEntry = []fs.DirEntry{}

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

	result.FullPath = null.StringFrom(filepath.ToSlash(absPath))

	files, err = Dir(absPath, target, order)

	if err != nil {
		return result, err
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if len(ext) > 0 {
			ext = ext[1:]
		}

		finfo, _ := file.Info()

		fileInfo := FileInfo{
			Name:     null.StringFrom(file.Name()),
			Type:     null.StringFrom(ext),
			Size:     null.IntFrom(finfo.Size()),
			DateTime: null.StringFrom(finfo.ModTime().Format("2006-01-02 15:04:05")),
			DTTM:     null.StringFrom(finfo.ModTime().Format("20060102150405")),
			IsDir:    null.BoolFrom(file.IsDir()),
		}

		result.Files = append(result.Files, fileInfo)
	}

	return result, nil
}

func DirectoryList(c echo.Context) error {
	pathRequest := new(PathRequest)

	if err := c.Bind(pathRequest); err != nil {
		r := make(map[string]interface{})
		r["msg"] = "Wrong request"
		return c.JSON(http.StatusOK, r)
	}

	// fList, err := GetDirectoryList("..", NAME, DESC)
	// fList, err := GetDirectoryList("..", TYPE, ASC)
	// fList, err := GetDirectoryList("..", SIZE, ASC)
	fList, err := GetDirectoryList(pathRequest.Path, NAME, ASC)
	if err != nil {
		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, fList)
}

func initEcho() {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, string(index))
	})
	e.POST("/dir-list", DirectoryList)

	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}

func initLorca() {
	cwd, _ := os.Getwd()
	profilePath := cwd + `\profile`

	lorca.DefaultChromeArgs = []string{
		"--disable-background-networking",
		"--disable-background-timer-throttling",
		"--disable-backgrounding-occluded-windows",
		"--disable-breakpad",
		"--disable-client-side-phishing-detection",
		"--disable-default-apps",
		"--disable-dev-shm-usage",
		"--disable-infobars",
		"--disable-extensions",
		"--disable-features=site-per-process",
		"--disable-hang-monitor",
		"--disable-ipc-flooding-protection",
		"--disable-popup-blocking",
		"--disable-prompt-on-repost",
		"--disable-renderer-backgrounding",
		"--disable-sync",
		"--disable-translate",
		// "--disable-windows10-custom-titlebar",
		"--metrics-recording-only",
		// "--no-first-run",
		"--no-default-browser-check",
		"--safebrowsing-disable-auto-update",
		// "--enable-automation",
		"--password-store=basic",
		"--use-mock-keychain",
	}

	// args := []string{"--ash-force-desktop", "--force-app-mode"}
	args := []string{}

	ui, err := lorca.New("http://localhost:1323", profilePath, 1024, 768, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	<-ui.Done()

	os.RemoveAll(cwd + `/profile`)
}

func main() {
	go func() { initEcho() }()
	initLorca()
}
