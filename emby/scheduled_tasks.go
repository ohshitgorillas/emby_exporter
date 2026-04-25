package emby

import (
	"fmt"
	"net/http"
)

type ScheduledTaskResult struct {
	StartTimeUtc string `json:"StartTimeUtc"`
	EndTimeUtc   string `json:"EndTimeUtc"`
	Status       string `json:"Status"`
	Name         string `json:"Name"`
	Key          string `json:"Key"`
	Id           string `json:"Id"`
	ErrorMessage string `json:"ErrorMessage"`
	LongErrorMessage string `json:"LongErrorMessage"`
}

type ScheduledTask struct {
	Name                string               `json:"Name"`
	State               string               `json:"State"`
	CurrentProgressPct  float64              `json:"CurrentProgressPercentage"`
	Id                  string               `json:"Id"`
	LastExecutionResult *ScheduledTaskResult `json:"LastExecutionResult"`
	Description         string               `json:"Description"`
	Category            string               `json:"Category"`
	IsHidden            bool                 `json:"IsHidden"`
	Key                 string               `json:"Key"`
}

func (c *Client) GetScheduledTasks() ([]ScheduledTask, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ScheduledTasks", c.URL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(c.ctx)

	res := []ScheduledTask{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}
