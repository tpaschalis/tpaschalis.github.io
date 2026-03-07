package patreon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://www.patreon.com/api/oauth2/v2"

// Client is an authenticated Patreon API v2 client.
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient creates a Patreon API client using a Creator Access Token.
// Obtain your token at https://www.patreon.com/portal/registration/register-clients.
func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetCampaigns returns the campaigns owned by the authenticated creator.
func (c *Client) GetCampaigns() ([]Campaign, error) {
	req, err := http.NewRequest("GET", baseURL+"/campaigns", nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching campaigns: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned %s for /campaigns", resp.Status)
	}

	var result CampaignsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding campaigns response: %w", err)
	}
	return result.Data, nil
}

// GetAllPosts fetches every post for a campaign, following pagination cursors.
// Posts are returned in reverse-chronological order (newest first) as returned
// by the API; callers may re-sort as needed.
func (c *Client) GetAllPosts(campaignID string, since time.Time, limit int) ([]Post, error) {
	var all []Post
	cursor := ""

	for {
		params := url.Values{}
		params.Set("filter[campaign_id]", campaignID)
		// Request the fields we care about.
		params.Set("fields[post]", "title,content,published_at,url,image")
		// Newest first so we can stop early when --since is used.
		params.Set("sort", "-published_at")
		if cursor != "" {
			params.Set("page[cursor]", cursor)
		}

		u := baseURL + "/posts?" + params.Encode()
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return nil, fmt.Errorf("building request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+c.token)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetching posts: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("API returned %s for /posts", resp.Status)
		}

		var result PostsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("decoding posts response: %w", err)
		}
		resp.Body.Close()

		for _, post := range result.Data {
			// Stop fetching older pages once we pass the --since threshold.
			if !since.IsZero() && post.Attributes.PublishedAt.Before(since) {
				return all, nil
			}
			all = append(all, post)
			if limit > 0 && len(all) >= limit {
				return all, nil
			}
		}

		// No more pages.
		if result.Meta.Pagination.Cursors == nil || result.Meta.Pagination.Cursors.Next == "" {
			break
		}
		cursor = result.Meta.Pagination.Cursors.Next
	}

	return all, nil
}
