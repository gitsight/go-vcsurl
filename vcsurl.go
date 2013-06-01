package vcsurl

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

type RepoHost string

const (
	Github     RepoHost = "github.com"
	GoogleCode RepoHost = "code.google.com"
)

type VCS string

const (
	Git       VCS = "git"
	Mercurial VCS = "hg"
)

// RepoInfo describes a VCS repository.
type RepoInfo struct {
	CloneURL string   // clone URL
	VCS      VCS      // VCS type
	RepoHost RepoHost // repo hosting site
	Username string   // username of repo owner on repo hosting site
	Name     string   // base name of repo on repo hosting site
	FullName string   // full name of repo on repo hosting site
	Rev      string   // a specific revision (commit ID, branch, etc.)
}

var removeDotGit = regexp.MustCompile(`\.git$`)

// Parses a string that resembles a VCS repository URL. See TestParse for a list of supported URL
// formats.
func Parse(spec string) (info *RepoInfo, err error) {
	if strings.HasPrefix(spec, "git@github.com:") {
		spec = strings.Replace(spec, "git@github.com:", "git://github.com/", 1)
	}

	var parsedURL *url.URL
	if parsedURL, err = url.Parse(spec); err == nil {
		if parsedURL.Scheme == "" {
			spec = "https://" + spec
			if parsedURL, err = url.Parse(spec); err != nil {
				return nil, err
			}
		}

		info = new(RepoInfo)

		info.CloneURL = parsedURL.String()
		info.RepoHost = RepoHost(parsedURL.Host)
		info.Rev = parsedURL.Fragment

		if info.RepoHost == Github || parsedURL.Scheme == "git" {
			info.VCS = Git
		} else if info.RepoHost == GoogleCode && parsedURL.Scheme == "https" {
			info.VCS = Mercurial
		}

		path := parsedURL.Path
		switch info.RepoHost {
		case Github:
			path = removeDotGit.ReplaceAllLiteralString(path, "")
			parts := strings.Split(path, "/")
			if len(parts) == 3 {
				info.Username = parts[1]
				info.Name = parts[2]
				info.FullName = parts[1] + "/" + parts[2]
				info.CloneURL = "git://github.com/" + info.FullName + ".git"
			}
		case GoogleCode:
			prefix := "/p/"
			if strings.HasPrefix(path, prefix) {
				info.Name = path[len(prefix):]
				info.FullName = info.Name
			}
		default:
			if len(path) == 0 {
				return nil, fmt.Errorf("empty path in repo spec: %q", spec)
			}
			info.FullName = path[1:] // remove leading slash
			info.Name = filepath.Base(path)
		}

		if info.Name == "" || info.FullName == "" {
			return nil, fmt.Errorf("unable to determine name or full name for repo spec %q", spec)
		}

		return info, nil
	}
	return nil, err
}
