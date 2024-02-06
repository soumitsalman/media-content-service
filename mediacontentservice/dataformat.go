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

type DataItem interface {
	GetGlobalId() string
}

type MediaContentItem struct {
	// unique identifier across media source. every reddit item has one. In reddit this is the name
	// in azure cosmos DB every record/item has to have an id.
	// In case of media content the media content itself comes with an unique identifier that we can use
	aztables.Entity

	GlobalId string `json:"_id,omitempty" bson:"_id, omitempty"`
	// which social media source is this coming from
	Source string `json:"source,omitempty" bson:"source,omitempty"`
	// unique id across Source
	Id string `json:"cid,omitempty" bson:"cid,omitempty"`

	// represents text title of the item. Applies to subreddits and posts but not comments
	Title string `json:"title,omitempty" bson:"title,omitempty"`
	// unique short name across the Source
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	// Subreddit, Post or Comment. This is not directly serialized
	Kind string `json:"kind,omitempty" bson:"kind,omitempty"`

	// Applies to comments and posts.
	// For comments: this represents which post or comment does this comment respond to.
	// for posts: this is the same value as the channel
	ChannelName string `json:"channel,omitempty" bson:"channel,omitempty"`

	//post text
	Text string `json:"text,omitempty" bson:"text,omitempty"`
	// for posts this is url posted by the post
	// for subreddit this is link
	Url string `json:"url,omitempty" bson:"url,omitempty"`

	//subreddit category
	Category string `json:"category,omitempty" bson:"category,omitempty"`

	// author of posts or comments. Empty for subreddits
	Author string `json:"author,omitempty" bson:"author,omitempty"`
	// date of creation of the post or comment. Empty for subreddits
	Created float64 `json:"created,omitempty" bson:"created,omitempty"`

	// Applies to posts and comments. Doesn't apply to subreddits
	Score int `json:"score,omitempty" bson:"score,omitempty"`
	// Number of comments to a post or a comment. Doesn't apply to subreddit
	Comments int `json:"comments,omitempty" bson:"comments,omitempty"`
	// Number of subscribers to a channel (subreddit). Doesn't apply to posts or comments
	Subscribers int `json:"subscribers,omitempty" bson:"subscribers,omitempty"`
	// number of likes, claps, thumbs-up
	ThumbsupCount int `json:"likes,omitempty" bson:"likes,omitempty"`
	// Applies to subreddit posts and comments. Doesn't apply to subreddits
	ThumbsupRatio float64 `json:"likes_ratio,omitempty" bson:"likes_ratio,omitempty"`

	Digest     string    `json:"digest,omitempty" bson:"digest,omitempty"`
	Embeddings []float32 `json:"embeddings,omitempty" bson:"embeddings,omitempty"`
}

func (item *MediaContentItem) CreateKeys() (string, string) {
	item.Entity.PartitionKey, item.Entity.RowKey = item.Source, item.Id
	return item.Entity.PartitionKey, item.Entity.RowKey
}

func (item *MediaContentItem) CreateGlobalId() string {
	item.GlobalId = fmt.Sprintf("%s@%s", item.Id, item.Source)
	return item.GlobalId
}

func (item MediaContentItem) GetGlobalId() string {
	return item.GlobalId
}

func compareMediaContentItems(a, b *MediaContentItem) bool {
	if a.GlobalId == "" {
		a.CreateGlobalId()
	}
	if b.GlobalId == "" {
		a.CreateGlobalId()
	}
	return a.GlobalId == b.GlobalId
}

type UserEngagementItem struct {
	aztables.Entity
	GlobalId  string `json:"_id,omitempty" bson:"_id, omitempty"`
	Username  string `json:"username,omitempty"`
	Source    string `json:"source,omitempty"`
	ContentId string `json:"cid,omitempty"`
	Action    string `json:"action,omitempty"`
}

func (item *UserEngagementItem) CreateKeys() (string, string) {
	item.Entity.PartitionKey, item.Entity.RowKey = item.Username, fmt.Sprintf("%s@%s:%s", item.ContentId, item.Source, item.Action)
	return item.Entity.PartitionKey, item.Entity.RowKey
}

func (item *UserEngagementItem) CreateGlobalId() string {
	item.GlobalId = fmt.Sprintf("%s->%s@%s:%s", item.Username, item.ContentId, item.Source, item.Action)
	return item.GlobalId
}

func (item UserEngagementItem) GetGlobalId() string {
	return item.GlobalId
}

// type UserActionData struct {
// 	// in cosmos DB every item has to have an id. Here the id will be synthetic
// 	// other than azure cosmos DB literally no one cares about this field
// 	RecordId      string `json:"id"`
// 	ContentId     string `json:"content_id"`
// 	Source        string `json:"source"`
// 	UserId        string `json:"user_id"`
// 	Processed     bool   `json:"processed,omitempty"`
// 	Action        string `json:"action,omitempty"`
// 	ActionContent string `json:"content,omitempty"`
// }

type UserInterestItem struct {
	UserId     string    `json:"uid,omitempty" bson:"uid,omitempty"`
	Category   string    `json:"category,omitempty" bson:"category,omitempty"`
	Embeddings []float32 `json:"embeddings,omitempty" bson:"embeddings,omitempty"`
	Direction  string    `json:"direction,omitempty" bson:"direction,omitempty"` // determining if this a positive interest or an explicit disinterest
	Timestamp  float64   `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

type CategoryItem struct {
	Category   string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Embeddings []float32 `json:"embeddings,omitempty" bson:"embeddings,omitempty"`
}
