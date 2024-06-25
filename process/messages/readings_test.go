package messages

import (
	"encoding/base64"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"testing"
)

func TestReadMilesiteCT(t *testing.T) {
	ttnData, err := base64.StdEncoding.DecodeString("A5cNAQAABJipAA==")
	if err != nil {
		log.Error().Err(err).Msg("Error decoding base64")
	}
	st := base64.StdEncoding.EncodeToString([]byte("A5c1FwAABJjWAA=="))
	log.Info().Str("encoding", string(st)).Msg("Test")
	log.Debug().Str("data", fmt.Sprintf("%+v", ttnData)).Msg("webhook data")
	whData, err := base64.StdEncoding.DecodeString("A5c1FwAABJjWAA==")
	if err != nil {
		t.Errorf("could not decode base64 data")
		return
	}
	log.Debug().Str("data", fmt.Sprintf("%+v", whData)).Msg("webhook data")
	type args struct {
		payload []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *MilesightCTReading
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "read all",
			args: args{
				payload: []byte{0xFF, 0x0B, 0xFF,
					0xFF, 0x01, 0x01,
					0xFF, 0x16, 0x67, 0x46, 0xD3, 0x88, 0x02, 0x58, 0x00, 0x00,
					0xFF, 0x09, 0x01, 0x00,
					0xFF, 0x0A, 0x01, 0x01,
					0x03, 0x97, 0x10, 0x27, 0x00, 0x00,
					0x84, 0x98, 0xB8, 0x0B, 0xD0, 0x07, 0xC4, 0x09, 0x05},
			},
			want: &MilesightCTReading{
				UID:   "6746d38802580000",
				Power: true,
				Version: Version{
					Ipso:     "0.1",
					Hardware: "1.0",
					Firmware: "1.1",
				},
				Current: Current{
					Total: 100,
					Value: 25,
					Max:   30,
					Min:   20,
					Alarms: Alarms{
						t:  true,
						tr: false,
						r:  true,
						rr: false,
					},
				},
			},
		},
		{
			name: "read all",
			args: args{
				payload: whData,
			},
			want: &MilesightCTReading{
				UID:   "6746d38802580000",
				Power: true,
				Version: Version{
					Ipso:     "0.1",
					Hardware: "1.0",
					Firmware: "1.1",
				},
				Current: Current{
					Total: 100,
					Value: 25,
					Max:   30,
					Min:   20,
					Alarms: Alarms{
						t:  true,
						tr: false,
						r:  true,
						rr: false,
					},
				},
			},
		},
		{
			name: "read all",
			args: args{
				payload: ttnData,
			},
			want: &MilesightCTReading{
				UID:   "6746d38802580000",
				Power: true,
				Version: Version{
					Ipso:     "0.1",
					Hardware: "1.0",
					Firmware: "1.1",
				},
				Current: Current{
					Total: 100,
					Value: 25,
					Max:   30,
					Min:   20,
					Alarms: Alarms{
						t:  true,
						tr: false,
						r:  true,
						rr: false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadMilesightCT(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMilesightCT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadMilesightCT() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readAlarm(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want Alarms
	}{
		// TODO: Add test cases.
		{
			name: "read payload",
			args: args{
				b: 0x01,
			},
			want: Alarms{
				t: true,
			},
		},
		{
			name: "read payload",
			args: args{
				b: 0x05,
			},
			want: Alarms{
				t: true,
				r: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readAlarm(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readAlarm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readSlice(t *testing.T) {
	type args struct {
		r       *MilesightCTReading
		payload []byte
		offset  int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "power",
			args: args{
				r:       &MilesightCTReading{},
				payload: []byte{0xFF, 0x0B, 0xFF},
				offset:  0,
			},
			want: 3,
		},
		{
			name: "ipso",
			args: args{
				r:       &MilesightCTReading{},
				payload: []byte{0xFF, 0x01, 0x01},
				offset:  0,
			},
			want: 3,
		},
		{
			name: "serial",
			args: args{
				r:       &MilesightCTReading{},
				payload: []byte{0xFF, 0x16, 0x67, 0x46, 0xD3, 0x88, 0x02, 0x58, 0x00, 0x00},
				offset:  0,
			},
			want: 10,
		},
		{
			name: "hardware",
			args: args{
				r:       &MilesightCTReading{},
				payload: []byte{0xFF, 0x09, 0x01, 0x00},
				offset:  0,
			},
			want: 4,
		},
		{
			name: "firmware",
			args: args{
				r:       &MilesightCTReading{},
				payload: []byte{0xFF, 0x0A, 0x01, 0x01},
				offset:  0,
			},
			want: 4,
		},
		{
			name: "total current",
			args: args{
				r:       &MilesightCTReading{},
				payload: []byte{0x03, 0x97, 0x10, 0x27, 0x00, 0x00},
				offset:  0,
			},
			want: 6,
		},
		{
			name: "current alarm",
			args: args{
				r:       &MilesightCTReading{},
				payload: []byte{0x84, 0x98, 0xB8, 0x0B, 0xD0, 0x07, 0xC4, 0x09, 0x05},
				offset:  0,
			},
			want: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readSlice(tt.args.r, tt.args.payload, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("readSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readVersion(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readVersion(tt.args.b); got != tt.want {
				t.Errorf("readVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
