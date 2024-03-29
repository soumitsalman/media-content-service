package api

import (
	ctx "context"
	"fmt"
	"log"
	"time"

	utils "github.com/soumitsalman/data-utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMediaContents(contents []MediaContentItem) {
	log.Println(len(contents), "contents being processed")
	// pre-emptive error handling for mongo pipeline
	if len(contents) == 0 {
		log.Println("Terminating early. No contents provided.")
		return
	}

	// get the list of sources and ids into an array
	content_sources, content_ids := make([]string, 0, len(contents)), make([]string, 0, len(contents))
	utils.ForEach[MediaContentItem](contents, func(item *MediaContentItem) {
		content_sources = append(content_sources, item.Source)
		content_ids = append(content_ids, item.Id)
	})

	// check which of these exist
	// for the ones that exist ONLY update the number fields
	// TODO: techically the search filter logic is flawed but since we can assume that the sources will always be the same in the array in practice it wont result in an error
	pipeline := mongo.Pipeline{
		{{
			"$match", bson.M{
				"source": bson.M{"$in": content_sources},
				"cid":    bson.M{"$in": content_ids},
			},
		}},
		{{
			"$project", bson.M{
				"source":      1,
				"cid":         1,
				"score":       1,
				"comments":    1,
				"subscribers": 1,
				"likes":       1,
				"likes_ratio": 1,
			},
		}},
	}
	// existing_contents := findItems[MediaContentItem](MEDIA_CONTENTS, filter, projection)
	existing_contents := findMany[MediaContentItem](MEDIA_CONTENTS, pipeline)
	log.Println(len(existing_contents), "existing contents found")

	// for the existing items, update the numbers with the newly retrieved numbers
	updateMany[MediaContentItem](
		MEDIA_CONTENTS,
		existing_contents,
		getMediaContentIdFilter,
		func(item *MediaContentItem) bson.M {
			return getMediaContentUpdateObj(&contents[utils.Index[MediaContentItem](*item, contents, compareMediaContents)])
		})

	// for the ones that do not exist
	// create embeddings
	// create categorization
	// create new entry in mongo
	new_contents := utils.Filter[MediaContentItem](contents, func(item *MediaContentItem) bool {
		return !utils.In(*item, existing_contents, compareMediaContents)
	})
	new_contents = utils.Filter[MediaContentItem](new_contents, func(item *MediaContentItem) bool {
		if item.Digest != "" {
			item.Excerpt = utils.TruncateTextWithEllipsis(item.Text, MAX_EXCEPRT_SIZE)
			item.Embeddings = CreateEmbeddingsForOne(item.Digest)
			item.Tags = createMediaContentTags(item.Embeddings)
			if item.Category != "" {
				item.Tags = append(item.Tags, item.Category)
			}
			item.Digest = "" //clear out the content. No need to present this
			return true
		}
		return false
	})
	insertMany[MediaContentItem](MEDIA_CONTENTS, new_contents)
}

func NewEnagements(engagements []UserEngagementItem) {
	log.Println(len(engagements), "engagements being processed")
	// rectify UID and filter out the ones without valid user
	engagements = utils.Filter[UserEngagementItem](engagements, func(item *UserEngagementItem) bool {
		if uid, ok := getGlobalUID(item.UserSource, item.Username); ok {
			item.UID = uid
			return true
		}
		log.Printf("User does NOT exist: %s@%s", item.Username, item.UserSource)
		return false
	})
	// pre-emptive error handling for mongo pipeline
	if len(engagements) == 0 {
		log.Println("Terminating Early. No engagements found with valid user names")
		return
	}

	pipeline := mongo.Pipeline{
		{{
			"$match", bson.M{
				"$or": utils.Transform[UserEngagementItem, bson.M](engagements, func(item *UserEngagementItem) bson.M {
					// if UID is not specified then get compute it
					return bson.M{
						"uid":    item.UID,
						"source": item.Source,
						"cid":    item.ContentId,
						"action": item.Action,
					}
				}),
			},
		}},
	}
	existing_engagements := findMany[UserEngagementItem](USER_ENGAGEMENTS, pipeline)
	log.Println(len(existing_engagements), "existing engagements found")

	// for the ones that do not exist
	// create new entry in mongo
	new_engs := utils.Filter[UserEngagementItem](engagements, func(item *UserEngagementItem) bool {
		return !utils.In(*item, existing_engagements, compareUserEngagements)
	})
	insertMany[UserEngagementItem](USER_ENGAGEMENTS, new_engs)
}

