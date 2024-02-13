package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	ds "github.com/soumitsr/media-content-service/mediacontentservice"
)

func newContentsHandler(ctx *gin.Context) {
	// request for posting new contents
	var contents []ds.MediaContentItem
	ctx.BindJSON(&contents)

	ds.PrintTable[ds.MediaContentItem](ds.SafeSlice[ds.MediaContentItem](contents, 0, 5),
		[]string{"Kind", "Channel", "Id", "Created", "Subscribers", "Comments", "Likes"},
		func(item *ds.MediaContentItem) []string {
			return []string{item.Kind, item.ChannelName, item.Id, ds.DateToString(item.Created), strconv.Itoa(item.Subscribers), strconv.Itoa(item.Comments), strconv.Itoa(item.ThumbsupCount)}
		})

	go ds.NewMediaContents_Mongo(contents)
	ctx.JSON(http.StatusOK, gin.H{"message": "process initiated"})
}

func getContentsHandler(ctx *gin.Context) {
	uid, kind := ctx.Param("uid"), ctx.Query("kind")
	log.Println(uid, kind)
	contents := ds.GetUserContentSuggestions(uid, kind)

	ds.PrintTable[ds.MediaContentItem](ds.SafeSlice[ds.MediaContentItem](contents, 0, 5),
		[]string{"Kind", "Channel", "Id", "Created", "Subscribers", "Comments", "Likes"},
		func(item *ds.MediaContentItem) []string {
			return []string{item.Kind, item.ChannelName, item.Id, ds.DateToString(item.Created), strconv.Itoa(item.Subscribers), strconv.Itoa(item.Comments), strconv.Itoa(item.ThumbsupCount)}
		})

	ctx.JSON(http.StatusOK, contents)
}

func getCredsHandler(ctx *gin.Context) {
	source := ctx.Param("source")
	creds := ds.GetAllUserCredentials(source)

	ds.PrintTable[ds.UserCredentialItem](creds,
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

	ds.PrintTable[ds.UserEngagementItem](ds.SafeSlice[ds.UserEngagementItem](engagements, 0, 5),
		[]string{"Engagement"},
		func(item *ds.UserEngagementItem) []string {
			return []string{fmt.Sprintf("%s->%s@%s:%s", item.Username, item.ContentId, item.Source, item.Action)}
		})

	go ds.NewEnagements_Mongo(engagements)
	ctx.JSON(http.StatusOK, gin.H{"message": "process initiated"})
}

func main() {
	router := gin.Default()

	auth_group := router.Group("/")
	// authn middleware
	auth_group.Use(func(ctx *gin.Context) {
		token_str := ctx.GetHeader("Authorization")
		if token_str == getInternalAuthToken() {
			ctx.Next()
		} else {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	})

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
