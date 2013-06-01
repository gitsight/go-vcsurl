package vcsurl

import (
	"github.com/kr/pretty"
	"testing"
)

func TestParse(t *testing.T) {
	githubUserRepo := RepoInfo{
		CloneURL: "git://github.com/user/repo.git",
		VCS:      Git,
		RepoHost: Github,
		Username: "user",
		Name:     "repo",
		FullName: "user/repo",
		Rev:      "asdf",
	}
	googleCodeRepo := RepoInfo{
		CloneURL: "https://code.google.com/p/go",
		VCS:      Mercurial,
		RepoHost: GoogleCode,
		Name:     "go",
		FullName: "go",
	}
	tests := []struct {
		url  string
		rid  string
		info RepoInfo
	}{
		{"github.com/user/repo#asdf", "github.com/user/repo", githubUserRepo},
		{"http://github.com/user/repo#asdf", "github.com/user/repo", githubUserRepo},
		{"http://github.com/user/repo.git#asdf", "github.com/user/repo", githubUserRepo},
		{"https://github.com/user/repo#asdf", "github.com/user/repo", githubUserRepo},
		{"https://github.com/user/repo.git#asdf", "github.com/user/repo", githubUserRepo},
		{"git://github.com/user/repo#asdf", "github.com/user/repo", githubUserRepo},
		{"git://github.com/user/repo.git#asdf", "github.com/user/repo", githubUserRepo},
		{"git+ssh://github.com/user/repo#asdf", "github.com/user/repo", githubUserRepo},
		{"git+ssh://github.com/user/repo.git#asdf", "github.com/user/repo", githubUserRepo},
		{"git@github.com:user/repo#asdf", "github.com/user/repo", githubUserRepo},
		{"git@github.com:user/repo.git#asdf", "github.com/user/repo", githubUserRepo},

		{"code.google.com/p/go", "code.google.com/p/go", googleCodeRepo},
		{"https://code.google.com/p/go", "code.google.com/p/go", googleCodeRepo},
	}

	for _, test := range tests {
		info, err := Parse(test.url)
		if err != nil {
			t.Errorf("clone URL %q: got error: %s", test.url, err)
			continue
		}
		if test.info != *info {
			t.Errorf("%s: %v", test.url, pretty.Diff(test.info, *info))
		}
	}
}
