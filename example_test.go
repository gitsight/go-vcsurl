package vcsurl_test

import (
	"fmt"

	"github.com/gitsight/go-vcsurl"
)

func ExampleParse() {
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

	// output:
	// 1. git github.com/alice/libfoo
	//    name: libfoo
	//    host: github.com
	//    remote: git@github.com/alice/libfoo.git
	// 2. git github.com/bob/libbar
	//    name: libbar
	//    host: github.com
	//    remote: git@github.com/bob/libbar.git
	// 3. git gitlab.com/foo/bar
	//    name: bar
	//    host: gitlab.com
	//    remote: git@gitlab.com/foo/bar.git
	// 4. git github.com/go-enry/go-enry
	//    name: go-enry
	//    host: github.com
	//    remote: git@github.com/go-enry/go-enry.git
	//    commit-ish: v2.4.1
}
