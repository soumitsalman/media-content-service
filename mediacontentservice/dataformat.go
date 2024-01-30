package mediacontentservice

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

const (
	CHANNEL = "channel"
	POST    = "post"
	COMMENT = "comment"
)

type TableItem interface {
	GetKeys() (string, string)
}

type MediaContentItem struct {
	// unique identifier across media source. every reddit item has one. In reddit this is the name
	// in azure cosmos DB every record/item has to have an id.
	// In case of media content the media content itself comes with an unique identifier that we can use
	aztables.Entity

	// which social media source is this coming from
	Source string `json:"source"`
	// unique id across Source
	Id string `json:"id"`

	// represents text title of the item. Applies to subreddits and posts but not comments
	Title string `json:"title,omitempty"`
	// unique short name across the Source
	Name string `json:"name,omitempty"`
	// Subreddit, Post or Comment. This is not directly serialized
	Kind string `json:"kind"`

	// Applies to comments and posts.
	// For comments: this represents which post or comment does this comment respond to.
	// for posts: this is the same value as the channel
	ChannelName string `json:"channel,omitempty"`

	//post text
	Text string `json:"text"`
	// for posts this is url posted by the post
	// for subreddit this is link
	Url string `json:"url,omitempty"`

	//subreddit category
	Category string `json:"category,omitempty"`

	// author of posts or comments. Empty for subreddits
	Author string `json:"author,omitempty"`
	// date of creation of the post or comment. Empty for subreddits
	Created float64 `json:"created,omitempty"`

	// Applies to posts and comments. Doesn't apply to subreddits
	Score int `json:"score,omitempty"`
	// Number of comments to a post or a comment. Doesn't apply to subreddit
	Comments int `json:"comments,omitempty"`
	// Number of subscribers to a channel (subreddit). Doesn't apply to posts or comments
	Subscribers int `json:"subscribers,omitempty"`
	// number of likes, claps, thumbs-up
	ThumbsupCount int `json:"likes,omitempty"`
	// Applies to subreddit posts and comments. Doesn't apply to subreddits
	ThumbsupRatio float64 `json:"likes_ratio,omitempty"`

	Children []MediaContentItem `json:"children,omitempty"`
}

func (item *MediaContentItem) GetKeys() (string, string) {
	item.Entity.PartitionKey, item.Entity.RowKey = item.Source, item.Id
	return item.Entity.PartitionKey, item.Entity.RowKey
}

type UserEngagementItem struct {
	aztables.Entity
	Username  string `json:"username"`
	Source    string `json:"source"`
	ContentId string `json:"cid"`
	Action    string `json:"action"`
}

func (item *UserEngagementItem) GetKeys() (string, string) {
	item.Entity.PartitionKey, item.Entity.RowKey = item.Username, fmt.Sprintf("%s@%s:%s", item.ContentId, item.Source, item.Action)
	return item.Entity.PartitionKey, item.Entity.RowKey
}

type UserActionData struct {
	// in cosmos DB every item has to have an id. Here the id will be synthetic
	// other than azure cosmos DB literally no one cares about this field
	RecordId      string `json:"id"`
	ContentId     string `json:"content_id"`
	Source        string `json:"source"`
	UserId        string `json:"user_id"`
	Processed     bool   `json:"processed,omitempty"`
	Action        string `json:"action,omitempty"`
	ActionContent string `json:"content,omitempty"`
}

type UserMetadata struct {
	UserId    string   `json:"user_id"`
	Interests []string `json:"interests"`
}

type UserItem struct {
}
