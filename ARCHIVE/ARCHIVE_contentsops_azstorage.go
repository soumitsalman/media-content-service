// archive this file
package ARCHIVE

// import (
// 	ctx "context"
// 	"encoding/json"
// 	"log"
// 	"os"

// 	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
// 	utils "github.com/soumitsalman/data-utils"
// 	ds "github.com/soumitsalman/media-content-service/api"
// )

// const TRANSACTION_BATCH_SIZE = 99

// func NewEnagements_AzTables(engagements []ds.UserEngagementItem) {
// 	engagements_table := getAzureTable(ds.USER_ENGAGEMENTS)
// 	transactions := make([]aztables.TransactionAction, len(engagements))
// 	for i, eng := range engagements {
// 		partition_key, row_key := eng.CreateKeys()
// 		serialize_func := func() []byte {
// 			ent, _ := json.Marshal(eng)
// 			return ent
// 		}
// 		transactions[i] = newTransactionAction(engagements_table,
// 			partition_key, row_key,
// 			serialize_func)
// 	}
// 	submitSafeTransaction(transactions, engagements_table)
// }

// func NewContents_AzTable(contents []ds.MediaContentItem) {
// 	contents_table := getAzureTable(ds.MEDIA_CONTENTS)
// 	// create a long transaction batch
// 	transactions := make([]aztables.TransactionAction, len(contents))
// 	for i, cnt := range contents {
// 		partition_key, row_key := cnt.CreateKeys()
// 		ensureEmbeddingsAndCategorizies(&cnt)
// 		serialize_func := func() []byte {
// 			ent, _ := json.Marshal(cnt)
// 			return ent
// 		}
// 		transactions[i] = newTransactionAction(contents_table,
// 			partition_key, row_key,
// 			serialize_func)
// 	}
// 	submitSafeTransaction(transactions, contents_table)
// }

// func ensureEmbeddingsAndCategorizies(item *ds.MediaContentItem) {
// 	vectors := ds.CreateEmbeddingsForOne(item.Digest)
// 	log.Println(len(vectors), vectors[0], vectors[len(vectors)-1])
// }

// func newTransactionAction(table *aztables.Client,
// 	partition_key, row_key string,
// 	entity_marshal_func func() []byte) aztables.TransactionAction {

// 	var action = aztables.TransactionAction{
// 		ActionType: aztables.TransactionTypeInsertMerge,
// 		Entity:     entity_marshal_func(),
// 	}
// 	// look for etag
// 	if existing, err := table.GetEntity(ctx.Background(), partition_key, row_key, nil); err == nil {
// 		// item exists. so assign the etag
// 		action.IfMatch = &existing.ETag
// 	}
// 	return action
// }

// func submitSafeTransaction(transactions []aztables.TransactionAction, table *aztables.Client) {
// 	for i := 0; i < len(transactions); i += TRANSACTION_BATCH_SIZE {
// 		batch := utils.SafeSlice[aztables.TransactionAction](transactions, i, i+TRANSACTION_BATCH_SIZE)
// 		if _, err := table.SubmitTransaction(ctx.Background(), batch, nil); err != nil {
// 			log.Println(err)
// 		} else {
// 			log.Println("Succeeded")
// 		}
// 	}
// }

// // azure storage DB and tables client

// var azure_storage_client *aztables.ServiceClient

// func getAzureStorageClient() *aztables.ServiceClient {
// 	if azure_storage_client == nil {
// 		if client, err := aztables.NewServiceClientFromConnectionString(os.Getenv("DB_CONNECTION_STRING"), nil); err != nil {
// 			log.Println(err)
// 		} else {
// 			azure_storage_client = client
// 		}
// 	}
// 	return azure_storage_client
// }

// func getAzureTable(table string) *aztables.Client {
// 	client := getAzureStorageClient()
// 	if client == nil {
// 		return nil
// 	}
// 	// return nil

// 	return client.NewClient(table)
// }
