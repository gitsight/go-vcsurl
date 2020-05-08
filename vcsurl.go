package vcsurl

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

type Host string

const (
	GitHub     Host = "github.com"
	Bitbucket  Host = "bitbucket.org"
	GoogleCode Host = "code.google.com"
	PythonOrg  Host = "hg.python.org"
	Launchpad  Host = "launchpad.net"
	GitLab     Host = "gitlab.com"

	gitHubAPI Host = "api.github.com"
)

type Kind string

const (
	Git       Kind = "git"
	Mercurial Kind = "hg"
	Bazaar    Kind = "bzr"
)

// VCS describes a VCS repository.
type VCS struct {
	// ID unique repository identification.
	ID string
	// CloneURL git remote format.
	CloneURL string
	// VCS type
	VCS Kind
	// Host is the public web of the repository.
	Host Host
	// Username of repo owner on repo hosting site.
	Username string
	// Name base name of repo on repo hosting site.
	Name string
	// FullName full name of repo on repo hosting site.
	FullName string
	// Rev a specific revision (commit ID, branch, etc.).
	Rev string
}

// Link returns the URL to the repository that is intended for access by humans
// using a Web browser (i.e., not the URL to the API resource).
func (r *VCS) Link() string {
	switch r.Host {
	case GoogleCode:
		return fmt.Sprintf("https://code.google.com/p/%s", r.FullName)
	default:
		return (&url.URL{Scheme: "https", Host: string(r.Host), Path: "/" + r.FullName}).String()
	}
}

var (
	removeDotGit    = regexp.MustCompile(`\.git$`)
	gitPreprocessRE = regexp.MustCompile("^git@([a-zA-Z0-9-_\\.]+)\\:(.*)$")
)

// Parses a string that resembles a VCS repository URL. See TestParse for a list of supported URL
// formats.
func Parse(spec string) (info *VCS, err error) {
	if parts := gitPreprocessRE.FindStringSubmatch(spec); len(parts) == 3 {
		spec = fmt.Sprintf("git://%s/%s", parts[1], parts[2])
	}

	var parsedURL *url.URL
	if parsedURL, err = url.Parse(spec); err == nil {
		if parsedURL.Scheme == "" {
			spec = "https://" + spec
			if parsedURL, err = url.Parse(spec); err != nil {
				return nil, err
			}
		}

		info = new(VCS)

		info.CloneURL = parsedURL.String()
		info.Host = Host(parsedURL.Host)
		info.Rev = parsedURL.Fragment

		if info.Host == GitHub || parsedURL.Scheme == "git" {
			info.VCS = Git
		} else if info.Host == GoogleCode && parsedURL.Scheme == "https" {
			info.VCS = Mercurial
		} else if info.Host == Bitbucket && (parsedURL.Scheme == "https" || parsedURL.Scheme == "http") {
			if !strings.HasSuffix(parsedURL.Path, ".git") {
				info.VCS = Mercurial
			}
		} else if info.Host == Launchpad {
			info.VCS = Bazaar
		}

		path := parsedURL.Path
		switch info.Host {
		case GitHub:
			parts := strings.Split(path, "/")
			if len(parts) >= 3 {
				info.Username = parts[1]
				info.Name = removeDotGit.ReplaceAllLiteralString(parts[2], "")
				info.FullName = info.Username + "/" + info.Name
				info.CloneURL = "git://github.com/" + info.FullName + ".git"
			}
		case gitHubAPI:
			parts := strings.Split(path, "/")
			if len(parts) >= 4 && parts[1] == "repos" {
				info.VCS = Git
				info.Host = GitHub
				info.Username = parts[2]
				info.Name = removeDotGit.ReplaceAllLiteralString(parts[3], "")
				info.FullName = info.Username + "/" + info.Name
				info.CloneURL = "git://github.com/" + info.FullName + ".git"
			}

			if len(parts) >= 6 && parts[4] == "commits" {
				info.Rev = parts[5]
			}

		case GoogleCode:
			parts := strings.Split(path, "/")
			if len(parts) >= 3 && parts[1] == "p" {
				info.Name = parts[2]
				info.FullName = info.Name
				info.CloneURL = "https://code.google.com/p/" + info.FullName
			}
		case PythonOrg:
			parts := strings.Split(path, "/")
			if len(parts) >= 2 {
				info.CloneURL = "http://hg.python.org" + path
				info.VCS = Mercurial
				info.Name = parts[len(parts)-1]
				info.FullName = strings.Join(parts[1:], "/")
			}
		case Bitbucket:
			parts := strings.Split(path, "/")
			if len(parts) >= 3 {
				info.Username = parts[1]
				if strings.HasSuffix(parts[2], ".git") {
					info.VCS = Git
					parts[2] = strings.TrimSuffix(parts[2], ".git")
				}
				info.Name = parts[2]
				info.FullName = parts[1] + "/" + parts[2]
				info.CloneURL = "https://bitbucket.org/" + info.FullName
				if info.VCS == Git {
					info.CloneURL += ".git"
				}
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
				info.VCS = Git
			} else if strings.Contains(spec, "hg") || strings.Contains(spec, "mercurial") {
				info.VCS = Mercurial
			}
		}

		if info.Host == Launchpad {
			parsedURL.Scheme = "bzr"
			info.CloneURL = parsedURL.String()
		}

		if info.Name == "" || info.FullName == "" {
			return nil, fmt.Errorf("unable to determine name or full name for repo spec %q", spec)
		}

		if info.ID == "" {
			info.ID = fmt.Sprintf("%s/%s", string(info.Host), info.FullName)
		}

		return info, nil
	}
	return nil, err
}
