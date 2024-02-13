package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	utils "github.com/soumitsalman/data-utils"
	ds "github.com/soumitsalman/media-content-service/api"
	"golang.org/x/time/rate"
)

func getPort() string {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		return ":" + port
	}
	return ":8080"
}

func getInternalAuthToken() string {
	return os.Getenv("INTERNAL_AUTH_TOKEN")
}

func newContentsHandler(ctx *gin.Context) {
	// request for posting new contents
	var contents []ds.MediaContentItem
	ctx.BindJSON(&contents)

	// TODO: remove this
	utils.PrintTable[ds.MediaContentItem](utils.SafeSlice[ds.MediaContentItem](contents, 0, 5),
		[]string{"Kind", "Channel", "Id", "Created", "Subscribers", "Comments", "Likes"},
		func(item *ds.MediaContentItem) []string {
			return []string{item.Kind, item.ChannelName, item.Id, utils.DateToString(item.Created), strconv.Itoa(item.Subscribers), strconv.Itoa(item.Comments), strconv.Itoa(item.ThumbsupCount)}
		})

	go ds.NewMediaContents_Mongo(contents)
	ctx.JSON(http.StatusOK, gin.H{"message": "process initiated"})
}

func getContentsHandler(ctx *gin.Context) {
	uid, kind := ctx.Param("uid"), ctx.Query("kind")
	log.Println(uid, kind)
	contents := ds.GetUserContentSuggestions(uid, kind)

	// TODO: remove this
	utils.PrintTable[ds.MediaContentItem](utils.SafeSlice[ds.MediaContentItem](contents, 0, 5),
		[]string{"Kind", "Channel", "Id", "Created", "Subscribers", "Comments", "Likes"},
		func(item *ds.MediaContentItem) []string {
			return []string{item.Kind, item.ChannelName, item.Id, utils.DateToString(item.Created), strconv.Itoa(item.Subscribers), strconv.Itoa(item.Comments), strconv.Itoa(item.ThumbsupCount)}
		})

	ctx.JSON(http.StatusOK, contents)
}

func getCredsHandler(ctx *gin.Context) {
	source := ctx.Param("source")
	creds := ds.GetAllUserCredentials(source)

	// TODO: remove this
	utils.PrintTable[ds.UserCredentialItem](creds,
		[]string{"UID", "Source", "Username", "Password"},
		func(item *ds.UserCredentialItem) []string {
			return []string{item.UID, item.Source, item.Username, item.Password}
		})

	ctx.JSON(http.StatusOK, creds)
}

func newCredsHandler(ctx *gin.Context) {
	var cred ds.UserCredentialItem
	ctx.BindJSON(&cred)
	go ds.NewUserCredential_Mongo(cred)
	ctx.JSON(http.StatusOK, gin.H{"message": "process initiated"})
}

func newEngagementHandler(ctx *gin.Context) {
	var engagements []ds.UserEngagementItem
	ctx.BindJSON(&engagements)

	// TODO: remove this
	utils.PrintTable[ds.UserEngagementItem](utils.SafeSlice[ds.UserEngagementItem](engagements, 0, 5),
		[]string{"Engagement"},
		func(item *ds.UserEngagementItem) []string {
			return []string{fmt.Sprintf("%s->%s@%s:%s", item.Username, item.ContentId, item.Source, item.Action)}
		})

	go ds.NewEnagements_Mongo(engagements)
	ctx.JSON(http.StatusOK, gin.H{"message": "process initiated"})
}

func main() {
	log.Println(os.Environ())

	router := gin.Default()

	auth_group := router.Group("/")

	// 100 tokens per sec with burst of 1000
	// TODO: make this check per user in future
	rate_limiter := rate.NewLimiter(100, 1000)
	rate_limit_handler := func(ctx *gin.Context) {
		if rate_limiter.Allow() {
			ctx.Next()
		} else {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
		}
	}

	auth_handler := func(ctx *gin.Context) {
		if ctx.GetHeader("Authorization") == getInternalAuthToken() {
			ctx.Next()
		} else {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}

	// authn and ratelimit middleware
	auth_group.Use(rate_limit_handler, auth_handler)

	// routes
	auth_group.GET("/contents/:uid", getContentsHandler)
	auth_group.POST("/contents/", newContentsHandler)
	// router.GET("/engagements/", contentsHandler)
	auth_group.POST("/engagements/", newEngagementHandler)
	auth_group.GET("/users/:source", getCredsHandler)
	auth_group.POST("/users/", newCredsHandler)

	// run
	router.Run(getPort())
}
