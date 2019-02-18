package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrl = `https://api.stackexchange.com/2.2/questions?pagesize=%d&order=desc&sort=activity&tagged=%s&site=%s`
)

type question struct {
	CreationDate float64 `json:"creation_date"`
	QuestionID   int64   `json:"question_id"`
	Title        string  `json:"title"`
	Link         string  `json:"link"`
}

type apiClient struct {
	httpClient *http.Client
	site       string
}

func (c *apiClient) getQuestions(max uint, tag string) (result []question, err error) {
	r, err := c.httpClient.Get(fmt.Sprintf(apiUrl, max, tag, c.site))
	if err != nil {
		return nil, fmt.Errorf("failed to get data from %s API: %s", c.site, err)
	}

	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&result)
}
