/* -------------------------------------------------------------------------- */
/*     This file contains function that can be used by multiple commands.     */
/* -------------------------------------------------------------------------- */

package controller

import (
	"testing"
)

func Test_checkURL(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid URL without https",
			args: args{
				str: "github.com/julien040/gut.git",
			},
			want: false,
		},
		{
			name: "Valid URL with https",
			args: args{
				str: "https://github.com/julien040/gut.git",
			},
			want: true,
		},
		{
			name: "Invalid URL",
			args: args{
				str: "",
			},
			want: false,
		},
		// Check URL doesn't support ssh yet
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkURL(tt.args.str); got != tt.want {
				t.Errorf("checkURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRepoNameFromURL(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid URL with .git",
			args: args{
				str: "https://github.com/julien040/gut.git",
			},
			want: "gut",
		},
		{
			name: "Valid URL without .git",
			args: args{
				str: "https://github.com/julien040/gut",
			},
			want: "gut",
		},
		{
			name: "Invalid URL",
			args: args{
				str: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRepoNameFromURL(tt.args.str); got != tt.want {
				t.Errorf("getRepoNameFromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getHost(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid URL from github",
			args: args{
				str: "https://github.com/julien040/gut",
			},
			want: "github.com",
		},
		{
			name: "Valid URL from gitlab",
			args: args{
				str: "https://gitlab.com/",
			},
			want: "gitlab.com",
		},
		{
			name: "Valid URL with port",
			args: args{
				str: "https://gitlab.com:8080/",
			},
			want: "gitlab.com:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHost(tt.args.str); got != tt.want {
				t.Errorf("getHost() = %v, want %v", got, tt.want)
			}
		})
	}
}
