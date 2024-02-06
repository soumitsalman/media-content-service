package mediacontentservice

import "os"

const (
	DB_CONNECTION_STRING = "DB_CONNECTION_STRING"
	DB_NAME              = "mediapulserepo"
	MEDIA_CONTENTS       = "mediacontents"
	USER_ENGAGEMENTS     = "userengagements"
	USER_IDS             = "userids"
	USER_INTERESTS       = "userinterests"
	INTEREST_CATEGORIES  = "categories"
)

func getDBConnectionString() string {
	return os.Getenv(DB_CONNECTION_STRING)
}

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
