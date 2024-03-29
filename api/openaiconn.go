package api

import (
	ctx "context"
	"log"

	openai "github.com/otiai10/openaigo"
	utils "github.com/soumitsalman/data-utils"
	"github.com/tiktoken-go/tokenizer"
)

var openai_client *openai.Client

func getOpenaiClient() *openai.Client {
	if openai_client == nil {
		openai_client = openai.NewClient(getOpenAIApiKey())
		openai_client.BaseURL = getOpenAIBaseUrl()
		openai_client.Organization = getOpenAIOrgID()
	}
	return openai_client
}

func CreateEmbeddingsForMany(text_array []string) [][]float32 {
	// this is done for pre-emptive error handling for APIs
	if len(text_array) == 0 {
		return nil
	}
	resp, err := getOpenaiClient().CreateEmbedding(
		ctx.Background(),
		openai.EmbeddingCreateRequestBody{
			Model: EMBEDDINGS_MODEL,
			Input: text_array,
		})
	if err != nil {
		log.Println(err)
		return nil
	}
	return utils.Transform[openai.EmbeddingData, []float32](
		resp.Data,
		func(data *openai.EmbeddingData) []float32 { return data.Embedding })
}

func CreateEmbeddingsForOne(text string) []float32 {
	resp, err := getOpenaiClient().CreateEmbedding(
		ctx.Background(),
		openai.EmbeddingCreateRequestBody{
			Model: EMBEDDINGS_MODEL,
			// TODO: may be I should create multiple chunks
			Input: []string{truncateTextForModel(text, EMBEDDINGS_MODEL)},
		})
	if err != nil {
		return nil
	}
	return resp.Data[0].Embedding
}

// all openai currently uses cl100kBase
// TODO: update this for 4 chars a token
func truncateTextForModel(text string, model string) string {
	// the library is not updated with all the new models like text-embedding-3-small
	enc, err := tokenizer.ForModel(tokenizer.Model(model))
	if err != nil {
		// use cl100k_base as default
		enc, _ = tokenizer.Get(tokenizer.Cl100kBase)
	}
	tokens, _, _ := enc.Encode(text)
	res, _ := enc.Decode(utils.SafeSlice[uint](tokens, 0, MAX_TOKEN_LIMIT))
	return res
}
