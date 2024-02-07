use mediapulserepo

db.categories.countDocuments({})

db.runCommand({
  createIndexes: 'categories',
  indexes: [
    {
      name: 'category_search',
      key: {
        "embeddings": "cosmosSearch"
      },
      cosmosSearchOptions: {
        kind: 'vector-ivf',
        numLists: 1,
        similarity: 'COS',
        dimensions: 1536
      }
    }
  ]
})

db.runCommand({
  createIndexes: 'mediacontents',
  indexes: [
    {
      name: 'mediacontent_search',
      key: {
        "embeddings": "cosmosSearch"
      },
      cosmosSearchOptions: {
        kind: 'vector-ivf',
        numLists: 1,
        similarity: 'COS',
        dimensions: 1536
      }
    }
  ]
})

