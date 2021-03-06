package main

import (
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/c2fo/testify/require"
)

type MockFileInfo struct {
	FileName    string
	IsDirectory bool
}

func (mf MockFileInfo) Name() string       { return mf.FileName }
func (mf MockFileInfo) Size() int64        { return int64(8) }
func (mf MockFileInfo) Mode() os.FileMode  { return os.ModePerm }
func (mf MockFileInfo) ModTime() time.Time { return time.Now() }
func (mf MockFileInfo) IsDir() bool        { return mf.IsDirectory }
func (mf MockFileInfo) Sys() interface{}   { return nil }

func Test_sortByName(t *testing.T) {
	type args struct {
		a fs.FileInfo
		b fs.FileInfo
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test_sortByName_files",
			args: args{
				a: MockFileInfo{FileName: "a", IsDirectory: false},
				b: MockFileInfo{FileName: "b", IsDirectory: false},
			},
			want: true,
		},
		{
			name: "Test_sortByName_files_dir_file",
			args: args{
				a: MockFileInfo{FileName: "a", IsDirectory: true},
				b: MockFileInfo{FileName: "b", IsDirectory: false},
			},
			want: true,
		},
		{
			name: "Test_sortByName_files_file_dir",
			args: args{
				a: MockFileInfo{FileName: "a", IsDirectory: true},
				b: MockFileInfo{FileName: "b", IsDirectory: false},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortByName(tt.args.a, tt.args.b)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_sortByTime(t *testing.T) {
	type args struct {
		a fs.FileInfo
		b fs.FileInfo
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test_sortByTime_files",
			args: args{
				a: MockFileInfo{FileName: "a", IsDirectory: false},
				b: MockFileInfo{FileName: "b", IsDirectory: false},
			},
			want: false,
		},
		{
			name: "Test_sortByTime_files_dir_file",
			args: args{
				a: MockFileInfo{FileName: "a", IsDirectory: true},
				b: MockFileInfo{FileName: "b", IsDirectory: false},
			},
			want: true,
		},
		{
			name: "Test_sortByTime_files_file_dir",
			args: args{
				a: MockFileInfo{FileName: "a", IsDirectory: true},
				b: MockFileInfo{FileName: "b", IsDirectory: false},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortByTime(tt.args.a, tt.args.b)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_getDir(t *testing.T) {
	type args struct {
		path      string
		sortby    int
		direction int
	}
	tests := []struct {
		name string
		args args
		want []fs.FileInfo
	}{
		{
			name: "Test_getDir_name_asc",
			args: args{
				path:      "./",
				sortby:    NAME,
				direction: ASC,
			},
			want: []fs.FileInfo{
				MockFileInfo{FileName: ".git", IsDirectory: true},
				MockFileInfo{FileName: "bin", IsDirectory: true},
				MockFileInfo{FileName: ".gitignore", IsDirectory: false},
				MockFileInfo{FileName: "README.md", IsDirectory: false},
				MockFileInfo{FileName: "go.mod", IsDirectory: false},
				MockFileInfo{FileName: "go.sum", IsDirectory: false},
				MockFileInfo{FileName: "main.go", IsDirectory: false},
				MockFileInfo{FileName: "main_test.go", IsDirectory: false},
			},
		},
		{
			name: "Test_getDir_name_desc",
			args: args{
				path:      "./",
				sortby:    NAME,
				direction: DESC,
			},
			want: []fs.FileInfo{
				MockFileInfo{FileName: "main_test.go", IsDirectory: false},
				MockFileInfo{FileName: "main.go", IsDirectory: false},
				MockFileInfo{FileName: "go.sum", IsDirectory: false},
				MockFileInfo{FileName: "go.mod", IsDirectory: false},
				MockFileInfo{FileName: "README.md", IsDirectory: false},
				MockFileInfo{FileName: ".gitignore", IsDirectory: false},
				MockFileInfo{FileName: "bin", IsDirectory: true},
				MockFileInfo{FileName: ".git", IsDirectory: true},
			},
		},
		{
			name: "Test_getDir_size_asc",
			args: args{
				path:      "./",
				sortby:    SIZE,
				direction: ASC,
			},
			want: []fs.FileInfo{
				MockFileInfo{FileName: ".git", IsDirectory: true},
				MockFileInfo{FileName: "bin", IsDirectory: true},
				MockFileInfo{FileName: ".gitignore", IsDirectory: false},
				MockFileInfo{FileName: "README.md", IsDirectory: false},
				MockFileInfo{FileName: "go.mod", IsDirectory: false},
				MockFileInfo{FileName: "go.sum", IsDirectory: false},
				MockFileInfo{FileName: "main.go", IsDirectory: false},
				MockFileInfo{FileName: "main_test.go", IsDirectory: false},
			},
		},
		{
			name: "Test_getDir_size_desc",
			args: args{
				path:      "./",
				sortby:    SIZE,
				direction: DESC,
			},
			want: []fs.FileInfo{
				MockFileInfo{FileName: "main_test.go", IsDirectory: false},
				MockFileInfo{FileName: "main.go", IsDirectory: false},
				MockFileInfo{FileName: "go.sum", IsDirectory: false},
				MockFileInfo{FileName: "go.mod", IsDirectory: false},
				MockFileInfo{FileName: "README.md", IsDirectory: false},
				MockFileInfo{FileName: ".gitignore", IsDirectory: false},
				MockFileInfo{FileName: ".git", IsDirectory: true},
				MockFileInfo{FileName: "bin", IsDirectory: true},
			},
		},
		{
			name: "Test_getDir_time_asc",
			args: args{
				path:      "./",
				sortby:    TIME,
				direction: ASC,
			},
			want: []fs.FileInfo{
				MockFileInfo{FileName: ".git", IsDirectory: true},
				MockFileInfo{FileName: "bin", IsDirectory: true},
				MockFileInfo{FileName: ".gitignore", IsDirectory: false},
				MockFileInfo{FileName: "README.md", IsDirectory: false},
				MockFileInfo{FileName: "go.mod", IsDirectory: false},
				MockFileInfo{FileName: "go.sum", IsDirectory: false},
				MockFileInfo{FileName: "main.go", IsDirectory: false},
				MockFileInfo{FileName: "main_test.go", IsDirectory: false},
			},
		},
		{
			name: "Test_getDir_time_desc",
			args: args{
				path:      "./",
				sortby:    TIME,
				direction: DESC,
			},
			want: []fs.FileInfo{
				MockFileInfo{FileName: "main_test.go", IsDirectory: false},
				MockFileInfo{FileName: "main.go", IsDirectory: false},
				MockFileInfo{FileName: "go.sum", IsDirectory: false},
				MockFileInfo{FileName: "go.mod", IsDirectory: false},
				MockFileInfo{FileName: "README.md", IsDirectory: false},
				MockFileInfo{FileName: ".gitignore", IsDirectory: false},
				MockFileInfo{FileName: "bin", IsDirectory: true},
				MockFileInfo{FileName: ".git", IsDirectory: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Chtimes("./main_test.go", time.Now(), time.Now())
			got, err := Dir(tt.args.path, tt.args.sortby, tt.args.direction)
			if err != nil {
				t.Error(err)
				return
			}

			for i, f := range got {
				require.Equal(t, tt.want[i].Name(), f.Name())
			}
		})
	}
}
