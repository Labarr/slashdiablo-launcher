package ladder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ladderAddress = "https://ladder.slashdiablo.net"
)

// Client encapsulates the details of the Slashdiablo API.
type Client interface {
	GetLadder(mode string) ([]ladderCharacter, error)
}

type client struct{}

type ladderResponse struct {
	Characters []ladderCharacter `json:"characters"`
}

type ladderCharacter struct {
	Name  string
	Class string
	Level int
	Rank  int
}

func (c *client) GetLadder(mode string) ([]ladderCharacter, error) {
	response, err := c.do(http.MethodGet, ladderAddress+"/ladder/rankings/"+mode+"?class=overall", nil)
	if err != nil {
		return nil, err
	}

	var resp ladderResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, err
	}

	return resp.Characters, nil
}

func (c *client) do(method string, addr string, payload []byte) ([]byte, error) {
	// Setup http client.
	client := http.Client{}

	r, err := http.NewRequest(method, addr, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	// Add headers.
	r.Header.Add("Content-Type", "application/json; charset=UTF-8")

	resp, err := client.Do(r)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status code %v", resp.StatusCode)
	}

	return responseBody, nil
}

// NewClient returns a new bigbank client with all dependencies.
func NewClient() Client {
	return &client{}
}
