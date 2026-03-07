// Package epubbuilder creates EPUB 3 files from a list of HTML articles.
//
// The generated EPUBs embed all referenced images (downloaded over HTTP) and
// include an NCX table of contents for compatibility with older e-readers such
// as the Kobo Aura.
package epubbuilder

import (
	"archive/zip"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"
)

// Article is a single piece of content to be added to the EPUB.
type Article struct {
	Title       string
	PublishedAt time.Time
	// Content is the raw HTML body of the article (not a full document).
	Content  string
	CoverURL string // optional thumbnail/cover image URL
}

// Book holds everything needed to render an EPUB file.
type Book struct {
	Title    string
	Author   string
	Language string
	Articles []Article
}

// -----------------------------------------------------------------------
// Entry point
// -----------------------------------------------------------------------

// Write renders the Book as an EPUB file at dest.
func (b *Book) Write(dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	// "mimetype" MUST be the first file and MUST be stored uncompressed.
	if err := addStored(zw, "mimetype", "application/epub+zip"); err != nil {
		return err
	}

	if err := addFile(zw, "META-INF/container.xml", containerXML); err != nil {
		return err
	}

	uid := newUUID()

	// Download and embed images, rewriting src= attributes in each article.
	imgCache := &imageCache{client: &http.Client{Timeout: 30 * time.Second}}
	chapters := make([]chapterMeta, 0, len(b.Articles))

	for i, art := range b.Articles {
		slug := fmt.Sprintf("article_%04d", i+1)
		html, imgs, err := processContent(art, imgCache)
		if err != nil {
			// Non-fatal: log and continue with original content.
			fmt.Fprintf(os.Stderr, "warning: processing %q: %v\n", art.Title, err)
			html = art.Content
			imgs = nil
		}

		ch := chapterMeta{
			ID:       slug,
			Filename: "articles/" + slug + ".xhtml",
			Title:    art.Title,
			Date:     art.PublishedAt,
			Images:   imgs,
		}
		chapters = append(chapters, ch)

		xhtml, err := renderArticle(art.Title, art.PublishedAt, html)
		if err != nil {
			return fmt.Errorf("rendering article %d: %w", i+1, err)
		}
		if err := addFile(zw, "OEBPS/"+ch.Filename, xhtml); err != nil {
			return err
		}
	}

	// Write embedded images collected from all articles.
	for path, data := range imgCache.files {
		w, err := zw.Create("OEBPS/images/" + path)
		if err != nil {
			return fmt.Errorf("adding image %s: %w", path, err)
		}
		if _, err := w.Write(data); err != nil {
			return fmt.Errorf("writing image %s: %w", path, err)
		}
	}

	if err := addFile(zw, "OEBPS/style.css", defaultCSS); err != nil {
		return err
	}

	// Render OPF and NCX using templates.
	opf, err := renderOPF(b, uid, chapters, imgCache)
	if err != nil {
		return fmt.Errorf("rendering OPF: %w", err)
	}
	if err := addFile(zw, "OEBPS/content.opf", opf); err != nil {
		return err
	}

	ncx, err := renderNCX(b, uid, chapters)
	if err != nil {
		return fmt.Errorf("rendering NCX: %w", err)
	}
	if err := addFile(zw, "OEBPS/toc.ncx", ncx); err != nil {
		return err
	}

	nav, err := renderNav(b, chapters)
	if err != nil {
		return fmt.Errorf("rendering nav: %w", err)
	}
	if err := addFile(zw, "OEBPS/nav.xhtml", nav); err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------
// Image downloading & caching
// -----------------------------------------------------------------------

type imageCache struct {
	mu     sync.Mutex
	client *http.Client
	// files maps a local filename (e.g. "abc123.jpg") to its raw bytes.
	files map[string][]byte
	// urls maps original URL to local filename.
	urls map[string]string
}

func (c *imageCache) fetch(rawURL string) (localName string, err error) {
	c.mu.Lock()
	if c.files == nil {
		c.files = make(map[string][]byte)
		c.urls = make(map[string]string)
	}
	if name, ok := c.urls[rawURL]; ok {
		c.mu.Unlock()
		return name, nil
	}
	c.mu.Unlock()

	resp, err := c.client.Get(rawURL) //nolint:noctx
	if err != nil {
		return "", fmt.Errorf("downloading image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("image server returned %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading image: %w", err)
	}

	ext := extFromContentType(resp.Header.Get("Content-Type"), rawURL)
	name := randHex(8) + ext

	c.mu.Lock()
	c.files[name] = data
	c.urls[rawURL] = name
	c.mu.Unlock()

	return name, nil
}

// extFromContentType returns an appropriate file extension for an image.
func extFromContentType(ct, rawURL string) string {
	ct = strings.ToLower(strings.Split(ct, ";")[0])
	switch ct {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	}
	// Fall back to URL extension.
	u, err := url.Parse(rawURL)
	if err == nil {
		ext := strings.ToLower(filepath.Ext(u.Path))
		if ext != "" {
			return ext
		}
	}
	return ".jpg"
}

// mimeFromExt returns the MIME type for a given file extension.
func mimeFromExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "image/jpeg"
	}
}

