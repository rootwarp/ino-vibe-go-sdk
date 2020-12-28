package device

import "time"

// StatusLog defines structure of log of inclination.
type StatusLog struct {
	Devid             string    `datastore:"Devid"`
	Time              time.Time `datastore:"Time"`
	Temperature       int       `datastore:"Temperature"`
	Battery           int       `datastore:"Battery"`
	RSSI              int       `datastore:"RSSI"`
	InstallSessionKey string    `datastore:"InstallSessionKey"`
}

// InclinationLog defines structures of log of inclination.
type InclinationLog struct {
	Devid             string    `datastore:"devid"`
	Time              time.Time `datastore:"time_created"`
	AccXMg            float64   `datastore:"acc_x_mg"`
	AccYMg            float64   `datastore:"acc_y_mg"`
	AccZMg            float64   `datastore:"acc_z_mg"`
	AngleZ            float64   `datastore:"angle_z"`
	InstallSessionKey string    `datastore:"install_session_key"`
}
