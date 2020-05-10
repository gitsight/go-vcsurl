package vcsurl

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

// Errors returned by Parse function.
var (
	ErrUnknownURL  = errors.New("unknown URL format")
	ErrUnableParse = errors.New("unable to determine name or full name")
	ErrEmptyURL    = errors.New("empty URL")
	ErrEmptyPath   = errors.New("empty path in URL")
)

// Host VCS provider.
type Host string

// Supported VCS host provider.
const (
	GitHub    Host = "github.com"
	Bitbucket Host = "bitbucket.org"
	GitLab    Host = "gitlab.com"

	gitHubAPI Host = "api.github.com"
)

// Kind of VCS
type Kind string

// Supported VCS kinds.
const (
	Git Kind = "git"
)

// Protocol of remote
type Protocol string

// Supported VCS protocols.
const (
	SSH   Protocol = "ssh"
	HTTPS Protocol = "https"
)

var kindByHost = map[Host]Kind{
	GitHub:    Git,
	gitHubAPI: Git,
	GitLab:    Git,
	Bitbucket: Git,
}

// VCS describes a VCS repository.
type VCS struct {
	// ID unique repository identification.
	ID string
	// CloneURL git remote format.
	CloneURL string
	// Kind of VCS.
	Kind Kind
	// Host is the public web of the repository.
	Host Host
	// Username of repo owner on repo hosting site.
	Username string
	// Name base name of repo on repo hosting site.
	Name string
	// FullName full name of repo on repo hosting site.
	FullName string
	// Committish is a reference to an object that can be recursively
	// dereferenced to a commit object. They can be commits, tags or branches.
	Committish string
}

var (
	removeDotGit    = regexp.MustCompile(`\.git$`)
	gitPreprocessRE = regexp.MustCompile("^git@([a-zA-Z0-9-_\\.]+)\\:(.*)$")
)

// Parse parses a string that resembles a VCS repository URL. See TestParse for
// a list of supported URL formats.
func Parse(raw string) (*VCS, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("parse %q: %w", raw, ErrEmptyURL)
	}

	spec := raw
	if parts := gitPreprocessRE.FindStringSubmatch(spec); len(parts) == 3 {
		spec = fmt.Sprintf("git://%s/%s", parts[1], parts[2])
	}

	parsedURL, err := url.Parse(spec)
	if err != nil {
		return nil, err
	}

	if parsedURL.Scheme == "" {
		spec = "https://" + spec
		if parsedURL, err = url.Parse(spec); err != nil {
			return nil, err
		}
	}

	vcs := &VCS{}
	vcs.Host = Host(parsedURL.Host)
	vcs.Kind = kindByHost[vcs.Host]

	switch vcs.Host {
	case GitHub, gitHubAPI:
		err = vcs.parseGitHub(parsedURL)
	case Bitbucket:
		err = vcs.parseBitbucket(parsedURL)
	case GitLab:
		err = vcs.parseGitlab(parsedURL)
	default:
		err = vcs.parseDefault(parsedURL)
	}

	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", raw, err)
	}

	if vcs.ID == "" {
		vcs.ID = fmt.Sprintf("%s/%s", string(vcs.Host), vcs.FullName)
	}

	return vcs, nil

}

func (v *VCS) parseGitHub(url *url.URL) error {
	parts := strings.Split(url.Path, "/")
	if v.Host == gitHubAPI {
		v.Host = GitHub
		if len(parts) < 2 || parts[1] != "repos" {
			return ErrUnknownURL
		}

		parts = parts[1:]
	}

	if len(parts) < 3 {
		return ErrUnknownURL
	}

	v.Username = parts[1]
	v.Name = removeDotGit.ReplaceAllLiteralString(parts[2], "")
	v.FullName = v.Username + "/" + v.Name

	if len(parts) < 5 {
		return nil
	}

	if _, ok := githubCommittishParts[parts[3]]; ok {
		v.Committish = strings.Join(parts[4:], "/")
		return nil
	}

	if len(parts) >= 6 && parts[3] == "releases" {
		v.Committish = parts[5]
	}

	return nil
}

var githubCommittishParts = map[string]struct{}{
	"commits":  struct{}{},
	"commit":   struct{}{},
	"tree":     struct{}{},
	"branches": struct{}{},
}

func (v *VCS) parseBitbucket(url *url.URL) error {
	parts := strings.Split(url.Path, "/")
	if len(parts) < 3 {
		return ErrUnknownURL
	}

	v.Username = parts[1]
	v.Name = removeDotGit.ReplaceAllLiteralString(parts[2], "")
	v.FullName = v.Username + "/" + v.Name

	if len(parts) >= 5 && (parts[3] == "src" || parts[3] == "commits" || parts[3] == "branch") {
		v.Committish = parts[4]
	}

	return nil
}

func (v *VCS) parseGitlab(url *url.URL) error {
	parts := strings.Split(url.Path, "/")
	if len(parts) < 3 {
		return ErrUnknownURL
	}

	var last int
	for _, p := range parts {
		if p == "-" {
			break
		}
		last++
	}

	v.Username = strings.Join(parts[1:last-1], "/")
	v.Name = removeDotGit.ReplaceAllLiteralString(parts[last-1], "")
	v.FullName = v.Username + "/" + v.Name

	if len(parts) >= (last + 2) {
		object := parts[last+1]
		if object == "tags" || object == "commit" || object == "tree" {
			v.Committish = strings.Join(parts[last+2:], "/")
		}
	}

	return nil
}

func (v *VCS) parseDefault(url *url.URL) error {
	path := url.Path
	if len(path) == 0 {
		return ErrEmptyPath
	}

	path = path[1:] // remove leading slash
	path = removeDotGit.ReplaceAllLiteralString(path, "")
	v.FullName = path
	v.Name = filepath.Base(path)
	if strings.Contains(url.String(), "git") {
		v.Kind = Git
	}

	if v.Name == "" || v.FullName == "" {
		return ErrUnableParse
	}

	return nil
}

func (v *VCS) Remote(p Protocol) string {
	switch p {
	case SSH:
		return v.sshRemote()
	case HTTPS:
		return v.httpsRemote()
	default:
		return ""
	}
}

// git@gitlab.com:commento/docs.git
// git@github.com:go-git/go-git.git
// git clone git@bitbucket.org:mcuadros/discovery-rest.git
func (v *VCS) sshRemote() string {
	return fmt.Sprintf("git@%s/%s/%s.git", v.Host, v.Username, v.Name)
}

// https://mcuadros@bitbucket.org/mcuadros/discovery-rest.git
// https://gitlab.com/commento/docs.git
// https://github.com/go-git/go-git.git
func (v *VCS) httpsRemote() string {
	var auth string
	if v.Host == Bitbucket {
		auth = fmt.Sprintf("%s@", v.Username)
	}

	return fmt.Sprintf("https://%s%s/%s/%s.git", auth, v.Host, v.Username, v.Name)
}
