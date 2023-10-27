package nationalize

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	BaseURL = "https://api.nationalize.io/"
	Timeout = 5 * time.Second
)

type Client struct {
	cli     *http.Client
	baseURL string
}

type CountryStruct struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type NationalityResponse struct {
	Count   int             `json:"count"`
	Name    string          `json:"name"`
	Country []CountryStruct `json:"country"`
}

func New() *Client {
	return &Client{
		cli: &http.Client{
			Timeout: Timeout,
		},
		baseURL: BaseURL,
	}
}

func (c *Client) Nationality(name string) (string, error) {
	// https://api.nationalize.io/?name=Dmitriy

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/?name=%s", c.baseURL, name), nil)
	if err != nil {
		return "", err
	}

	resp, err := c.cli.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	user := &NationalityResponse{}
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return "", err
	}

	if len(user.Country) == 0 {
		return "", errors.New("returned empty country list")
	}

	return user.Country[0].CountryID, nil
}
