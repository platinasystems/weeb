package weeb

import (
	"path"
	"strings"

	"github.com/platinasystems/weeb/html"
)

type Site struct {
	DocByPath  map[string]*html.Doc
	PageByPath map[string]Page
}

type Page interface {
	PageBody(path string, d *html.Doc) html.BodyVec
}

// Does path match pattern?
func pathMatch(pattern, path string) bool {
	n := len(pattern)
	if n > 1 && pattern[n-1] == '/' {
		return len(path) >= n && path[0:n] == pattern
	} else {
		return pattern == path
	}
}

// Return the canonical path for p, eliminating . and .. elements.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)

	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}

// Find a handler on a handler map given a path string
// Most-specific (longest) pattern wins
func (s *Site) Match(key string) (p Page, d *html.Doc, pattern string) {
	n := 0
	key = cleanPath(key)
	for path, page := range s.PageByPath {
		if !pathMatch(path, key) {
			continue
		}
		if p == nil || len(path) > n {
			n = len(path)
			p = page
			pattern = path
			d = s.DocByPath[path]
		}
	}
	return
}

type Content struct {
	// URL path used to index this content (url.Path)
	URLPath string

	// Either FilePath is set or Data is provided inline.
	FilePath string
	Data     []byte

	UnixTimeLastModified int64
	ContentType          string
	ContentEncoding      string
}

var ContentByPath = make(map[string]*Content)

func (c *Content) Register() {
	ContentByPath[c.URLPath] = c
}

func GoPackageNameForPath(p string) (n string) {
	n = path.Base(p)
	n = strings.Replace(n, ".", "_", -1)
	return
}
