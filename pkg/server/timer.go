package server

import (
	"challenge/pkg/logger"
	"challenge/pkg/proto"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var timerAPI = API{
	BaseURL: "https://timercheck.io",
	Client:  &http.Client{},
}

// create remote timer with name: `name` and duration: `duration`
func (api *API) createTimer(name string, duration int) error {
	resp, err := api.Client.Get(fmt.Sprintf("%s/%s/%d", api.BaseURL, name, duration))

	if err != nil || resp.StatusCode != 200 {
		return err
	}

	return nil
}

func (api *API) pingTimer(name string) (*proto.Timer, error) {
	resp, err := api.Client.Get(fmt.Sprintf("%s/%s", api.BaseURL, name))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// treat 504 code differently cause when timer is out body is different
	if resp.StatusCode == 504 {
		return &proto.Timer{
			Name:    name,
			Seconds: 0,
		}, nil
	}

	values := struct {
		Timer            string  `json:"timer"`
		SecondsRemaining float64 `json:"seconds_remaining"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&values); err != nil {
		logger.Error.Fatalln(values)
		return nil, err
	}

	return &proto.Timer{
		Name:    values.Timer,
		Seconds: int64(values.SecondsRemaining),
	}, nil
}

func (s *ChallengeServer) StartTimer(timer *proto.Timer, stream proto.ChallengeService_StartTimerServer) error {
	logger.Info.Println("Starting timer...")
	if err := timerAPI.createTimer(timer.GetName(), int(timer.GetSeconds())); err != nil {
		logger.Error.Fatalln("Cannot create remote timer")
	}

	var running = true

	// it's very very dirty, but I'll try to fix it
	for running {
		logger.Info.Println("tick")

		time.Sleep(time.Duration(timer.Frequency) * time.Second)
		resp, err := timerAPI.pingTimer(timer.GetName())
		if err != nil {
			logger.Error.Fatalln(err)
		}

		if resp.GetSeconds() == 0 {
			running = false
		}

		if err := stream.Send(resp); err != nil {
			logger.Error.Fatalln("Error while sending message to timer stream")
		}
	}

	return nil
}
