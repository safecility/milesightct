package messages

import (
	"encoding/json"
	"time"
)

type ENMessage struct {
	Params json.RawMessage
	Meta   Meta   `json:"meta"`
	Type   ENType `json:"type"`
}

type ENType string

const (
	InfoType            ENType = "info"
	UplinkType          ENType = "uplink"
	DownlinkType        ENType = "downlink"
	DownlinkRequestType ENType = "downlink_request"
)

type Meta struct {
	Network     string  `json:"network"`
	PacketHash  string  `json:"packet_hash"`
	Application string  `json:"application"`
	DeviceAddr  string  `json:"device_addr"`
	Time        float64 `json:"time"`
	Device      string  `json:"device"`
	PacketId    string  `json:"packet_id"`
	Gateway     string  `json:"gateway"`
}

type DeviceMeta struct {
	Application string    `json:"application"`
	DeviceAddr  string    `json:"device_addr"`
	Time        time.Time `json:"time"`
	Device      string    `json:"device"`
	PacketId    string    `json:"packet_id"`
	Source      string
	CompanyID   int64
}

type EverynetUplink struct {
	DeviceMeta
	Uplink UplinkParams
}

type UplinkParams struct {
	RxTime    float64 `json:"rx_time"`
	Port      int     `json:"port"`
	Duplicate bool    `json:"duplicate"`
	Radio     Radio   `json:"radio"`
	CounterUp int     `json:"counter_up"`
	Lora      struct {
		Header      LoraHeader    `json:"header"`
		MacCommands []interface{} `json:"mac_commands"`
	} `json:"lora"`
	Payload          string `json:"payload"`
	EncryptedPayload string `json:"encrypted_payload"`
}

type LoraHeader struct {
	ClassB    bool `json:"class_b"`
	Confirmed bool `json:"confirmed"`
	Adr       bool `json:"adr"`
	Ack       bool `json:"ack"`
	AdrAckReq bool `json:"adr_ack_req"`
	Version   int  `json:"version"`
	Type      int  `json:"type"`
}

type Radio struct {
	Delay      float64    `json:"delay"`
	DataRate   int        `json:"datarate"`
	Modulation Modulation `json:"modulation"`
	Hardware   Hardware   `json:"hardware"`
	Time       float64    `json:"time"`
	Freq       float64    `json:"freq"`
	Size       int        `json:"size"`
}

type Modulation struct {
	Bandwidth int    `json:"bandwidth"`
	Type      string `json:"type"`
	Spreading int    `json:"spreading"`
	CodeRate  string `json:"coderate"`
}

type Hardware struct {
	Status  int     `json:"status"`
	Chain   int     `json:"chain"`
	Tmst    int64   `json:"tmst"`
	Snr     float64 `json:"snr"`
	Rssi    float64 `json:"rssi"`
	Channel int     `json:"channel"`
}

type DownlinkRequestParams struct {
	CounterDown int     `json:"counter_down"`
	MaxSize     int     `json:"max_size"`
	TxTime      float64 `json:"tx_time"`
}

type InfoParams struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
