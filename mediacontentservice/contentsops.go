package mediacontentservice

import (
	ctx "context"
	"encoding/json"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

const TRANSACTION_BATCH_SIZE = 99

func NewEnagements(engagements []UserEngagementItem) {
	engagements_table := getTable(USER_ENGAGEMENTS_TABLE)
	transactions := make([]aztables.TransactionAction, len(engagements))
	for i, eng := range engagements {
		partition_key, row_key := eng.GetKeys()
		serialize_func := func() []byte {
			ent, _ := json.Marshal(eng)
			return ent
		}
		transactions[i] = newTransactionAction(engagements_table,
			partition_key, row_key,
			func() {},
			serialize_func)
	}
	submitSafeTransaction(transactions, engagements_table)
}

func NewContents(contents []MediaContentItem) {
	contents_table := getTable(MEDIASTORE_CONTENTS_TABLE)
	// create a long transaction batch
	transactions := make([]aztables.TransactionAction, len(contents))
	for i, cnt := range contents {
		partition_key, row_key := cnt.GetKeys()
		create_embeddings_func := func() {
			log.Println("TODO: entity does not exist. send for embeddings and categories")
		}
		serialize_func := func() []byte {
			cnt.Children = nil
			ent, _ := json.Marshal(cnt)
			return ent
		}
		transactions[i] = newTransactionAction(contents_table,
			partition_key, row_key,
			create_embeddings_func,
			serialize_func)
	}
	submitSafeTransaction(transactions, contents_table)
}

func newTransactionAction(table *aztables.Client,
	partition_key, row_key string,
	new_item_func func(),
	entity_marshal_func func() []byte) aztables.TransactionAction {

	var action = aztables.TransactionAction{
		ActionType: aztables.TransactionTypeInsertMerge,
		Entity:     entity_marshal_func(),
	}
	if existing, err := table.GetEntity(ctx.Background(), partition_key, row_key, nil); err != nil {
		new_item_func()
	} else {
		action.IfMatch = &existing.ETag
	}
	return action
}

func submitSafeTransaction(transactions []aztables.TransactionAction, table *aztables.Client) {
	for i := 0; i < len(transactions); i += TRANSACTION_BATCH_SIZE {
		batch := safetSlice[aztables.TransactionAction](transactions, i, i+TRANSACTION_BATCH_SIZE)
		if _, err := table.SubmitTransaction(ctx.Background(), batch, nil); err != nil {
			log.Println(err)
		} else {
			log.Println("Succeeded")
		}
	}
}

func safetSlice[T any](array []T, start, noninclusive_end int) []T {
	if start < 0 {
		start = 0
	}
	if noninclusive_end < 0 {
		noninclusive_end = 0
	}
	return array[min(start, len(array)):min(noninclusive_end, len(array))]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
