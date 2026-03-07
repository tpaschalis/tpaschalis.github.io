# patreon-epub

A command-line tool that downloads posts from a Patreon creator's campaign and
packages them as EPUB files — ready to sideload onto your Kobo, Kindle, or any
other e-reader.

Images referenced inside posts are downloaded and embedded in the EPUB so
articles read correctly offline.

---

## Prerequisites

- Go 1.21 or later
- A **Creator's Access Token** from the Patreon developer portal

### Getting a Creator's Access Token

1. Go to <https://www.patreon.com/portal/registration/register-clients>
2. Create a new client (or open an existing one)
3. Copy the **Creator's Access Token** shown on that page

> This token gives read access to your own campaigns and their posts.
> It does **not** grant access to other creators' patron-only content.

---

## Build

```bash
git clone <this repo>
cd patreon-epub
go build -o patreon-epub .
```

---

## Usage

```
patreon-epub [flags]

Flags:
  --token         string   Creator access token (or set PATREON_ACCESS_TOKEN)
  --campaign-id   string   Campaign ID — auto-detected when omitted
  --output        string   Output directory (default: current directory)
  --group-by      string   Split posts into EPUBs: all | year | month  (default: all)
  --limit         int      Max posts to fetch, 0 = unlimited
  --since         string   Only posts on/after this date YYYY-MM-DD
  --author        string   Author name embedded in EPUB metadata
  --title         string   Base title for the EPUB(s)
```

### Examples

**Download everything into a single EPUB:**

```bash
export PATREON_ACCESS_TOKEN=your_token_here
patreon-epub --output ~/books
```

**One EPUB per calendar month, starting from January 2024:**

```bash
patreon-epub \
  --token your_token_here \
  --group-by month \
  --since 2024-01-01 \
  --output ~/books/patreon
```

**Specify a campaign explicitly (useful if you run multiple campaigns):**

```bash
patreon-epub --token your_token_here --campaign-id 12345678
```

---

## Output

Each EPUB is named after the base title and the grouping key, e.g.:

| `--group-by` | Output filenames |
|---|---|
| `all` | `My_Blog.epub` |
| `year` | `My_Blog_2023.epub`, `My_Blog_2024.epub` |
| `month` | `My_Blog_2024-01.epub`, `My_Blog_2024-02.epub`, … |

---

## Sideloading onto a Kobo Aura

1. Connect the Kobo via USB — it mounts as a drive.
2. Copy the `.epub` files into the `Digital Editions` folder on the device
   (or the root if that folder doesn't exist).
3. Eject and unplug. The Kobo library will update automatically.

---

## How it works

1. **Fetch**: calls `GET /api/oauth2/v2/posts` with cursor-based pagination
   until all posts are retrieved.
2. **Process**: for each post, downloads any `<img src="…">` images and the
   post's thumbnail. Images are stored inside the EPUB at `OEBPS/images/`.
3. **Build**: assembles a valid EPUB 3 ZIP archive — `mimetype`, `container.xml`,
   `content.opf`, `toc.ncx`, `nav.xhtml`, per-article XHTML chapters, CSS,
   and embedded images.

---

## Limitations

- Patron-only posts are accessible only if the token belongs to the creator
  of that campaign.
- Post HTML is lightly sanitised for XHTML compliance (void elements are
  self-closed). Heavily malformed HTML may render oddly.
- Videos and audio embedded in posts are not downloaded; only `<img>` tags
  are processed.
