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

var activeTimers map[string]bool = map[string]bool{}

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

	if _, ok := activeTimers[timer.Name]; ok {
		newTimer, err := timerAPI.pingTimer(timer.Name)
		if err != nil {
			logger.Error.Println("Cannot update remote timer")
			return err
		}
		timer.Seconds = newTimer.Seconds
	} else {
		activeTimers[timer.Client.Id] = true
		if err := timerAPI.createTimer(timer.Name, int(timer.Seconds)); err != nil {
			logger.Error.Println("Cannot create remote timer")
			return err
		}
	}

	ticker := time.NewTicker(time.Duration(timer.Frequency) * time.Second)
	done := make(chan bool)

	var streamerr error = nil
	go func() {
		for {
			select {
			case <-done:
				activeTimers[timer.Client.Id] = false
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
