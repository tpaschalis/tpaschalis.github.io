package patreon

import "time"

// PostsResponse is the top-level response from the /posts endpoint.
type PostsResponse struct {
	Data []Post `json:"data"`
	Meta Meta   `json:"meta"`
}

// Post represents a single Patreon post.
type Post struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes PostAttributes `json:"attributes"`
}

// PostAttributes holds the content fields of a post.
type PostAttributes struct {
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	PublishedAt time.Time  `json:"published_at"`
	URL         string     `json:"url"`
	Image       *PostImage `json:"image"`
}

// PostImage holds the URLs for a post's thumbnail/cover image.
type PostImage struct {
	URL      string `json:"url"`
	LargeURL string `json:"large_url"`
	ThumbURL string `json:"thumb_url"`
}

// Meta holds pagination metadata.
type Meta struct {
	Pagination Pagination `json:"pagination"`
}

// Pagination describes the current page and the next cursor.
type Pagination struct {
	Cursors *Cursors `json:"cursors"`
	Total   int      `json:"total"`
}

// Cursors holds the opaque cursor for the next page.
type Cursors struct {
	Next string `json:"next"`
}

// CampaignsResponse is the top-level response from the /campaigns endpoint.
type CampaignsResponse struct {
	Data []Campaign `json:"data"`
}

// Campaign represents a Patreon campaign (creator page).
type Campaign struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes CampaignAttributes `json:"attributes"`
}

// CampaignAttributes holds metadata about a campaign.
type CampaignAttributes struct {
	CreationName string `json:"creation_name"`
	Summary      string `json:"summary"`
}