func NewInterests(interests []UserInterestItem) {
	// creating embeddings is more expensive so do double check
	// find existing embeddings so that there is no need to do multiple embedding calls since thats more expensive
	cat_names := utils.Transform[UserInterestItem, string](interests, func(item *UserInterestItem) string { return item.Category })
	log.Println(len(cat_names), "interests being processed")
	categories := NewCategories(cat_names)

	// for all items all have embeddings - just put them in user interest table
	interests = utils.ForEach[UserInterestItem](interests, func(item *UserInterestItem) {
		cat_i := utils.IndexAny[CategoryItem](categories, func(cat *CategoryItem) bool { return cat.Category == item.Category })
		item.Embeddings = categories[cat_i].Embeddings
		item.Timestamp = float64(time.Now().UnixNano()) / float64(time.Second)
	})
	insertMany[UserInterestItem](USER_INTERESTS, interests)
}

func NewCategories(cat_names []string) []CategoryItem {
	// existing_cats := findMany[CategoryItem](INTEREST_CATEGORIES, bson.M{"_id": bson.M{"$in": cat_names}}, nil)
	pipeline := mongo.Pipeline{
		{{
			"$match", bson.M{
				"_id": bson.M{"$in": cat_names},
			},
		}},
	}
	existing_cats := findMany[CategoryItem](INTEREST_CATEGORIES, pipeline)
	log.Println(len(existing_cats), "categories found")

	new_cat_names := utils.Filter[string](
		cat_names,
		func(item *string) bool {
			return !utils.Any[CategoryItem](existing_cats, func(cat *CategoryItem) bool { return cat.Category == *item })
		})
	log.Println(len(new_cat_names), "new categories need embeddings")

	if new_embeddings := CreateEmbeddingsForMany(new_cat_names); new_embeddings != nil {
		new_cats := make([]CategoryItem, len(new_cat_names))
		for i := range new_cat_names {
			new_cats[i] = CategoryItem{
				Category:   new_cat_names[i],
				Embeddings: new_embeddings[i],
			}
		}
		insertMany[CategoryItem](INTEREST_CATEGORIES, new_cats)
		existing_cats = append(existing_cats, new_cats...)
	}
	return existing_cats
}

func NewUserCredential(credential UserCredentialItem) string {
	// this is a brand new user
	// just make one up using source and username
	if credential.UID == "" {
		credential.UID = fmt.Sprintf("%s@%s", credential.Username, credential.Source)
	}
	insertOne[UserCredentialItem](USER_IDS, credential)
	return credential.UID
}

func GetAllUserCredentials(source string) []UserCredentialItem {
	pipeline := mongo.Pipeline{
		{{
			"$match", bson.M{"source": source},
		}},
	}
	return findMany[UserCredentialItem](USER_IDS, pipeline)
}

func GetUserContentSuggestions(uid, kind string) []MediaContentItem {
	interests := GetUserInterests(uid)
	engagements := GetUserContentEngagements(uid)

	// use index to pull in contents from media contents
	var match_clause = bson.M{}
	if kind != "*" && kind != "" {
		match_clause["kind"] = kind
	}
	if len(interests) > 0 {
		match_clause["tags"] = bson.M{"$in": interests}
	}
	if len(engagements) > 0 {
		match_clause["$nor"] = utils.Transform[UserEngagementItem, bson.M](engagements, func(item *UserEngagementItem) bson.M {
			return bson.M{
				"source": item.Source,
				"cid":    item.ContentId,
			}
		})
	}

	pipeline := mongo.Pipeline{
		{{
			"$match", match_clause,
		}}, // filter
		{{
			"$project", bson.M{
				"entity":     0,
				"category":   0,
				"score":      0,
				"digest":     0,
				"embeddings": 0,
				"text":       0,
			},
		}}, // projection
		{{
			"$sort", bson.M{
				"comments":    -1,
				"likes":       -1,
				"likes_ratio": -1,
				"subscribers": -1,
				"created":     -1,
			},
		}}, // sort
		{{
			"$limit", 5,
		}}, // top 5
	}

	return findMany[MediaContentItem](MEDIA_CONTENTS, pipeline)
}

func GetUserContentEngagements(uid string) []UserEngagementItem {
	pipeline := mongo.Pipeline{
		{{
			"$match", bson.M{"uid": uid},
		}},
	}
	return findMany[UserEngagementItem](USER_ENGAGEMENTS, pipeline)
}

func GetUserInterests(uid string) []string {
	pipeline := mongo.Pipeline{
		{{
			"$match", bson.M{"uid": uid},
		}},
		{{
			"$project", bson.M{"category": 1},
		}},
	}
	items := findMany[UserInterestItem](USER_INTERESTS, pipeline)
	return utils.Transform[UserInterestItem, string](items, func(item *UserInterestItem) string {
		return item.Category
	})
}

