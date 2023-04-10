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
		{
			name: "Valid URL with ssh",
			args: args{
				str: "git@github.com:julien040/gut.git",
			},
			want: true,
		},
		{
			name: "Valid URL with ssh without .git",
			args: args{
				str: "git@github.com:julien040/gut",
			},
			want: true,
		},
		{
			name: "Valid URL with ssh without host",
			args: args{
				str: "git@julien040/gut.git",
			},
			want: false,
		},
		{
			name: "Valid URL with ssh without scheme",
			args: args{
				str: "julien040@git/github.com/julien040/gut.git",
			},
			want: false,
		},
		{
			name: "Valid Gitlab URL with ssh",
			args: args{
				str: "git@gitlab.com:gitlab-org/gitlab.git",
			},
			want: true,
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
		{
			name: "Valid URL with ssh",
			args: args{
				str: "git@github.com:julien040/gut.git",
			},
			want: "gut",
		},
		{
			name: "Valid URL with ssh without .git",
			args: args{
				str: "git@github.com:julien040/gut",
			},
			want: "gut",
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
		{
			name: "Valid URL with ssh",
			args: args{
				str: "git@github.com:julien040/gut",
			},
			want: "github.com",
		},
		{
			name: "Invalid URL",
			args: args{
				str: "",
			},
			want: "",
		},
		{
			name: "Valid URL with ssh without path",
			args: args{
				str: "git@github.com:julien040",
			},
			want: "github.com",
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

func Test_isEmailValid(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid email",
			args: args{
				email: "contact@acme.com",
			},
			want: true,
		},
		{
			name: "Invalid email 1",
			args: args{
				email: "contact@acme",
			},
			want: true,
		},
		{
			name: "Invalid email 2",
			args: args{
				email: "contactacme.com",
			},
			want: false,
		},
		{
			name: "Invalid email 3",
			args: args{
				email: "@acme.com",
			},
			want: false,
		},
		{
			name: "Empty email",
			args: args{
				email: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEmailValid(tt.args.email); got != tt.want {
				t.Errorf("isEmailValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isDomainValid(t *testing.T) {
	type args struct {
		domain string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid domain",
			args: args{
				domain: "acme.com",
			},
			want: true,
		},
		{
			name: "Invalid domain",
			args: args{
				domain: "acme",
			},
			want: false,
		},
		{
			name: "Invalid domain",
			args: args{
				domain: "acme.",
			},
			want: false,
		},
		{
			name: "Invalid domain",
			args: args{
				domain: ".com",
			},
			want: false,
		},
		{
			name: "Empty domain",
			args: args{
				domain: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDomainValid(tt.args.domain); got != tt.want {
				t.Errorf("isDomainValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
