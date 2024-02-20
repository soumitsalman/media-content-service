package api

const (
	CHANNEL = "channel"
	POST    = "post"
	COMMENT = "comment"
)

type MediaContentItem struct {
	Source string `json:"source,omitempty" bson:"source,omitempty"` // which social media source is this coming from
	Id     string `json:"cid,omitempty" bson:"cid,omitempty"`       // unique id across Source

	Title       string `json:"title,omitempty" bson:"title,omitempty"` // represents text title of the item. Applies to subreddits and posts but not comments
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Kind        string `json:"kind,omitempty" bson:"kind,omitempty"`
	ChannelName string `json:"channel,omitempty" bson:"channel,omitempty"` // fancy name of the channel represented by the channel itself or the channel where the post/comment is
	Excerpt     string `json:"excerpt,omitempty" bson:"excerpt,omitempty"`
	Text        string `json:"text,omitempty" bson:"text,omitempty"`
	Url         string `json:"url,omitempty" bson:"url,omitempty"`

	//subreddit category
	Category      string   `json:"category,omitempty" bson:"category,omitempty"`
	Tags          []string `json:"tags,omitempty" bson:"tags,omitempty"`
	Author        string   `json:"author,omitempty" bson:"author,omitempty"`           // author of posts or comments. Empty for subreddits
	Created       float64  `json:"created,omitempty" bson:"created,omitempty"`         // date of creation of the post or comment. Empty for subreddits
	Score         int      `json:"score,omitempty" bson:"score,omitempty"`             // Applies to posts and comments. Doesn't apply to subreddits
	Comments      int      `json:"comments,omitempty" bson:"comments,omitempty"`       // Number of comments to a post or a comment. Doesn't apply to subreddit
	Subscribers   int      `json:"subscribers,omitempty" bson:"subscribers,omitempty"` // Number of subscribers to a channel (subreddit). Doesn't apply to posts or comments
	ThumbsupCount int      `json:"likes,omitempty" bson:"likes,omitempty"`             // number of likes, claps, thumbs-up
	ThumbsupRatio float64  `json:"likes_ratio,omitempty" bson:"likes_ratio,omitempty"` // Applies to subreddit posts and comments. Doesn't apply to subreddits

	Digest     string    `json:"digest,omitempty" bson:"digest,omitempty"`
	Embeddings []float32 `json:"embeddings,omitempty" bson:"embeddings,omitempty"`
}

// TODO: remove this
// func (item *MediaContentItem) CreateKeys() (string, string) {
// 	item.Entity.PartitionKey, item.Entity.RowKey = item.Source, item.Id
// 	return item.Entity.PartitionKey, item.Entity.RowKey
// }

func compareMediaContents(a, b *MediaContentItem) bool {
	return (a.Source == b.Source) && (a.Id == b.Id)
}

type CategoryItem struct {
	Category   string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Embeddings []float32 `json:"embeddings,omitempty" bson:"embeddings,omitempty"`
}

type UserEngagementItem struct {
	UID        string `json:"uid,omitempty" bson:"uid,omitempty"`
	Username   string `json:"username,omitempty" bson:"username,omitempty"`
	UserSource string `json:"usersource,omitempty" bson:"usersource,omitempty"`
	Source     string `json:"source,omitempty" bson:"source,omitempty"`
	ContentId  string `json:"cid,omitempty" bson:"cid,omitempty"`
	Action     string `json:"action,omitempty" bson:"action,omitempty"`
}

// func (item *UserEngagementItem) CreateKeys() (string, string) {
// 	item.Entity.PartitionKey, item.Entity.RowKey = item.Username, fmt.Sprintf("%s@%s:%s", item.ContentId, item.Source, item.Action)
// 	return item.Entity.PartitionKey, item.Entity.RowKey
// }

func compareUserEngagements(a, b *UserEngagementItem) bool {
	return (a.Username == b.Username) &&
		(a.Source == b.Source) &&
		(a.ContentId == b.ContentId) &&
		(a.Action == b.Action)
}

type UserInterestItem struct {
	UID        string    `json:"uid,omitempty" bson:"uid,omitempty"`
	Category   string    `json:"category,omitempty" bson:"category,omitempty"`
	ContentId  string    `json:"cid,omitempty" bson:"cid,omitempty"`
	Embeddings []float32 `json:"embeddings,omitempty" bson:"embeddings,omitempty"`
	Direction  string    `json:"direction,omitempty" bson:"direction,omitempty"` // determining if this a positive interest or an explicit disinterest
	Timestamp  float64   `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

type UserCredentialItem struct {
	UID       string `json:"uid,omitempty" bson:"uid,omitempty"`
	Source    string `json:"source,omitempty" bson:"source,omitempty"`
	Username  string `json:"username,omitempty" bson:"username,omitempty"`
	Password  string `json:"password,omitempty" bson:"password,omitempty"`
	AuthToken string `json:"auth_token,omitempty" bson:"auth_token,omitempty"`
}
