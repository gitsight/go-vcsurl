package vcsurl

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	ErrUnknownURL = errors.New("unknown URL format")
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
func Parse(spec string) (*VCS, error) {
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

	info := &VCS{}

	info.CloneURL = parsedURL.String()
	info.Host = Host(parsedURL.Host)
	info.Committish = parsedURL.Fragment
	info.Kind = kindByHost[info.Host]

	path := parsedURL.Path
	switch info.Host {
	case GitHub, gitHubAPI:
		if err := info.parseGitHub(parsedURL); err != nil {
			return nil, err
		}

	case Bitbucket:
		if err := info.parseBitbucket(parsedURL); err != nil {
			return nil, err
		}
	case GitLab:
		if err := info.parseGitlab(parsedURL); err != nil {
			return nil, err
		}
	default:
		if len(path) == 0 {
			return nil, fmt.Errorf("empty path in repo spec: %q", spec)
		}
		path = path[1:] // remove leading slash
		path = removeDotGit.ReplaceAllLiteralString(path, "")
		info.FullName = path
		info.Name = filepath.Base(path)
		if strings.Contains(spec, "git") {
			info.Kind = Git
		}
	}

	if info.Name == "" || info.FullName == "" {
		return nil, fmt.Errorf("unable to determine name or full name for repo spec %q", spec)
	}

	if info.ID == "" {
		info.ID = fmt.Sprintf("%s/%s", string(info.Host), info.FullName)
	}

	return info, nil

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