// -----------------------------------------------------------------------
// HTML processing
// -----------------------------------------------------------------------

// imgSrcRe matches <img ... src="URL" ...> in HTML.
var imgSrcRe = regexp.MustCompile(`(?i)<img([^>]*?)\bsrc="([^"]+)"([^>]*)>`)

// voidTagRe self-closes void elements that are invalid in XHTML.
var voidTagRe = regexp.MustCompile(`(?i)<(br|hr|input|link|meta)([^/]*[^/])?>`)

type imgInfo struct {
	localName string
	mimeType  string
}

// processContent downloads images referenced in article.Content and
// article.CoverURL, rewrites the src= attributes to point to the local copies,
// and returns modified HTML together with image metadata.
func processContent(art Article, cache *imageCache) (string, []imgInfo, error) {
	html := art.Content
	var imgs []imgInfo
	seen := make(map[string]bool)

	var fetchErr error
	html = imgSrcRe.ReplaceAllStringFunc(html, func(match string) string {
		sub := imgSrcRe.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		srcURL := sub[2]
		localName, err := cache.fetch(srcURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: skipping image %s: %v\n", srcURL, err)
			return match
		}
		if !seen[localName] {
			seen[localName] = true
			imgs = append(imgs, imgInfo{
				localName: localName,
				mimeType:  mimeFromExt(filepath.Ext(localName)),
			})
		}
		// Rewrite src to the relative path from articles/ directory.
		return fmt.Sprintf(`<img%s src="../images/%s"%s/>`, sub[1], localName, sub[3])
	})

	// Prepend cover image if present.
	if art.CoverURL != "" {
		localName, err := cache.fetch(art.CoverURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: skipping cover image: %v\n", err)
		} else {
			if !seen[localName] {
				seen[localName] = true
				imgs = append([]imgInfo{{
					localName: localName,
					mimeType:  mimeFromExt(filepath.Ext(localName)),
				}}, imgs...)
			}
			html = fmt.Sprintf(`<div class="cover-image"><img src="../images/%s" alt="cover"/></div>`, localName) + html
		}
	}

	if fetchErr != nil {
		return html, imgs, fetchErr
	}

	// Fix common XHTML issues: self-close void elements.
	html = voidTagRe.ReplaceAllString(html, "<$1$2/>")

	return html, imgs, nil
}

// -----------------------------------------------------------------------
// Template rendering
// -----------------------------------------------------------------------

type chapterMeta struct {
	ID       string
	Filename string // relative to OEBPS/, e.g. "articles/article_0001.xhtml"
	Title    string
	Date     time.Time
	Images   []imgInfo
}

func renderArticle(title string, date time.Time, body string) (string, error) {
	return renderTemplate(articleTpl, map[string]any{
		"Title": title,
		"Date":  date.Format("2 January 2006"),
		"Body":  body,
	})
}

func renderOPF(b *Book, uid string, chapters []chapterMeta, cache *imageCache) (string, error) {
	type imgManifestItem struct {
		ID       string
		Filename string
		MIME     string
	}
	var imgItems []imgManifestItem
	for name := range cache.files {
		imgItems = append(imgItems, imgManifestItem{
			ID:       "img-" + strings.TrimSuffix(name, filepath.Ext(name)),
			Filename: "images/" + name,
			MIME:     mimeFromExt(filepath.Ext(name)),
		})
	}

	return renderTemplate(opfTpl, map[string]any{
		"Title":    b.Title,
		"Author":   b.Author,
		"Language": b.Language,
		"UID":      uid,
		"Chapters": chapters,
		"Images":   imgItems,
	})
}

func renderNCX(b *Book, uid string, chapters []chapterMeta) (string, error) {
	return renderTemplate(ncxTpl, map[string]any{
		"Title":    b.Title,
		"UID":      uid,
		"Chapters": chapters,
	})
}

func renderNav(b *Book, chapters []chapterMeta) (string, error) {
	return renderTemplate(navTpl, map[string]any{
		"Title":    b.Title,
		"Chapters": chapters,
	})
}

