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
