package dir

import "gopkg.in/guregu/null.v4"

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
	Path     null.String `json:"path"`
	FullPath null.String `json:"full-path"`
	Children []FileInfo  `json:"children"`
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
