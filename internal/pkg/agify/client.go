package agify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	BaseURL = "https://api.agify.io/"
	Timeout = 5 * time.Second
)

type Client struct {
	cli     *http.Client
	baseURL string
}

func New() *Client {
	return &Client{
		cli: &http.Client{
			Timeout: Timeout,
		},
		baseURL: BaseURL,
	}
}

type AgeResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func (c *Client) Age(name string) (int, error) {
	// https://api.agify.io/?name=Dmitriy

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/?name=%s", c.baseURL, name), nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.cli.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	user := &AgeResponse{}
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return 0, err
	}

	return user.Age, nil
}
