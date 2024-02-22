package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func newContentsHandler(ctx *gin.Context) {
	// request for posting new contents
	var contents []MediaContentItem
	ctx.BindJSON(&contents)

	// TODO: remove this
	// utils.PrintTable[MediaContentItem](utils.SafeSlice[MediaContentItem](contents, 0, 5),
	// 	[]string{"Kind", "Channel", "Id", "Created", "Subscribers", "Comments", "Likes"},
	// 	func(item *MediaContentItem) []string {
	// 		return []string{item.Kind, item.ChannelName, item.Id, utils.DateToString(item.Created), strconv.Itoa(item.Subscribers), strconv.Itoa(item.Comments), strconv.Itoa(item.ThumbsupCount)}
	// 	})
	// log.Println("New Contents", len(contents))

	go NewMediaContents(contents)
	ctx.JSON(http.StatusOK, gin.H{"message": "process initiated"})
}

func getContentsHandler(ctx *gin.Context) {
	uid, kind, valid_user := ctx.Query("uid"), ctx.Query("kind"), false

	// if uid is not specified, use username and usersource to get the uid
	if uid == "" {
		username, usersource := ctx.Query("username"), ctx.Query("usersource")
		uid, valid_user = getGlobalUID(usersource, username)
	} else {
		valid_user = isValidUID(uid)
	}

	if valid_user {
		contents := GetUserContentSuggestions(uid, kind)
		log.Println("Got Contents", len(contents))
		ctx.JSON(http.StatusOK, contents)
	} else {
		log.Println("User does not exist")
		ctx.JSON(http.StatusNoContent, gin.H{"message": "User not found"})
	}
}

func getCredsHandler(ctx *gin.Context) {
	source := ctx.Param("source")
	creds := GetAllUserCredentials(source)

	// TODO: remove this
	// utils.PrintTable[UserCredentialItem](creds,
	// 	[]string{"UID", "Source", "Username", "Password"},
	// 	func(item *UserCredentialItem) []string {
	// 		return []string{item.UID, item.Source, item.Username, item.Password}
	// 	})
	// log.Println("Got Creds", len(creds))

	ctx.JSON(http.StatusOK, creds)
}

func newCredsHandler(ctx *gin.Context) {
	var cred UserCredentialItem
	ctx.BindJSON(&cred)
	// log.Println("New Creds")
	go NewUserCredential(cred)
	ctx.JSON(http.StatusOK, gin.H{"message": "process initiated"})
}

func newEngagementHandler(ctx *gin.Context) {
	var engagements []UserEngagementItem
	ctx.BindJSON(&engagements)

	// TODO: remove this
	// utils.PrintTable[UserEngagementItem](utils.SafeSlice[UserEngagementItem](engagements, 0, 5),
	// 	[]string{"Engagement"},
	// 	func(item *UserEngagementItem) []string {
	// 		return []string{fmt.Sprintf("%s->%s@%s:%s", item.Username, item.ContentId, item.Source, item.Action)}
	// 	})
	// log.Println("New Engagements", len(engagements))

	go NewEnagements(engagements)
	ctx.JSON(http.StatusOK, gin.H{"message": "process initiated"})
}

// func newInterestsHandler(ctx *gin.Context) {
// 	var interests []NewInterestRequest
// 	ctx.BindJSON(&interests)

// 	// TODO: remove this
// 	// utils.PrintTable[UserEngagementItem](utils.SafeSlice[UserEngagementItem](engagements, 0, 5),
// 	// 	[]string{"Engagement"},
// 	// 	func(item *UserEngagementItem) []string {
// 	// 		return []string{fmt.Sprintf("%s->%s@%s:%s", item.Username, item.ContentId, item.Source, item.Action)}
// 	// 	})
// 	log.Println("New Interests", len(interests))

// 	utils.Transform[NewInterestRequest, UserInterestItem]

// 	go NewInterests(interests)
// 	ctx.JSON(http.StatusOK, gin.H{"message": "process initiated"})
// }

func authenticationHandler(ctx *gin.Context) {
	// fmt.Println(ctx.GetHeader("X-API-Key"), getInternalAuthToken())
	if ctx.GetHeader("X-API-Key") == getInternalAuthToken() {
		ctx.Next()
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}

func createRateLimitHandler(r rate.Limit, b int) gin.HandlerFunc {
	rate_limiter := rate.NewLimiter(r, b)
	return func(ctx *gin.Context) {
		if rate_limiter.Allow() {
			ctx.Next()
		} else {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
		}
	}
}

func NewServer(r rate.Limit, b int) *gin.Engine {
	router := gin.Default()

	auth_group := router.Group("/")

	// authn and ratelimit middleware
	auth_group.Use(createRateLimitHandler(r, b), authenticationHandler)

	// routes
	auth_group.GET("/contents", getContentsHandler)
	auth_group.POST("/contents", newContentsHandler)
	auth_group.POST("/engagements", newEngagementHandler)
	// auth_group.POST("/interests", newInterestsHandler)
	auth_group.GET("/users/:source", getCredsHandler)
	auth_group.POST("/users", newCredsHandler)

	// run
	return router
}
