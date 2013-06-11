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

		// subpaths
		{"http://github.com/user/repo/subpath#asdf", "github.com/user/repo", githubUserRepo},
		{"git@github.com:user/repo.git/subpath#asdf", "github.com/user/repo", githubUserRepo},
		{"https://code.google.com/p/go/subpath", "code.google.com/p/go", googleCodeRepo},

		// other repo hosts
		{"git://example.com/foo", "example.com/foo", RepoInfo{
			CloneURL: "git://example.com/foo",
			VCS:      Git,
			RepoHost: "example.com",
			Name:     "foo",
			FullName: "foo",
		}},
		{"https://example.com/foo.git", "example.com/foo", RepoInfo{
			CloneURL: "https://example.com/foo.git",
			VCS:      Git,
			RepoHost: "example.com",
			Name:     "foo",
			FullName: "foo",
		}},
		{"https://example.com/git/foo", "example.com/foo", RepoInfo{
			CloneURL: "https://example.com/git/foo",
			VCS:      Git,
			RepoHost: "example.com",
			Name:     "foo",
			FullName: "git/foo",
		}},
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
