package feed

type Selector struct {
	Count   uint
	AfterID string
}

type Posts []Post

type Post struct {
	ID         string  `json:"id"`
	URL        string  `json:"url"`
	Title      string  `json:"title"`
	Text       string  `json:"text"`
	ImageURL   string  `json:"imageUrl"`
	SourceType string  `json:"source"`
	CreatedAt  float64 `json:"createdAt"`
}
