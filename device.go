package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Device struct {
	ID   uint64 `json:"id"`
	Tags []Tag  `json:"tags"`
}

type Tag struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func (c *Client) GetDevices(ctx context.Context) ([]Device, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.hyprcubd.com/v1/devices/search", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println("Hyprcubd SearchDevices returned ", resp.StatusCode)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received %d from Hyprcubd", resp.StatusCode)
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	devices := []Device{}
	err = json.Unmarshal(out, &devices)
	return devices, err
}

type CreateDevicesResponse struct {
	Devices []uint64 `json:"ids"`
}

func (c *Client) CreateDevices(ctx context.Context, devices []Device) ([]uint64, error) {
	out, err := json.Marshal(&devices)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.hyprcubd.com/v1/devices", bytes.NewReader(out))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received %d from Hyprcubd", resp.StatusCode)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cdr := CreateDevicesResponse{}
	err = json.Unmarshal(buf, &cdr)
	if err != nil {
		return nil, err
	}

	return cdr.Devices, nil
}
