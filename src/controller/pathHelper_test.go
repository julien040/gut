/* -------------------------------------------------------------------------- */
/*      This file contains function that can be used by multiple commands     */
/* -------------------------------------------------------------------------- */

package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/fatih/color"
)

func Test_getAbsPathFromInput(t *testing.T) {
	wd, _ := os.Getwd()

	type args struct {
		repoName string
		str      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Local path",
			args: args{
				repoName: "https://github.com/julien040/gut.git",
				str:      "test",
			},
			want: filepath.Join(wd, "test"),
		},
		{
			name: "Current folder",
			args: args{
				repoName: "gut",
				str:      ".",
			},
			want: wd,
		},
		{
			name: "Empty string",
			args: args{
				repoName: "gut",
				str:      "",
			},
			want: filepath.Join(wd, "gut"),
		},
		{
			name: "Absolute path",
			args: args{
				repoName: "gut",
				str:      "/home/user/repo",
			},
			want: "/home/user/repo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAbsPathFromInput(tt.args.repoName, tt.args.str); got != tt.want {
				t.Errorf("getAbsPathFromInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkIfPathExist(t *testing.T) {
	wd, _ := os.Getwd()
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Path exists",
			args: args{
				path: "/home",
			},
			want: true,
		},
		{
			name: "Path doesn't exist",
			args: args{
				path: "/home/this/path/doesnt/exist",
			},
			want: false,
		},
		{
			name: "Path is a file",
			args: args{
				path: "go.mod",
			},
			want: false,
		},
		{
			name: "Path is a local folder",
			args: args{
				path: wd,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkIfPathExist(tt.args.path); got != tt.want {
				fmt.Fprintf(color.Output, "Path: %s", tt.args.path)
				t.Errorf("checkIfPathExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isDirectoryEmpty(t *testing.T) {
	os.Mkdir("test", 0755)
	wd, _ := os.Getwd()
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Empty directory",
			args: args{
				path: filepath.Join(wd, "test"),
			},
			want: true,
		},
		{
			name: "Directory with files",
			args: args{
				path: filepath.Join(wd),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isDirectoryEmpty(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("isDirectoryEmpty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isDirectoryEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
	// Remove the test directory
	os.Remove("test")
}
