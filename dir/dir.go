package dir

import (
	"errors"
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

func GetFileList(path string, target, order int) (FileInfo, error) {
	var err error
	var files []fs.DirEntry = []fs.DirEntry{}

	result := FileInfo{
		Path:     null.StringFrom(path),
		FullPath: null.StringFrom(""),
		Children: []FileInfo{},
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
			Path:     null.StringFrom(filepath.ToSlash(filepath.Join(absPath, file.Name()))),
			FullPath: null.StringFrom(filepath.ToSlash(filepath.Join(absPath, file.Name()))),
		}

		result.Children = append(result.Children, fileInfo)
	}

	return result, nil
}

func findUpstreamDirectories(cwd FileInfo, target, order int) (FileInfo, error) {
	var err error

	errHereIsRoot := errors.New("here is root")
	parentPath := filepath.ToSlash(filepath.Dir(cwd.FullPath.String))

	if parentPath == cwd.FullPath.String {
		return FileInfo{}, errHereIsRoot
	}

	p, _ := os.Stat(parentPath)
	if !p.IsDir() {
		return FileInfo{}, errors.New("not a directory")
	}

	parent := FileInfo{
		Name:     null.StringFrom(filepath.Base(parentPath)),
		Path:     null.StringFrom(parentPath),
		FullPath: null.StringFrom(""),
		Children: []FileInfo{},
	}

	parent.Children = append(parent.Children, cwd)
	parent.FullPath = null.StringFrom(parentPath)
	parent.IsDir = null.BoolFrom(true)

	parent.Children, err = GetChildDirectories(parent.FullPath.String, target, order)
	if err != nil {
		return FileInfo{}, err
	}

	for i, child := range parent.Children {

		if filepath.ToSlash(child.Name.String) == filepath.ToSlash(cwd.Name.String) {
			parent.Children[i] = cwd
			break
		}
	}

	grandParent, err := findUpstreamDirectories(parent, target, order)
	if err != nil {
		if err.Error() == errHereIsRoot.Error() {
			return parent, nil
		} else {
			return FileInfo{}, err
		}
	}

	return grandParent, nil

}

func GetVisibleDirectories(path string, target, order int) (FileInfo, error) {
	var err error

	cwd := FileInfo{
		Path:     null.StringFrom(path),
		FullPath: null.StringFrom(""),
		Children: []FileInfo{},
	}

	p, _ := os.Stat(path)
	dir := path
	if !p.IsDir() {
		return cwd, errors.New("not a directory")
	}

	dir = path + "/"
	dir = filepath.Dir(dir)
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return cwd, err
	}

	cwd.Name = null.StringFrom(filepath.Base(absPath))
	cwd.FullPath = null.StringFrom(filepath.ToSlash(absPath))
	cwd.IsDir = null.BoolFrom(true)

	uplist, err := findUpstreamDirectories(cwd, target, order)
	if err != nil {
		return cwd, err
	}

	return uplist, nil
}

func GetChildDirectories(path string, target, order int) (result []FileInfo, err error) {
	p, _ := os.Stat(path)
	dir := path
	if p.IsDir() {
		dir = path + "/"
	}

	dir = filepath.Dir(dir)
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return
	}

	files, err := Dir(absPath, target, order)
	if err != nil {
		return
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
			Path:     null.StringFrom(filepath.ToSlash(filepath.Join(absPath, file.Name()))),
			FullPath: null.StringFrom(filepath.ToSlash(filepath.Join(absPath, file.Name()))),
		}

		if fileInfo.IsDir.Bool {
			result = append(result, fileInfo)
		}
	}

	return
}