func isValidUID(uid string) bool {
	_, err := findOne[UserCredentialItem](USER_IDS, bson.M{"uid": uid})
	return err == nil
}

func getGlobalUID(source, username string) (string, bool) {
	item, err := findOne[UserCredentialItem](USER_IDS, bson.M{"source": source, "username": username})
	return item.UID, (err == nil)
}

// data object transformers
func getMediaContentIdFilter(item *MediaContentItem) bson.M {
	return bson.M{
		"source": item.Source,
		"cid":    item.Id,
	}
}

func getMediaContentUpdateObj(item *MediaContentItem) bson.M {
	// only update the following fields
	return bson.M{
		"$set": MediaContentItem{
			Score:         item.Score,
			Comments:      item.Comments,
			Subscribers:   item.Subscribers,
			ThumbsupCount: item.ThumbsupCount,
			ThumbsupRatio: item.ThumbsupRatio,
		}}
}

func createMediaContentTags(media_content_embeddings []float32) []string {
	search_comm := mongo.Pipeline{
		{{
			"$search", bson.M{
				"cosmosSearch": bson.M{
					"vector": media_content_embeddings,
					"path":   "embeddings",
					"k":      2,
				}, // return the top item
			},
		}},
		{{
			"$project", bson.M{
				"_id": 1, //"$$ROOT._id"},
			},
		}},
	}
	if cursor, err := getMongoCollection(INTEREST_CATEGORIES).Aggregate(ctx.Background(), search_comm); err != nil {
		return nil
	} else {
		defer cursor.Close(ctx.Background())
		var items []CategoryItem
		cursor.All(ctx.Background(), &items)
		return utils.Transform[CategoryItem, string](items, func(item *CategoryItem) string { return item.Category })
	}
}

// mongo db specific operations
func insertMany[T any](table string, items []T) {
	// this is done for error handling for mongo db
	if len(items) == 0 {
		log.Println("empty list of items. nothing to insert")
		return
	}

	coll := getMongoCollection(table)
	if res, err := coll.InsertMany(
		ctx.Background(),
		utils.Transform[T, any](items, func(item *T) any { return item })); err != nil {
		log.Println("Insertion failed", err)
	} else {
		log.Println(len(res.InsertedIDs), "items inserted in Mongo DB", table)
	}
}

// mongo db specific operations
func insertOne[T any](table string, item T) {
	coll := getMongoCollection(table)
	if res, err := coll.InsertOne(ctx.Background(), item); err != nil {
		log.Println("Insertion failed", err)
	} else {
		log.Println(res.InsertedID, "items inserted in Mongo DB", table)
	}
}

func updateMany[T any](table string, items []T, filter_func func(item *T) bson.M, update_func func(item *T) bson.M) {
	// this is done for error handling for mongo db
	if len(items) == 0 {
		log.Println("empty list of items. nothing to update")
		return
	}

	coll := getMongoCollection(table)
	if res, err := coll.BulkWrite(
		ctx.Background(),
		utils.Transform[T, mongo.WriteModel](items, func(item *T) mongo.WriteModel {
			return mongo.NewUpdateOneModel().SetFilter(filter_func(item)).SetUpdate(update_func(item))
		})); err != nil {
		log.Println("Update failed", err)
	} else {
		log.Println(res.MatchedCount, "items matched in Mongo DB", table)
	}
}

func findOne[T any](table string, filter bson.M) (T, error) {
	var item T
	err := getMongoCollection(table).FindOne(ctx.Background(), filter).Decode(&item)
	return item, err
}

func findMany[T any](table string, pipeline mongo.Pipeline) []T {
	cursor, err := getMongoCollection(table).Aggregate(ctx.Background(), pipeline)
	if err == nil {
		defer cursor.Close(ctx.Background())
		var contents []T
		if err = cursor.All(ctx.Background(), &contents); err == nil {
			return contents
		}
	}
	log.Println("Couldn't retrieve items from", table, "| error:", err)
	return nil
}

// mongo DB and collections clients
var mongo_client *mongo.Client

func getMongoClient() *mongo.Client {
	if mongo_client == nil {
		client_options := options.Client().ApplyURI(getDBConnectionString())
		if new_client, err := mongo.Connect(ctx.Background(), client_options); err != nil {
			log.Println(err)
		} else {
			mongo_client = new_client
			// log.Println("connection succeeded")
		}
	}
	return mongo_client
}

func getMongoDatabase() *mongo.Database {
	var client = getMongoClient()
	if client == nil {
		return nil
	}
	return client.Database(DB_NAME)
}

func getMongoCollection(name string) *mongo.Collection {
	var db = getMongoDatabase()
	if db == nil {
		return nil
	}
	return db.Collection(name)
}
