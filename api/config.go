package api

import "os"

const (
	DB_NAME             = "mediapulserepo"
	MEDIA_CONTENTS      = "mediacontents"
	USER_ENGAGEMENTS    = "userengagements"
	USER_IDS            = "userids"
	USER_INTERESTS      = "userinterests"
	INTEREST_CATEGORIES = "categories"
)

const (
	COMPLETION_MODEL = "gpt-3.5-turbo"
	EMBEDDINGS_MODEL = "text-embedding-3-small"
	MAX_TOKEN_LIMIT  = 4096
	MAX_EXCEPRT_SIZE = 400 //400 chars is about 100 tokens
)

func getDBConnectionString() string {
	return os.Getenv("DB_CONNECTION_STRING")
}

func getInternalAuthToken() string {
	return os.Getenv("INTERNAL_AUTH_TOKEN")
}

func getOpenAIApiKey() string {
	return os.Getenv("OPENAI_API_KEY")
}

func getOpenAIBaseUrl() string {
	return os.Getenv("OPENAI_BASE_URL")
}

func getOpenAIOrgID() string {
	return os.Getenv("OPENAI_ORG_ID")
}
