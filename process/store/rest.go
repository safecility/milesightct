package store

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/safecility/iot/devices/milesightct/process/helpers"
	"github.com/safecility/iot/devices/milesightct/process/messages"
)

type DeviceClient struct {
	client *resty.Client
	server string
}

func CreateDeviceClient(config *helpers.Config) (*DeviceClient, error) {
	if config.Store.Rest == nil {
		return nil, fmt.Errorf("no rest config provided")
	}
	// Create a Resty Client
	client := resty.New()
	return &DeviceClient{
		client: client,
		server: config.Store.Rest.Address(),
	}, nil
}

func (dc *DeviceClient) GetDevice(uid string) (*messages.PowerDevice, error) {
	resp, err := dc.client.R().
		//SetQueryParams(map[string]string{
		//	"page_no": "1",
		//	"limit":   "20",
		//	"sort":    "name",
		//	"order":   "asc",
		//	"random":  strconv.FormatInt(time.Now().Unix(), 10),
		//}).
		SetHeader("Accept", "application/json").
		//SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F").
		Get(fmt.Sprintf("%s/device/%s", dc.server, uid))

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("%s", resp.Status())
	}
	pd := &messages.PowerDevice{}
	err = json.Unmarshal(resp.Body(), pd)
	return pd, err
}

func (dc *DeviceClient) Close() error {
	return nil
}
