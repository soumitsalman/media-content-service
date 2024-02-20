package api

type NewEngagementRequest struct {
	Username   string `json:"username,omitempty" bson:"username,omitempty"`
	UserSource string `json:"usersource,omitempty" bson:"usersource,omitempty"`
	Source     string `json:"source,omitempty" bson:"source,omitempty"`
	ContentId  string `json:"cid,omitempty" bson:"cid,omitempty"`
	Action     string `json:"action,omitempty" bson:"action,omitempty"`
}

type NewInterestRequest struct {
	Username   string `json:"username,omitempty" bson:"username,omitempty"`
	UserSource string `json:"usersource,omitempty" bson:"usersource,omitempty"`
	Category   string `json:"category,omitempty" bson:"category,omitempty"`
	Source     string `json:"source,omitempty" bson:"source,omitempty"`
	ContentId  string `json:"cid,omitempty" bson:"cid,omitempty"`
	Direction  string `json:"direction,omitempty" bson:"direction,omitempty"` // determining if this a positive interest or an explicit disinterest
}
