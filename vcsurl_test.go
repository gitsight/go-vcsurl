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
		vcs, err := vcsurl.Parse(url)
		require.NoError(t, err)
		AssertVCS_GitHub(t, vcs)
	}
}

func TestParse_GitHubRevision(t *testing.T) {
	urls := []string{
		"github.com/foo/bar#qux",
		"https://github.com/foo/bar/commit/qux",
		"https://api.github.com/repos/foo/bar/commits/qux",
	}

	for _, url := range urls {
		vcs, err := vcsurl.Parse(url)
		require.NoError(t, err)
		require.Equal(t, vcs.Rev, "qux")
		AssertVCS_GitHub(t, vcs)
	}
}

func AssertVCS_GitHub(t *testing.T, vcs *vcsurl.VCS) {
	t.Helper()
	require.Equal(t, vcs.Kind, vcsurl.Git)
	require.Equal(t, vcs.Host, vcsurl.GitHub)
	require.Equal(t, vcs.Username, "foo")
	require.Equal(t, vcs.Name, "bar")
	require.Equal(t, vcs.FullName, "foo/bar")
}
