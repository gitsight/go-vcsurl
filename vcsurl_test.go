package vcsurl_test

import (
	"testing"

	"github.com/gitsight/go-vcsurl"
	"github.com/stretchr/testify/require"
)

func TestParse_GitHub(t *testing.T) {
	urls := []string{
		"github.com/foo/bar",
		"http://github.com/foo/bar",
		"https://github.com/foo/bar",
		"https://github.com/foo/bar.git",
		"https://api.github.com/repos/foo/bar",
		"git@github.com:foo/bar",
		"git@github.com:foo/bar.git",
		"git+ssh://github.com/foo/bar",
	}

	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			vcs, err := vcsurl.Parse(url)
			require.NoError(t, err)
			AssertVCS_GitHub(t, vcs)
		})
	}
}

func TestParse_GitHubCommittish(t *testing.T) {
	urls := []string{
		"github.com/foo/bar#qux",
		"https://github.com/foo/bar/commit/qux",
		"https://api.github.com/repos/foo/bar/commits/qux",
		"https://api.github.com/repos/foo/bar/branches/qux",
		"https://github.com/foo/bar/tree/qux",
		"https://github.com/foo/bar/releases/tag/qux",
	}

	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			vcs, err := vcsurl.Parse(url)
			require.NoError(t, err)
			require.Equal(t, vcs.Committish, "qux")
			AssertVCS_GitHub(t, vcs)
		})
	}
}

func TestParse_GitHubCommittishSlash(t *testing.T) {
	vcs, err := vcsurl.Parse("https://github.com/foo/bar/tree/qux/baz")
	require.NoError(t, err)
	require.Equal(t, vcs.Committish, "qux/baz")
	AssertVCS_GitHub(t, vcs)
}

func AssertVCS_GitHub(t *testing.T, vcs *vcsurl.VCS) {
	t.Helper()
	require.Equal(t, vcs.Kind, vcsurl.Git)
	require.Equal(t, vcs.Host, vcsurl.GitHub)
	require.Equal(t, vcs.Username, "foo")
	require.Equal(t, vcs.Name, "bar")
	require.Equal(t, vcs.FullName, "foo/bar")
}

func TestParse_Bitbucket(t *testing.T) {
	urls := []string{
		"bitbucket.org/foo/bar",
		"https://bitbucket.org/foo/bar",
		"http://bitbucket.org/foo/bar",
		"http://bitbucket.org/foo/bar.git",
		"https://baz@bitbucket.org/foo/bar.git",
		"git@bitbucket.org:foo/bar.git",
	}

	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			vcs, err := vcsurl.Parse(url)
			require.NoError(t, err)
			AssertVCS_Bitbucket(t, vcs)
		})
	}
}

func TestParse_BitbucketCommittish(t *testing.T) {
	urls := []string{
		"bitbucket.org/foo/bar#qux",
		"https://bitbucket.org/foo/bar/src/qux/",
		"https://bitbucket.org/foo/bar/commits/qux",
		"https://bitbucket.org/foo/bar/branch/qux",
	}

	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			vcs, err := vcsurl.Parse(url)
			require.NoError(t, err)
			require.Equal(t, vcs.Committish, "qux")
			AssertVCS_Bitbucket(t, vcs)
		})
	}
}

//

func AssertVCS_Bitbucket(t *testing.T, vcs *vcsurl.VCS) {
	t.Helper()
	require.Equal(t, vcs.Kind, vcsurl.Git)
	require.Equal(t, vcs.Host, vcsurl.Bitbucket)
	require.Equal(t, vcs.Username, "foo")
	require.Equal(t, vcs.Name, "bar")
	require.Equal(t, vcs.FullName, "foo/bar")
}

func TestParse_Gitlab(t *testing.T) {
	urls := []string{
		"gitlab.com/foo/bar",
		"https://gitlab.com/foo/bar",
		"https://gitlab.com/foo/bar.git",
		"git@gitlab.com:foo/bar.git",
	}

	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			vcs, err := vcsurl.Parse(url)
			require.NoError(t, err)
			require.Equal(t, vcs.Kind, vcsurl.Git)
			require.Equal(t, vcs.Host, vcsurl.GitLab)
			require.Equal(t, vcs.Username, "foo")
			require.Equal(t, vcs.Name, "bar")
			require.Equal(t, vcs.FullName, "foo/bar")
		})
	}
}

func TestParse_GitlabSubGroup(t *testing.T) {
	urls := []string{
		"gitlab.com/foo/bar/qux",
		"https://gitlab.com/foo/bar/qux",
		"https://gitlab.com/foo/bar/qux.git",
		"git@gitlab.com:foo/bar/qux.git",
	}

	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			vcs, err := vcsurl.Parse(url)
			require.NoError(t, err)
			require.Equal(t, vcs.Kind, vcsurl.Git)
			require.Equal(t, vcs.Host, vcsurl.GitLab)
			require.Equal(t, vcs.Username, "foo/bar")
			require.Equal(t, vcs.Name, "qux")
			require.Equal(t, vcs.FullName, "foo/bar/qux")
		})
	}
}

func TestParse_GitlabCommittish(t *testing.T) {
	urls := []string{
		"https://gitlab.com/foo/qux/bar/-/commit/baz",
		"https://gitlab.com/foo/qux/bar/-/tags/baz",
		"https://gitlab.com/foo/bar/-/tree/baz",
	}

	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			vcs, err := vcsurl.Parse(url)
			require.NoError(t, err)
			require.Equal(t, vcs.Name, "bar")
			require.Equal(t, vcs.Committish, "baz")
		})
	}
}
