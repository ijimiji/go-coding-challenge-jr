package server

import (
	"challenge/pkg/logger"
	"challenge/pkg/proto"
	"encoding/json"
	"errors"
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

// send http request to get remaining time of timer `name`
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

	// handle nonexistent timers on 404 error code for example
	if resp.StatusCode != 200 {
		return &proto.Timer{}, errors.New("non 200 return code")
	}

	values := struct {
		Timer            string  `json:"timer"`
		SecondsRemaining float64 `json:"seconds_remaining"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&values); err != nil {
		logger.Error.Println(values)
		return nil, err
	}

	return &proto.Timer{
		Name:    values.Timer,
		Seconds: int64(values.SecondsRemaining),
	}, nil
}

func (s *ChallengeServer) StartTimer(timer *proto.Timer, stream proto.ChallengeService_StartTimerServer) error {
	logger.Info.Println("Starting timer...")

	t, err := timerAPI.pingTimer(timer.Name)

	// timer doesn't even exist or run out
	if err != nil || t.Seconds == 0 {
		if err := timerAPI.createTimer(timer.Name, int(timer.Seconds)); err != nil {
			return err
		}
	} else {
		// timer exists, running and clients should resubscribe to it

		// replace seconds to create correct ticker below
		timer.Seconds = t.Seconds
	}

	ticker := time.NewTicker(time.Duration(timer.Frequency) * time.Second)
	done := make(chan bool)

	var streamerr error = nil
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				logger.Info.Println("tick")
				newTimer, err := timerAPI.pingTimer(timer.Name)
				if err != nil {
					streamerr = err
					done <- true
				}
				if newTimer.Seconds == 0 {
					done <- true
				}
				stream.SendMsg(newTimer)
			}
		}
	}()

	time.Sleep(time.Duration(timer.Seconds) * time.Second)
	ticker.Stop()
	done <- true

	return streamerr
}
