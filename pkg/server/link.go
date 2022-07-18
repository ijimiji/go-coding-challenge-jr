// This part of server module is responsible for shortening links
package server

import (
	"bytes"
	"challenge/pkg/config"
	"challenge/pkg/proto"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

var cfg = config.Get()

type API struct {
	Client  *http.Client
	BaseURL string
}

var bitlyAPI = API{
	Client:  &http.Client{},
	BaseURL: "https://api-ssl.bitly.com/v4",
}

func (api *API) shortenLink(link string) (string, error) {

	// "long_url" field is a link to be shortened
	form, _ := json.Marshal(map[string]string{"long_url": link})
	req, _ := http.NewRequest("POST", api.BaseURL+"/shorten", bytes.NewBuffer(form))

	// bitly requires authentication via token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+cfg.Token)

	// call Do 'cause I don't know any other way to pass header to request
	resp, err := api.Client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return link, errors.New("return code is not 200")
	}

	// ignore everything but "link"
	values := struct {
		Link string
	}{}
	err = json.NewDecoder(resp.Body).Decode(&values)
	defer resp.Body.Close()

	return values.Link, err
}

func (s *ChallengeServer) MakeShortLink(_ context.Context, link *proto.Link) (*proto.Link, error) {

	shortenedLink, err := bitlyAPI.shortenLink(link.GetData())
	if err != nil {
		return link, err
	}

	return &proto.Link{Data: shortenedLink}, nil
}
