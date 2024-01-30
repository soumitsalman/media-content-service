package main

import (
	"encoding/json"
	"log"
	"os"

	ds "github.com/soumitsr/media-content-service/mediacontentservice"
	// "github.com/otiai10/openaigo"
)

const REDDIT = "REDDIT"

func _loadSeedFile() {
	contents := readJson[[]ds.MediaContentItem]("/home/soumitsr/Codes/goreddit/contents.json")
	ds.NewContents(contents[16:27])

	loaded_e := readJson[[]map[string]any]("/home/soumitsr/Codes/goreddit/engagements.json")
	// preprocessing for now. in future this should be a different format given from the collector
	engagements := make([]ds.UserEngagementItem, 0, 30)
	for _, e := range loaded_e {
		for id, action := range e["engagements"].(map[string]any) {
			engagements = append(engagements, ds.UserEngagementItem{
				Username:  e["username"].(string),
				Source:    REDDIT,
				ContentId: id,
				Action:    action.(string),
			})
		}
	}
	log.Println(len(engagements))
	ds.NewEnagements(engagements)
}

func readJson[T any](path string) T {
	content_bytes, _ := os.ReadFile(path)
	var contents T
	json.Unmarshal(content_bytes, &contents)
	return contents
}

func main() {
	_loadSeedFile()
}

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
