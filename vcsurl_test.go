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

func TestParse_Default(t *testing.T) {
	urls := []struct {
		url []string
		vcs *vcsurl.VCS
	}{{
		[]string{
			"git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git",
			"https://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git",
		}, &vcsurl.VCS{
			Kind:     vcsurl.Git,
			ID:       "git.kernel.org/pub/scm/linux/kernel/git/stable/linux",
			Host:     "git.kernel.org",
			Name:     "linux",
			FullName: "pub/scm/linux/kernel/git/stable/linux",
		},
	}, {
		[]string{"https://kernel.googlesource.com/pub/scm/linux/kernel/git/stable/linux.git"}, &vcsurl.VCS{
			Kind:     vcsurl.Git,
			ID:       "kernel.googlesource.com/pub/scm/linux/kernel/git/stable/linux",
			Host:     "kernel.googlesource.com",
			Name:     "linux",
			FullName: "pub/scm/linux/kernel/git/stable/linux",
		},
	}, {
		[]string{"https://gitea.com/gitea/tea.git"}, &vcsurl.VCS{
			Kind:     vcsurl.Git,
			ID:       "gitea.com/gitea/tea",
			Host:     "gitea.com",
			Name:     "tea",
			FullName: "gitea/tea",
		},
	}, {
		[]string{"git://git.savannah.gnu.org/bash.git"}, &vcsurl.VCS{
			Kind:     vcsurl.Git,
			ID:       "git.savannah.gnu.org/bash",
			Host:     "git.savannah.gnu.org",
			Name:     "bash",
			FullName: "bash",
		},
	}, {
		[]string{"https://git.savannah.gnu.org/git/bash.git"}, &vcsurl.VCS{
			Kind:     vcsurl.Git,
			ID:       "git.savannah.gnu.org/git/bash",
			Host:     "git.savannah.gnu.org",
			Name:     "bash",
			FullName: "git/bash",
		},
	}, {
		[]string{"ssh://git.savannah.gnu.org/srv/git/bash.git"},
		&vcsurl.VCS{
			Kind:     vcsurl.Git,
			ID:       "git.savannah.gnu.org/srv/git/bash",
			Host:     "git.savannah.gnu.org",
			Name:     "bash",
			FullName: "srv/git/bash",
		},
	}}

	for _, test := range urls {
		for _, url := range test.url {
			t.Run(url, func(t *testing.T) {
				vcs, err := vcsurl.Parse(url)
				require.NoError(t, err)
				require.Equal(t, test.vcs, vcs)
			})
		}
	}
}

func TestParse_Empty(t *testing.T) {
	_, err := vcsurl.Parse("")
	require.Error(t, err)
}

func TestParse_Invalid(t *testing.T) {
	_, err := vcsurl.Parse("foo")
	require.Error(t, err)
}
