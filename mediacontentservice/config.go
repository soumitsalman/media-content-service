package mediacontentservice

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

var media_store_client *aztables.ServiceClient
var table_clients map[string]*aztables.Client = make(map[string]*aztables.Client)

const (
	MEDIASTORE_CONNECTION_STRING = "MEDIASTORE_CONNECTION_STRING"
	MEDIASTORE_CONTENTS_TABLE    = "MEDIASTORE_CONTENTS_TABLE"
	USER_ENGAGEMENTS_TABLE       = "USER_ENGAGEMENTS_TABLE"
	USER_IDS_TABLE               = "USER_IDS_TABLE"
	USER_INTERESTS_TABLE         = "USER_INTERESTS_TABLE"
)

func getMediaStoreClient() *aztables.ServiceClient {
	if media_store_client == nil {
		media_store_client, _ = aztables.NewServiceClientFromConnectionString(os.Getenv(MEDIASTORE_CONNECTION_STRING), nil)
	}
	return media_store_client
}

func getTable(table string) *aztables.Client {
	table_name := os.Getenv(table)
	// return the instance if it already exists
	if table_client, ok := table_clients[table_name]; ok {
		return table_client
	}
	// or else create an instance and return it
	if table_client := getMediaStoreClient().NewClient(table_name); table_client != nil {
		table_clients[table_name] = table_client
		return table_client
	}
	// return nil
	log.Println("Couldn't find the table:", table_name)
	return nil
}

// func getMediaContentsTableName() string {
// 	return os.Getenv("MEDIASTORE_CONTENTS_TABLE")
// }

// func getMediaContentsTable() *aztables.Client {
// 	media_contents_once.Do(func() {
// 		if media_contents_client = getMediaStoreClient().NewClient(getMediaContentsTableName()); media_contents_client == nil {
// 			log.Println("Couldn't find the table")
// 		} else {
// 			log.Println("Table client created")
// 		}
// 	})
// 	return media_contents_client
// }

// func getUserEngagementsTableName() string {
// 	return os.Getenv("USER_ENGAGEMENTS_TABLE")
// }

// func getUserEngagementsTable() *aztables.Client {
// 	user_engagements_once.Do(func() {
// 		if user_engagements_client = getMediaStoreClient().NewClient(getUserEngagementsTableName()); media_contents_client == nil {
// 			log.Println("Couldn't find the table")
// 		} else {
// 			log.Println("Table client created")
// 		}
// 	})
// 	return user_engagements_client
// }

// func getUserIdsTableName() string {
// 	return os.Getenv("USER_IDS_TABLE")
// }

// func getUserEngagementsTable() *aztables.Client {
// 	user_engagements_once.Do(func() {
// 		if user_engagements_client = getMediaStoreClient().NewClient(getUserEngagementsTableName()); media_contents_client == nil {
// 			log.Println("Couldn't find the table")
// 		} else {
// 			log.Println("Table client created")
// 		}
// 	})
// 	return user_engagements_client
// }

// func getUserInterestsTableName() string {
// 	return os.Getenv("USER_INTERESTS_TABLE")
// }
