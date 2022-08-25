package dir

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/guregu/null.v4"
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
