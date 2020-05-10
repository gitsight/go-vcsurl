# go-vcsurl [![GoDoc](https://godoc.org/github.com/gitsight/go-vcsurl?status.svg)](https://pkg.go.dev/github.com/gitsight/go-vcsurl) [![Test](https://github.com/gitsight/go-vcsurl/workflows/Test/badge.svg)](https://github.com/gitsight/go-vcsurl/actions?query=workflow%3ATest)


`go-vcsurl` library provides a VCS URL parser for HTTP, git or ssh remote URLs
and also frontend URLs from providers like GitHub, GitLab, and Bitbucket. 

This library is based on the previous work done by [@sourcegraph](https://github.com/gitsight/go-vcsurl).

Installation
------------

The recommended way to install go-vcsurl

```
go get github.com/gitsight/go-vcsurl
```

Usage
-----

```go
urls := []string{
	"github.com/alice/libfoo",
	"git://github.com/bob/libbar",
	"https://gitlab.com/foo/bar",
	"https://github.com/go-enry/go-enry/releases/tag/v2.4.1",
}

for i, url := range urls {
	info, err := vcsurl.Parse(url)
	if err != nil {
		fmt.Printf("error parsing %s\n", err)
	}

	fmt.Printf("%d. %s %s\n", i+1, info.Kind, info.ID)
	fmt.Printf("   name: %s\n", info.Name)
	fmt.Printf("   host: %s\n", info.Host)

	remote, _ := info.Remote(vcsurl.SSH)
	fmt.Printf("   remote: %s\n", remote)

	if info.Committish != "" {
		fmt.Printf("   commit-ish: %s\n", info.Committish)
	}
}
```


```
1. git github.com/alice/libfoo
   name: libfoo
   host: github.com
   remote: git@github.com/alice/libfoo.git
2. git github.com/bob/libbar
   name: libbar
   host: github.com
   remote: git@github.com/bob/libbar.git
3. git gitlab.com/foo/bar
   name: bar
   host: gitlab.com
   remote: git@gitlab.com/foo/bar.git
4. git github.com/go-enry/go-enry
   name: go-enry
   host: github.com
   remote: git@github.com/go-enry/go-enry.git
   commit-ish: v2.4.1
```



License
-------

MIT, see [LICENSE](LICENSE)