package client

import "time"

type SearchDataRequest struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`

	Measurements []string   `json:"measurements"`
	Tags         []TagMatch `json:"tags"`
}

type TagMatch struct {
	Name     string      `json:"name"`
	Value    interface{} `json:"value"`
	Operator string      `json:"operator"`
}
