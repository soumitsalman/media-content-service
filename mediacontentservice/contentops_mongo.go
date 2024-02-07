package mediacontentservice

import (
	ctx "context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMediaContents_Mongo(contents []MediaContentItem) {
	// initialize with global id
	contents = ForEach[MediaContentItem](contents, func(item *MediaContentItem) { item.CreateGlobalId() })

	// check which of these exist
	// for the ones that exist ONLY update the number fields
	content_ids := Extract[MediaContentItem, string](contents, func(item *MediaContentItem) string { return item.GetGlobalId() })
	log.Println(len(content_ids), "contents being processed")
	existing_contents := findItems[MediaContentItem](
		MEDIA_CONTENTS,
		&bson.M{"_id": bson.M{"$in": content_ids}},
		&bson.M{"_id": 1, "score": 1, "comments": 1, "subscribers": 1, "likes": 1, "likes_ratio": 1})
	log.Println(len(existing_contents), "existing contents found")
	existing_contents = ForEach[MediaContentItem](existing_contents, func(item *MediaContentItem) {
		i := Index[MediaContentItem](*item, contents, compareMediaContentItems)
		item.Score = contents[i].Score
		item.Comments = contents[i].Comments
		item.Subscribers = contents[i].Subscribers
		item.ThumbsupCount = contents[i].ThumbsupCount
		item.ThumbsupRatio = contents[i].ThumbsupRatio
	})
	updateMany[MediaContentItem](MEDIA_CONTENTS, existing_contents)

	// for the ones that do not exist
	// create embeddings
	// create categorization
	// create new entry in mongo
	new_contents := Filter[MediaContentItem](contents, func(item MediaContentItem) bool {
		return !In(item, existing_contents, compareMediaContentItems)
	})
	new_contents = ForEach[MediaContentItem](new_contents, func(item *MediaContentItem) {
		item.Embeddings = CreateEmbeddingsForOne(item.Digest)
		item.Category = createMediaContentCategory(item.Embeddings)
		log.Println(item.Kind, item.ChannelName, item.Title, item.Category)
	})
	insertMany[MediaContentItem](MEDIA_CONTENTS, new_contents)
}

func NewEnagements_Mongo(engagements []UserEngagementItem) {
	engagements = ForEach[UserEngagementItem](engagements, func(item *UserEngagementItem) { item.CreateGlobalId() })
	// check which of these exist
	eng_ids := Extract[UserEngagementItem, string](engagements, func(item *UserEngagementItem) string { return item.GlobalId })
	log.Println(len(eng_ids), "engagements being processed")
	existing_contents := findItems[UserEngagementItem](USER_ENGAGEMENTS, &bson.M{"_id": bson.M{"$in": eng_ids}}, nil)
	log.Println(len(existing_contents), "existing engagements found")

	// for the ones that do not exist
	// create new entry in mongo
	new_engs := Filter[UserEngagementItem](engagements, func(item UserEngagementItem) bool {
		return !In(item, existing_contents, func(a, b *UserEngagementItem) bool { return a.GlobalId == b.GlobalId })
	})
	insertMany[UserEngagementItem](USER_ENGAGEMENTS, new_engs)
}

func NewInterests_Mongo(interests []UserInterestItem) {
	// creating embeddings is more expensive so do double check
	// find existing embeddings so that there is no need to do multiple embedding calls since thats more expensive
	cat_names := Extract[UserInterestItem, string](interests, func(item *UserInterestItem) string { return item.Category })
	log.Println(len(cat_names), "interests being processed")
	categories := NewCategories_Mongo(cat_names)

	// for all items all have embeddings - just put them in user interest table
	interests = ForEach[UserInterestItem](interests, func(item *UserInterestItem) {
		cat_i := IndexAny[CategoryItem](categories, func(cat *CategoryItem) bool { return cat.Category == item.Category })
		item.Embeddings = categories[cat_i].Embeddings
		item.Timestamp = float64(time.Now().UnixNano()) / float64(time.Second)
	})
	insertMany[UserInterestItem](USER_INTERESTS, interests)
}

func NewCategories_Mongo(cat_names []string) []CategoryItem {
	existing_cats := findItems[CategoryItem](
		INTEREST_CATEGORIES,
		&bson.M{"_id": bson.M{"$in": cat_names}},
		nil)
	log.Println(len(existing_cats), "categories found")

	new_cat_names := Filter[string](
		cat_names,
		func(item string) bool {
			return !Any[CategoryItem](existing_cats, func(cat *CategoryItem) bool { return cat.Category == item })
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

func createMediaContentCategory(media_content_embeddings []float32) string {
	search_comm := mongo.Pipeline{
		{{"$search", bson.D{
			{"cosmosSearch", bson.D{
				{"vector", media_content_embeddings},
				{"path", "embeddings"},
				{"k", 1}, // return the top item
			}},
		}}},
		{{"$project", bson.D{
			// {"similarityScore", bson.M{"$meta": "searchScore"}},
			{"_id", "$$ROOT._id"},
		}}},
	}
	if cursor, err := getMongoCollection(INTEREST_CATEGORIES).Aggregate(ctx.Background(), search_comm); err != nil {
		return ""
	} else {
		defer cursor.Close(ctx.Background())
		var items []CategoryItem
		cursor.All(ctx.Background(), &items)
		return items[0].Category
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
		Extract[T, any](items, func(item *T) any { return item })); err != nil {
		log.Println("Insertion failed", err)
	} else {
		log.Println(len(res.InsertedIDs), "items inserted in Mongo DB", table)
	}
}

func updateMany[T DataItem](table string, items []T) {
	// this is done for error handling for mongo db
	if len(items) == 0 {
		log.Println("empty list of items. nothing to update")
		return
	}

	coll := getMongoCollection(table)
	if res, err := coll.BulkWrite(
		ctx.Background(),
		Extract[T, mongo.WriteModel](items, func(item *T) mongo.WriteModel {
			return mongo.NewUpdateOneModel().
				SetFilter(bson.M{"_id": (*item).GetGlobalId()}).
				SetUpdate(bson.M{"$set": item})
		})); err != nil {
		log.Println("Update failed", err)
	} else {
		log.Println(res.ModifiedCount, "items updated in Mongo DB", table)
	}
}

func findItems[T any](table string, filter *bson.M, fields *bson.M) []T {
	coll := getMongoCollection(table)
	var cursor *mongo.Cursor
	if fields != nil {
		find_options := options.Find().SetProjection(fields)
		cursor, _ = coll.Find(ctx.Background(), filter, find_options)
	} else {
		cursor, _ = coll.Find(ctx.Background(), filter)
	}
	defer cursor.Close(ctx.Background())
	var items []T
	cursor.All(ctx.Background(), &items)
	return items
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
