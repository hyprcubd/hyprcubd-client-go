package hyprcubd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// TimeFormat is a nanosecond format
const TimeFormat = "2006-01-02 15:04:05.999999999"

// InsertDataRequest defines a set of series to insert
type InsertDataRequest struct {
	IntSeries []IntSeries `json:"intSeries"`
}

// IntSeries is an integer based series
type IntSeries struct {
	Name         string           `json:"name"`
	Measurements []IntMeasurement `json:"measurements"`
}

// IntMeasurement is a single time and value measurement
type IntMeasurement struct {
	Time  int64 `json:"time"`
	Value int64 `json:"val"`
}

// InsertData sends data to Hyprcubd for the given device
func (c *Client) InsertData(ctx context.Context, deviceID uint64, idr InsertDataRequest) error {
	buf, err := json.Marshal(&idr)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.hyprcubd.com/v1/device/%d/data", deviceID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		// body, err := ioutil.ReadAll(resp.Body)
		// if err == nil {
		// 	resp.Body.Close()
		// }
		return fmt.Errorf("received %d from Hyprcubd", resp.StatusCode)
	}
	return nil
}

type SearchDataRequest struct {
	StartTime    time.Time  `json:"startTime"`
	EndTime      time.Time  `json:"endTime"`
	Measurements []string   `json:"measurements"`
	Tags         []TagMatch `json:"tags"`
}

type TagMatch struct {
	Name     string      `json:"name"`
	Value    interface{} `json:"value"`
	Operator string      `json:"operator"`
}

type SearchDataResponse struct {
	Columns []ResultColumn  `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

type ResultColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// SearchData queries Hyprcubd given the request parameters
func (c *Client) SearchData(ctx context.Context, sdr SearchDataRequest) (*SearchDataResponse, error) {
	buf, err := json.Marshal(&sdr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.hyprcubd.com/v1/data/search", bytes.NewReader(buf))
	if err != nil {
		log.Println("CreateRequest error")
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Hyprcubd returned %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	sdresp := SearchDataResponse{}
	err = json.Unmarshal(body, &sdresp)
	return &sdresp, err
}
