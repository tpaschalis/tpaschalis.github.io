# patreon-epub — implementation plan

## Goal

Build a Go CLI tool that downloads posts from a Patreon creator's campaign via
the Patreon API v2, embeds any images referenced in the post content, and
outputs EPUB files suitable for offline reading on an e-reader (Kobo Aura).

---

## Design decisions

### No third-party dependencies
The EPUB format is a ZIP file with a specific internal layout. Rather than
pulling in an external epub library (which requires internet access during
development), we implement EPUB generation directly using Go's `archive/zip`
and `text/template` packages.

### EPUB version
We target **EPUB 3** with a mandatory EPUB 2 NCX table of contents for
backwards compatibility. Kobo devices handle both well; the NCX fallback
ensures the TOC works even on older firmware.

### MIME type for the `mimetype` entry
Per the EPUB spec, the `mimetype` file must be the **first** entry in the ZIP
and must use the `Store` (uncompressed) method. The builder enforces this by
calling `zip.CreateHeader` with `Method: zip.Store` before any other entries.

### Image handling
1. `<img src="…">` tags in the post HTML are extracted with a regexp.
2. Each image is downloaded concurrently-safe via a shared `imageCache`.
3. Images are stored in `OEBPS/images/` inside the EPUB ZIP.
4. The `src` attribute is rewritten to the relative path `../images/<hash>.<ext>`.
5. Post cover/thumbnail images (the `image` field in the API response) are
   prepended as a cover image in each chapter.
6. MIME type is detected from the HTTP `Content-Type` response header, with
   a fallback to the URL path extension.

### HTML → XHTML conversion
Patreon post content is raw HTML. EPUB requires well-formed XHTML inside
chapters. We apply a minimal transformation:
- Void elements (`<br>`, `<hr>`, `<input>`, `<link>`, `<meta>`) are
  self-closed with a regexp replacement.
- Content is wrapped in a full XHTML document template.
Deliberately avoiding a full HTML parser keeps the tool dependency-free.

### Post grouping
The `--group-by` flag (values: `all`, `year`, `month`) controls how posts are
split across output files. This is useful for long-running creators where a
single EPUB of hundreds of posts would be unwieldy on a device.

---

## File layout

```
patreon-epub/
├── main.go              CLI: flags, campaign resolution, grouping, orchestration
├── patreon/
│   ├── models.go        API response types (Post, Campaign, pagination, …)
│   └── client.go        HTTP client: GetCampaigns, GetAllPosts (with cursor pagination)
├── epubbuilder/
│   └── epub.go          EPUB writer: ZIP layout, image cache, templates
├── plan.md              This file
├── README.md            User-facing documentation
└── go.mod
```

---

## Patreon API notes

- **Base URL**: `https://www.patreon.com/api/oauth2/v2`
- **Auth**: `Authorization: Bearer <creator-access-token>` header
- **Token source**: Patreon developer portal →
  https://www.patreon.com/portal/registration/register-clients
- **Posts endpoint**: `GET /posts?filter[campaign_id]=…&fields[post]=title,content,published_at,url,image&sort=-published_at`
- **Pagination**: cursor-based via `meta.pagination.cursors.next`; the client
  follows all pages automatically.

---

## Potential future improvements

- Concurrent image downloads (currently sequential per article; the cache is
  already concurrency-safe).
- HTML sanitisation / proper HTML→XHTML parser for more robust conversion.
- Cover page for the whole EPUB (separate from per-article covers).
- Tag/label filtering (e.g. only posts with a specific Patreon tag).
- Progress bar for large campaigns.
- JSON/CSV index of downloaded posts to enable incremental updates.
