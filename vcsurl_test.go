package vcsurl

import (
	"testing"

	"github.com/kr/pretty"
)

var (
	githubUserRepo = VCS{
		ID:       "github.com/user/repo",
		CloneURL: "git://github.com/user/repo.git",
		VCS:      Git,
		Host:     GitHub,
		Username: "user",
		Name:     "repo",
		FullName: "user/repo",
		Rev:      "asdf",
	}
	googleCodeRepo = VCS{
		ID:       "code.google.com/go",
		CloneURL: "https://code.google.com/p/go",
		VCS:      Mercurial,
		Host:     GoogleCode,
		Name:     "go",
		FullName: "go",
	}
	cpythonRepo = VCS{
		ID:       "hg.python.org/cpython",
		CloneURL: "http://hg.python.org/cpython",
		VCS:      Mercurial,
		Host:     PythonOrg,
		Name:     "cpython",
		FullName: "cpython",
	}
	bitbucketHgRepo = VCS{
		ID:       "bitbucket.org/user/repo",
		CloneURL: "https://bitbucket.org/user/repo",
		VCS:      Mercurial,
		Host:     Bitbucket,
		Username: "user",
		Name:     "repo",
		FullName: "user/repo",
	}
	bitbucketGitRepo = VCS{
		ID:       "bitbucket.org/user/repo",
		CloneURL: "https://bitbucket.org/user/repo.git",
		VCS:      Git,
		Host:     Bitbucket,
		Username: "user",
		Name:     "repo",
		FullName: "user/repo",
	}
	launchpadRepo = VCS{
		ID:       "launchpad.net/repo",
		CloneURL: "bzr://launchpad.net/repo",
		VCS:      Bazaar,
		Host:     Launchpad,
		Username: "",
		Name:     "repo",
		FullName: "repo",
	}
)

func TestParse(t *testing.T) {
	tests := []struct {
		url  string
		rid  string
		info VCS
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
		{"https://api.github.com/repos/user/repo/commits/asdf", "github.com/user/repo", githubUserRepo},

		{"code.google.com/p/go", "code.google.com/p/go", googleCodeRepo},
		{"https://code.google.com/p/go", "code.google.com/p/go", googleCodeRepo},

		{"hg.python.org/cpython", "hg.python.org/cpython", cpythonRepo},
		{"http://hg.python.org/cpython", "hg.python.org/cpython", cpythonRepo},

		{"bitbucket.org/user/repo", "bitbucket.org/user/repo", bitbucketHgRepo},
		{"https://bitbucket.org/user/repo", "bitbucket.org/user/repo", bitbucketHgRepo},
		{"http://bitbucket.org/user/repo", "bitbucket.org/user/repo", bitbucketHgRepo},

		{"bitbucket.org/user/repo.git", "bitbucket.org/user/repo", bitbucketGitRepo},
		{"https://bitbucket.org/user/repo.git", "bitbucket.org/user/repo", bitbucketGitRepo},
		{"http://bitbucket.org/user/repo.git", "bitbucket.org/user/repo", bitbucketGitRepo},

		{"http://launchpad.net/repo", "launchpad.net/repo", launchpadRepo},
		{"bzr://launchpad.net/repo", "launchpad.net/repo", launchpadRepo},
		{"bzr+ssh://launchpad.net/repo", "launchpad.net/repo", launchpadRepo},

		// subpaths
		{"http://github.com/user/repo/subpath#asdf", "github.com/user/repo", githubUserRepo},
		{"git@github.com:user/repo.git/subpath#asdf", "github.com/user/repo", githubUserRepo},
		{"https://code.google.com/p/go/subpath", "code.google.com/p/go", googleCodeRepo},

		// other repo hosts
		{"git://example.com/foo", "example.com/foo", VCS{
			ID:       "example.com/foo",
			CloneURL: "git://example.com/foo",
			VCS:      Git,
			Host:     "example.com",
			Name:     "foo",
			FullName: "foo",
		}},
		{"https://example.com/foo.git", "example.com/foo", VCS{
			ID:       "example.com/foo",
			CloneURL: "https://example.com/foo.git",
			VCS:      Git,
			Host:     "example.com",
			Name:     "foo",
			FullName: "foo",
		}},
		{"https://example.com/git/foo", "example.com/foo", VCS{
			ID:       "example.com/git/foo",
			CloneURL: "https://example.com/git/foo",
			VCS:      Git,
			Host:     "example.com",
			Name:     "foo",
			FullName: "git/foo",
		}},
		{"git@git.private.com:org/repo.git", "git.private.com/org/repo", VCS{
			ID:       "git.private.com/org/repo",
			CloneURL: "git://git.private.com/org/repo.git",
			VCS:      Git,
			Host:     "git.private.com",
			Name:     "repo",
			FullName: "org/repo",
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

func TestLink(t *testing.T) {
	tests := []struct {
		repo VCS
		link string
	}{
		{githubUserRepo, "https://github.com/user/repo"},
		{bitbucketHgRepo, "https://bitbucket.org/user/repo"},
		{bitbucketGitRepo, "https://bitbucket.org/user/repo"},
		{googleCodeRepo, "https://code.google.com/p/go"},
	}

	for _, test := range tests {
		link := test.repo.Link()
		if test.link != link {
			t.Errorf("%s: want link %q, got %q", test.repo.CloneURL, test.link, link)
		}
	}
}
