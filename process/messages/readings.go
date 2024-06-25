package messages

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

const (
	metaChannel    byte = 0xff
	totalChannel   byte = 0x03
	currentChannel byte = 0x04
	alarmChannel   byte = 0x84
	power          byte = 0x0b
	ipso           byte = 0x01
	serial         byte = 0x16
	hardware       byte = 0x09
	firmware       byte = 0x0a
)

type MilesightCTReading struct {
	*PowerDevice
	UID   string
	Power bool
	Time  time.Time
	Version
	Current
}

type Version struct {
	Ipso     string
	Hardware string
	Firmware string
}

type Current struct {
	Total float32
	Value float32
	Max   float32
	Min   float32
	Alarms
}

type Alarms struct {
	t  bool
	tr bool
	r  bool
	rr bool
}

func ReadMilesightCT(payload []byte) (*MilesightCTReading, error) {
	dpi := &MilesightCTReading{}

	read := 0
	available := len(payload)
	for available > read {
		nr, err := readSlice(dpi, payload, read)
		if err != nil {
			return nil, err
		}
		read += nr
	}

	return dpi, nil
}

func readSlice(r *MilesightCTReading, payload []byte, offset int) (int, error) {
	available := len(payload) - offset
	if available < 2 {
		return 0, fmt.Errorf("payload too small")
	}
	channelID := payload[offset]
	channelType := payload[offset+1]
	switch channelID {
	case metaChannel:
		switch channelType {
		case power:
			r.Power = true
			return 3, nil
		case ipso:
			if available < 3 {
				return 0, fmt.Errorf("payload too small")
			}
			r.Version.Ipso = readVersion(payload[offset+2])
			return 3, nil
		case hardware:
			if available < 3 {
				return 0, fmt.Errorf("payload too small")
			}
			r.Version.Hardware = readVersionLong(payload[offset+2:])
			return 4, nil
		case firmware:
			if available < 4 {
				return 0, fmt.Errorf("payload too small")
			}
			r.Version.Firmware = readVersionLong(payload[offset+2:])
			return 4, nil
		case serial:
			if available < 10 {
				return 0, fmt.Errorf("payload too small")
			}
			s := ""
			for i := 0; i < 8; i++ {
				s = fmt.Sprintf("%s%02x", s, payload[offset+2+i]&0xff)
			}
			r.UID = s
			return 10, nil
		default:
			return 0, fmt.Errorf("unknown type")
		}
	case totalChannel:
		if available < 6 {
			return 0, fmt.Errorf("payload too small")
		}
		r.Current.Total = float32(binary.LittleEndian.Uint32(payload[offset+2:offset+6])) / 100
		return 6, nil
	case currentChannel:
		if available < 4 {
			return 0, fmt.Errorf("payload too small")
		}
		r.Current.Value = float32(binary.LittleEndian.Uint16(payload[offset+2:offset+4])) / 100
		return 4, nil
	case alarmChannel:
		if available < 9 {
			return 0, fmt.Errorf("payload too small")
		}
		r.Current.Max = float32(binary.LittleEndian.Uint16(payload[offset+2:offset+4])) / 100
		r.Current.Min = float32(binary.LittleEndian.Uint16(payload[offset+4:offset+6])) / 100
		r.Current.Value = float32(binary.LittleEndian.Uint16(payload[offset+6:offset+8])) / 100
		r.Current.Alarms = readAlarm(payload[offset+8])
		return 9, nil
	default:
		return 0, fmt.Errorf("unknown channel")
	}
}

func readVersion(b byte) string {
	major := (b & 0xf0) >> 4
	minor := b & 0x0f
	return fmt.Sprintf("%d.%d", major, minor)
}

func readVersionLong(b []byte) string {
	return fmt.Sprintf("%d.%d", b[0], b[1])
}

func readAlarm(b byte) Alarms {
	a := Alarms{}
	if ((b >> 0) & 0x01) == 1 {
		//alarm = append(alarm, "threshold alarm");
		a.t = true
	}
	if ((b >> 1) & 0x01) == 1 {
		//alarm = append(alarm, "threshold alarm release");
		a.tr = true
	}
	if ((b >> 2) & 0x01) == 1 {
		//alarm = append(alarm, "over range alarm");
		a.r = true
	}
	if ((b >> 3) & 0x01) == 1 {
		//alarm = append(alarm, "over range alarm release");
		a.rr = true
	}
	return a
}

func BytesToFloat32(b []byte) float32 {
	beI := binary.BigEndian.Uint32(b)

	return math.Float32frombits(beI)
}
