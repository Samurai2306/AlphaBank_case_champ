package rag

// Chunk is one knowledge unit from the demo corpus.
type Chunk struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Source string   `json:"source"`
	Tags   []string `json:"tags"`
	Text   string   `json:"text"`
}

// Hit is a ranked retrieval result matching the retrieve_kb contract.
type Hit struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Text   string  `json:"text"`
	Source string  `json:"source"`
	Score  float64 `json:"score"`
}
