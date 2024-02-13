// delete this file later
package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	ds "github.com/soumitsr/media-content-service/mediacontentservice"
	// "github.com/otiai10/openaigo"
)

const REDDIT = "REDDIT"

func _loadSeedFile() {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoFormatHeaders(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	// contents := readJson[[]ds.MediaContentItem]("/home/soumitsr/Codes/goreddit/contents.json")
	// ds.NewMediaContents_Mongo(contents[0:100])

	// engagements := readJson[[]ds.UserEngagementItem]("/home/soumitsr/Codes/goreddit/engagements.json")
	// ds.NewEnagements_Mongo(engagements)

	// raw_interests := readCsv("/home/soumitsr/Codes/media-content-service/Seed Data for Cafecito - user_interests.csv")
	// var user_interests = ds.Extract[[]string, ds.UserInterestItem](
	// 	raw_interests,
	// 	func(item *[]string) ds.UserInterestItem {
	// 		return ds.UserInterestItem{
	// 			GlobalUID: (*item)[0],
	// 			Category:  (*item)[1],
	// 		}
	// 	},
	// )
	// ds.NewInterests_Mongo(user_interests[1:])

	table.SetHeader([]string{"Kind", "Channel", "Tags", "Created", "Subscribers", "Comments", "Likes"})
	ds.ForEach[ds.MediaContentItem](ds.GetUserContentSuggestions("danny_004", "post"), func(item *ds.MediaContentItem) {
		table.Append([]string{item.Kind, item.ChannelName, strings.Join(item.Tags, ", "), ds.DateToString(item.Created), strconv.Itoa(item.Subscribers), strconv.Itoa(item.Comments), strconv.Itoa(item.ThumbsupCount)})
	})
	// table.Render()

	// table.SetHeader([]string{"UID", "Source", "Username", "PW"})
	// ds.ForEach[ds.UserCredentialItem](ds.GetAllUserCredentials("SLACK"), func(item *ds.UserCredentialItem) {
	// 	table.Append([]string{item.UID, item.Source, item.Username, item.Password})
	// })
	table.Render()

}

func _tryEmbeddings() {
	contents := readJson[[]ds.MediaContentItem]("/home/soumitsr/Codes/goreddit/contents.json")

	for _, cnt := range contents[0:5] {
		res := ds.CreateEmbeddingsForOne(cnt.Text)
		log.Println(len(res), res[0], res[len(res)-1])
	}

}

func readJson[T any](path string) T {
	content_bytes, _ := os.ReadFile(path)
	var contents T
	json.Unmarshal(content_bytes, &contents)
	return contents
}

func readCsv(path string) [][]string {
	file, _ := os.Open(path)
	defer file.Close()

	reader := csv.NewReader(file)
	items, _ := reader.ReadAll()
	return items
}

// cats := []string{
// 	"Life",
// 	"Life Lessons",
// 	"Politics",
// 	"Travel",
// 	"Poetry",
// 	"Entrepreneurship",
// 	"Education",
// 	"Health",
// 	"Love",
// 	"Design",
// 	"Writing",
// 	"Technology",
// 	"Self Improvement",
// 	"Business",
// 	"Music",
// 	"Social Media",
// 	"Sports",
// 	"Food",
// 	"Art",
// 	"Python",
// 	"JavaScript",
// 	"Java",
// 	"C#",
// 	"C++",
// 	"C",
// 	"Ruby",
// 	"PHP",
// 	"Swift",
// 	"Kotlin",
// 	"Go",
// 	"Rust",
// 	"TypeScript",
// 	"Scala",
// 	"Perl",
// 	"R",
// 	"Dart",
// 	"Model Architecture",
// 	"Natural Language Processing",
// 	"Machine Learning Algorithms",
// 	"Ethical Implications of AI",
// 	"AI in Healthcare",
// 	"AI in Education",
// 	"AI and Creativity",
// 	"AI Policy and Regulation",
// 	"Data Privacy in AI",
// 	"Future of Work with AI",
// 	"OpenAI's GPT",
// 	"Anyscale",
// 	"OctoAI",
// 	"Open Source LLM",
// 	"Open Source Software",
// 	"AI and Robotics",
// 	"Bias and Fairness in AI",
// 	"AI in Financial Services",
// 	"Environmental Impact of AI",
// }
// ds.NewCategories_Mongo(cats)

// creds, err := azidentity.NewDefaultAzureCredential(nil)

// if err != nil {
// 	log.Println("error authenticating")
// 	return
// }
// log.Println(creds)

// client, err := aztables.NewServiceClientFromConnectionString(getMediaStoreConnection(), nil)
// if err != nil {
// 	log.Println("Couldn't find the connection string")
// 	return
// }
// table := client.NewClient(getMediaContentsTableName())
// if table == nil {
// 	log.Println("Couldn't find the table")
// 	return
// }

// filter := "PartitionKey eq 'channel' and (num_subscribers ge 100)"
// listOptions := aztables.ListEntitiesOptions{
// 	Filter: &filter,
// }
// pager := table.NewListEntitiesPager(&listOptions)
// if pager == nil {
// 	log.Println("Couldn't find items")
// 	return
// }

// for pager.More() {
// 	resp, _ := pager.NextPage(context.Background())
// 	for _, entity := range resp.Entities {
// 		// var data aztables.EDMEntity
// 		// data.UnmarshalJSON(entity)
// 		// log.Println(data.Properties["name"], data.Properties["title"], data.Properties["num_subscribers"])
// 		var media_item mediacontentservice.MediaContentItem
// 		json.Unmarshal(entity, &media_item)
// 		log.Println(media_item.Name, media_item.Title, media_item.NumSubscribers)

// 		// var media_item map[string]any
// 		// json.Unmarshal(entity, &media_item)
// 		// log.Println(media_item)

// 	}
// }
// var data map[string]any
// json.Unmarshal(pager.Value, &data)
// log.Println(data)

// cnt.Children = nil
// cnt.Entity = aztables.Entity{
// 	PartitionKey: cnt.Source,
// 	RowKey:       cnt.Id,
// }
// ent, _ := json.Marshal(cnt)

// cnt_bytes, _ := json.Marshal(cnt)
// var json_obj map[string]any
// json.Unmarshal(cnt_bytes, &json_obj)

// new_entity := aztables.EDMEntity{
// 	Entity: aztables.Entity{
// 		PartitionKey: cnt.Source,
// 		RowKey:       cnt.Id,
// 	},
// 	Properties: json_obj,
// }

// ent, _ := new_entity.MarshalJSON()

// if _, err := contents_table.AddEntity(ctx.Background(), ent, nil); err != nil {
// 	log.Println(err)
// } else {
// 	log.Println("Added")
// }
// cnt.Children = nil
// cnt.Entity = aztables.Entity{
// 	PartitionKey: cnt.Source,
// 	RowKey:       cnt.Id,
// }
// ent, _ := json.Marshal(cnt)

// if _, err := contents_table.UpdateEntity(
// 	ctx.Background(),
// 	ent,
// 	&aztables.UpdateEntityOptions{UpdateMode: aztables.UpdateModeMerge, IfMatch: &existing.ETag}); err != nil {
// 	log.Println(err)
// } else {
// 	log.Println("Updated")
// }

// table cannot contain children
// batch[i] = aztables.TransactionAction{
// 	Entity:     ent,
// 	ActionType: aztables.TransactionTypeAdd, //this should merge
// }

// for i := 0; i < len(contents); i += TRANSACTION_BATCH_SIZE {
// 	//batch_end := i + TRANSACTION_BATCH_SIZE
// 	content_slice := safetSlice[ds.MediaContentItem](contents, i, i+TRANSACTION_BATCH_SIZE)
// 	transaction_batch := make([]aztables.TransactionAction, len(content_slice))

// 	for j, cnt := range content_slice {
// 		// log.Println("Current entity:", cnt.Id, cnt.Title)
// 		if existing, err := contents_table.GetEntity(ctx.Background(), cnt.Source, cnt.Id, nil); err != nil {
// 			// log.Println("TODO: entity does not exist. send for embeddings and categories", existing)
// 			transaction_batch[j] = aztables.TransactionAction{
// 				ActionType: aztables.TransactionTypeAdd,
// 			}
// 		} else {
// 			// log.Println("entity exists. Etag:", existing.ETag)
// 			transaction_batch[j] = aztables.TransactionAction{
// 				ActionType: aztables.TransactionTypeUpdateMerge,
// 				IfMatch:    &existing.ETag,
// 			}
// 		}
// 		cnt.Children = nil
// 		cnt.Entity = aztables.Entity{
// 			PartitionKey: cnt.Source,
// 			RowKey:       cnt.Id,
// 		}
// 		transaction_batch[j].Entity, _ = json.Marshal(cnt)
// 	}

// 	if _, err := contents_table.SubmitTransaction(ctx.Background(), transaction_batch, nil); err != nil {
// 		log.Println(err)
// 	} else {
// 		log.Println("Succeeded")
// 	}
// }

// var batch = make([]aztables.TransactionAction, len(contents))
// for i, cnt := range contents {
// 	log.Println("Current entity:", cnt.Id, cnt.Title)
// 	if existing, err := contents_table.GetEntity(ctx.Background(), cnt.Source, cnt.Id, nil); err != nil {

// 		log.Println("TODO: entity does not exist. send for embeddings and categories", existing)

// 		batch[i] = aztables.TransactionAction{
// 			ActionType: aztables.TransactionTypeAdd,
// 		}
// 	} else {

// 		log.Println("entity exists. Etag:", existing.ETag)
// 		batch[i] = aztables.TransactionAction{
// 			ActionType: aztables.TransactionTypeUpdateMerge,
// 			IfMatch:    &existing.ETag,
// 		}
// 	}
// 	cnt.Children = nil
// 	cnt.Entity = aztables.Entity{
// 		PartitionKey: cnt.Source,
// 		RowKey:       cnt.Id,
// 	}
// 	batch[i].Entity, _ = json.Marshal(cnt)

// }

// if resp, err := contents_table.SubmitTransaction(ctx.Background(), batch, nil); err != nil {
// 	log.Println(err)
// } else {
// 	log.Println(resp)
// }
