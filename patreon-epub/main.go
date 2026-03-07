// patreon-epub fetches posts from a Patreon campaign and writes them as EPUB
// files suitable for e-readers such as the Kobo Aura.
//
// Usage:
//
//	patreon-epub --token <creator-access-token> [flags]
//
// The creator access token can also be supplied via the PATREON_ACCESS_TOKEN
// environment variable.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/tpaschalis/patreon-epub/epubbuilder"
	"github.com/tpaschalis/patreon-epub/patreon"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	fs := flag.NewFlagSet("patreon-epub", flag.ExitOnError)

	token := fs.String("token", "", "Patreon creator access token (or set PATREON_ACCESS_TOKEN)")
	campaignID := fs.String("campaign-id", "", "Campaign ID (auto-detected when omitted)")
	outputDir := fs.String("output", ".", "Directory to write EPUB files into")
	groupBy := fs.String("group-by", "all", "How to split posts into EPUBs: all | year | month")
	limit := fs.Int("limit", 0, "Maximum number of posts to fetch (0 = unlimited)")
	sinceStr := fs.String("since", "", "Only include posts published on or after this date (YYYY-MM-DD)")
	author := fs.String("author", "", "Author name recorded in the EPUB metadata")
	title := fs.String("title", "", "Base title for the EPUB(s); defaults to campaign name")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `patreon-epub – download Patreon posts as EPUB files

Usage:
  patreon-epub [flags]

Flags:
`)
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Getting your creator access token:
  1. Go to https://www.patreon.com/portal/registration/register-clients
  2. Create a client or open an existing one.
  3. Copy the "Creator's Access Token" shown on that page.

Examples:
  # All posts into a single EPUB:
  patreon-epub --token YOUR_TOKEN --output ~/books

  # One EPUB per calendar month, from 2024 onwards:
  patreon-epub --token YOUR_TOKEN --group-by month --since 2024-01-01

  # Fetch a specific campaign (useful if you manage multiple):
  patreon-epub --token YOUR_TOKEN --campaign-id 12345678
`)
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	// Resolve token from env if not provided by flag.
	if *token == "" {
		*token = os.Getenv("PATREON_ACCESS_TOKEN")
	}
	if *token == "" {
		fs.Usage()
		return fmt.Errorf("--token is required (or set PATREON_ACCESS_TOKEN)")
	}

	// Validate --group-by.
	switch *groupBy {
	case "all", "year", "month":
	default:
		return fmt.Errorf("--group-by must be one of: all, year, month")
	}

	// Parse --since.
	var since time.Time
	if *sinceStr != "" {
		var err error
		since, err = time.Parse("2006-01-02", *sinceStr)
		if err != nil {
			return fmt.Errorf("--since must be in YYYY-MM-DD format: %w", err)
		}
	}

	client := patreon.NewClient(*token)

	// Resolve campaign ID.
	if *campaignID == "" {
		fmt.Fprintln(os.Stderr, "Fetching campaigns…")
		campaigns, err := client.GetCampaigns()
		if err != nil {
			return fmt.Errorf("fetching campaigns: %w", err)
		}
		if len(campaigns) == 0 {
			return fmt.Errorf("no campaigns found for this token")
		}
		if len(campaigns) > 1 {
			fmt.Fprintln(os.Stderr, "Multiple campaigns found; using the first one.")
			fmt.Fprintln(os.Stderr, "Pass --campaign-id to select a specific one:")
			for _, c := range campaigns {
				fmt.Fprintf(os.Stderr, "  %s  %s\n", c.ID, c.Attributes.CreationName)
			}
		}
		*campaignID = campaigns[0].ID
		if *title == "" {
			*title = campaigns[0].Attributes.CreationName
		}
	}
	if *title == "" {
		*title = "Patreon Posts"
	}

	fmt.Fprintf(os.Stderr, "Fetching posts for campaign %s…\n", *campaignID)
	posts, err := client.GetAllPosts(*campaignID, since, *limit)
	if err != nil {
		return fmt.Errorf("fetching posts: %w", err)
	}
	if len(posts) == 0 {
		fmt.Fprintln(os.Stderr, "No posts found.")
		return nil
	}
	fmt.Fprintf(os.Stderr, "Found %d posts.\n", len(posts))

	// Sort oldest-first so articles appear in reading order inside each EPUB.
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Attributes.PublishedAt.Before(posts[j].Attributes.PublishedAt)
	})

	// Group posts by the requested period.
	groups := groupPosts(posts, *groupBy)

	if err := os.MkdirAll(*outputDir, 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	for _, g := range groups {
		articles := make([]epubbuilder.Article, 0, len(g.posts))
		for _, p := range g.posts {
			a := p.Attributes
			coverURL := ""
			if a.Image != nil {
				if a.Image.LargeURL != "" {
					coverURL = a.Image.LargeURL
				} else {
					coverURL = a.Image.URL
				}
			}
			articles = append(articles, epubbuilder.Article{
				Title:       a.Title,
				PublishedAt: a.PublishedAt,
				Content:     a.Content,
				CoverURL:    coverURL,
			})
		}

		book := &epubbuilder.Book{
			Title:    g.title(*title),
			Author:   *author,
			Language: "en",
			Articles: articles,
		}

		filename := filepath.Join(*outputDir, g.filename(*title)+".epub")
		fmt.Fprintf(os.Stderr, "Writing %s (%d articles)…\n", filename, len(articles))
		if err := book.Write(filename); err != nil {
			return fmt.Errorf("writing %s: %w", filename, err)
		}
	}

	fmt.Fprintln(os.Stderr, "Done.")
	return nil
}

// -----------------------------------------------------------------------
// Post grouping
// -----------------------------------------------------------------------

type group struct {
	key   string // e.g. "2024", "2024-03"
	posts []patreon.Post
}

// title returns the EPUB title for this group.
func (g *group) title(base string) string {
	if g.key == "" {
		return base
	}
	return base + " – " + g.key
}

// filename returns a safe filename stem for this group's EPUB.
func (g *group) filename(base string) string {
	safe := strings.NewReplacer(" ", "_", "/", "-", ":", "-").Replace(base)
	if g.key == "" {
		return safe
	}
	return safe + "_" + g.key
}

func groupPosts(posts []patreon.Post, by string) []group {
	if by == "all" {
		return []group{{posts: posts}}
	}

	index := map[string]int{}
	var groups []group

	for _, p := range posts {
		var key string
		switch by {
		case "year":
			key = p.Attributes.PublishedAt.Format("2006")
		case "month":
			key = p.Attributes.PublishedAt.Format("2006-01")
		}
		if _, ok := index[key]; !ok {
			index[key] = len(groups)
			groups = append(groups, group{key: key})
		}
		i := index[key]
		groups[i].posts = append(groups[i].posts, p)
	}

	return groups
}
