package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib"
	"github.com/safecility/iot/devices/milesightct/process/messages"
)

// TODO: EXAMPLE query - adjust for local device storage
const (
	getDeviceStmt = `SELECT uid as DeviceUID, name as DeviceName, tag as DeviceTag, type as DeviceType,
       		companyID, locationID, power_factor, line_voltage
		FROM device
		JOIN safecility.power_device pd on device.id = pd.deviceId
		WHERE type='power' AND device.uid = ?`
)

// DeviceSql is accessed both directly and by the device Cache, direct access is only for uplinks which show Compliance events
type DeviceSql struct {
	sqlDB          *sql.DB
	getDeviceByUID *sql.Stmt
}

func NewDeviceSql(db *sql.DB) (*DeviceSql, error) {
	sqlDB := &DeviceSql{
		sqlDB: db,
	}
	var err error

	if sqlDB.getDeviceByUID, err = db.Prepare(getDeviceStmt); err != nil {
		return nil, err
	}

	return sqlDB, nil
}

func (db DeviceSql) GetDevice(uid string) (*messages.PowerDevice, error) {
	log.Debug().Str("uid", uid).Msg("getting device from sql")
	row := db.getDeviceByUID.QueryRow(uid)

	serverDevice, err := scanDevice(row)
	if err != nil {
		return nil, err
	}

	return serverDevice, nil
}

func (db DeviceSql) Close() error {
	return db.sqlDB.Close()
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanDevice(s rowScanner) (*messages.PowerDevice, error) {
	var (
		name        sql.NullString
		uid         sql.NullString
		tag         sql.NullString
		dType       sql.NullString
		companyUID  sql.NullInt64
		locationUID sql.NullInt64
		powerFactor sql.NullFloat64
		voltage     sql.NullFloat64
	)

	err := s.Scan(&name, &uid, &tag, &dType, &companyUID, &locationUID, &powerFactor, &voltage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	deviceInfo := messages.PowerDevice{
		Device: lib.Device{
			DeviceUID: uid.String,
			DeviceMeta: &lib.DeviceMeta{
				DeviceName: name.String,
				DeviceTag:  tag.String,
				DeviceType: lib.DeviceType(dType.String),
				Listing: &lib.Listing{
					CompanyUID:  fmt.Sprintf("%d", companyUID.Int64),
					LocationUID: fmt.Sprintf("%d", locationUID.Int64),
				},
			},
		},
		PowerFactor: powerFactor.Float64,
		Voltage:     voltage.Float64,
	}

	return &deviceInfo, nil
}