func renderTemplate(tpl *template.Template, data any) (string, error) {
	var sb strings.Builder
	if err := tpl.Execute(&sb, data); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// -----------------------------------------------------------------------
// ZIP helpers
// -----------------------------------------------------------------------

// addFile adds a deflated file to the ZIP archive.
func addFile(zw *zip.Writer, name, content string) error {
	w, err := zw.Create(name)
	if err != nil {
		return fmt.Errorf("creating zip entry %s: %w", name, err)
	}
	_, err = io.WriteString(w, content)
	return err
}

// addStored adds an uncompressed (stored) file to the ZIP archive.
func addStored(zw *zip.Writer, name, content string) error {
	h := &zip.FileHeader{
		Name:   name,
		Method: zip.Store,
	}
	w, err := zw.CreateHeader(h)
	if err != nil {
		return fmt.Errorf("creating stored zip entry %s: %w", name, err)
	}
	_, err = io.WriteString(w, content)
	return err
}

// -----------------------------------------------------------------------
// Misc helpers
// -----------------------------------------------------------------------

func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// -----------------------------------------------------------------------
// Static content
// -----------------------------------------------------------------------

const containerXML = `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:schemas:container">
  <rootfiles>
    <rootfile full-path="OEBPS/content.opf"
              media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>`

const defaultCSS = `
body {
  font-family: Georgia, serif;
  font-size: 1em;
  line-height: 1.6;
  margin: 1em 1.5em;
  color: #1a1a1a;
}
h1 { font-size: 1.4em; margin-bottom: 0.2em; }
h2 { font-size: 1.2em; }
.meta { font-size: 0.85em; color: #555; margin-bottom: 1.2em; }
img { max-width: 100%; height: auto; display: block; margin: 1em auto; }
.cover-image { text-align: center; margin-bottom: 1.5em; }
blockquote {
  border-left: 3px solid #ccc;
  margin: 1em 0;
  padding: 0 1em;
  color: #555;
}
pre, code { font-family: monospace; font-size: 0.9em; }
pre { overflow-x: auto; background: #f5f5f5; padding: 0.8em; }
a { color: #333; }
`

// -----------------------------------------------------------------------
// Templates
// -----------------------------------------------------------------------

var articleTpl = template.Must(template.New("article").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en">
<head>
  <meta charset="UTF-8"/>
  <title>{{.Title}}</title>
  <link rel="stylesheet" type="text/css" href="../style.css"/>
</head>
<body>
  <h1>{{.Title}}</h1>
  <p class="meta">{{.Date}}</p>
  <div class="content">
{{.Body}}
  </div>
</body>
</html>`))

var opfTpl = template.Must(template.New("opf").Funcs(template.FuncMap{
	"now": func() string { return time.Now().UTC().Format(time.RFC3339) },
}).Parse(`<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf"
         version="3.0"
         unique-identifier="BookId">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/"
            xmlns:opf="http://www.idpf.org/2007/opf">
    <dc:title>{{.Title}}</dc:title>
    <dc:creator>{{.Author}}</dc:creator>
    <dc:language>{{.Language}}</dc:language>
    <dc:identifier id="BookId">urn:uuid:{{.UID}}</dc:identifier>
    <meta property="dcterms:modified">{{now}}</meta>
  </metadata>
  <manifest>
    <item id="ncx"  href="toc.ncx"   media-type="application/x-dtbncx+xml"/>
    <item id="nav"  href="nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>
    <item id="css"  href="style.css" media-type="text/css"/>
{{- range .Chapters}}
    <item id="{{.ID}}" href="{{.Filename}}" media-type="application/xhtml+xml"/>
{{- end}}
{{- range .Images}}
    <item id="{{.ID}}" href="{{.Filename}}" media-type="{{.MIME}}"/>
{{- end}}
  </manifest>
  <spine toc="ncx">
{{- range .Chapters}}
    <itemref idref="{{.ID}}"/>
{{- end}}
  </spine>
</package>`))

var ncxTpl = template.Must(template.New("ncx").Funcs(template.FuncMap{
	"inc": func(i int) int { return i + 1 },
}).Parse(`<?xml version="1.0" encoding="UTF-8"?>
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1">
  <head>
    <meta name="dtb:uid" content="urn:uuid:{{.UID}}"/>
    <meta name="dtb:depth" content="1"/>
    <meta name="dtb:totalPageCount" content="0"/>
    <meta name="dtb:maxPageNum" content="0"/>
  </head>
  <docTitle><text>{{.Title}}</text></docTitle>
  <navMap>
{{- range $i, $ch := .Chapters}}
    <navPoint id="navPoint-{{inc $i}}" playOrder="{{inc $i}}">
      <navLabel><text>{{$ch.Title}}</text></navLabel>
      <content src="{{$ch.Filename}}"/>
    </navPoint>
{{- end}}
  </navMap>
</ncx>`))

var navTpl = template.Must(template.New("nav").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml"
      xmlns:epub="http://www.idpf.org/2007/ops"
      xml:lang="en">
<head>
  <meta charset="UTF-8"/>
  <title>{{.Title}}</title>
</head>
<body>
  <nav epub:type="toc" id="toc">
    <h1>Table of Contents</h1>
    <ol>
{{- range .Chapters}}
      <li><a href="{{.Filename}}">{{.Title}}</a></li>
{{- end}}
    </ol>
  </nav>
</body>
</html>`))

