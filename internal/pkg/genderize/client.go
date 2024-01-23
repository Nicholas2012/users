package genderize

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	Timeout = 5 * time.Second
)

type Client struct {
	cli     *http.Client
	baseURL string
}

func New(apiURL string) *Client {
	return &Client{
		cli: &http.Client{
			Timeout: Timeout,
		},
		baseURL: apiURL,
	}
}

type GenderResponse struct {
	Count  int    `json:"count"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

func (c *Client) Gender(name string) (string, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/?name=%s", c.baseURL, name), nil)
	if err != nil {
		return "", err
	}

	resp, err := c.cli.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	user := &GenderResponse{}
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return "", err
	}

	g := "f"
	if string(user.Gender) == "male" {
		g = "m"
	}

	return g, nil
}
