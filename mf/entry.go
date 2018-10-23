package mf

import (
	"net/url"
)

type EntryProps struct {
	Title       string    `json:"name"`
	Summary     string    `json:"summary"`
	Content     string    `json:"content"` // handle HTML
	Author      string    `json:"author"`  // handle embedded HCard
	URL         string    `json:"url"`
	Categories  []string  `json:"category"`
	SyndicateTo []url.URL `json:"mp-syndicate-to"`

	LikeOf    url.URL `json:"like-of"`
	InReplyTo url.URL `json:"in-reply-to"`
	RepostOf  url.URL `json:"repost-of"`
}
