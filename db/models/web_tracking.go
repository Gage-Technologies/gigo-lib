package models

import (
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/kisielk/sqlstruct"
)

type WebTrackingEvent string

const (
	PageVisit     WebTrackingEvent = "pagevisit"
	LoginStart    WebTrackingEvent = "loginstart"
	Login         WebTrackingEvent = "login"
	Logout        WebTrackingEvent = "logout"
	SignupStart   WebTrackingEvent = "signup"
	Signup        WebTrackingEvent = "signup"
	ResetPassword WebTrackingEvent = "resetpassword"
)

// WebTracking
//
//	An object representing web tracking information for a user's usage
type WebTracking struct {
	// ID Unique identifier of the usage
	ID int64 `json:"_id" sql:"_id"`

	// UserID The id of the user who made the usage
	UserID *int64 `json:"user_id" sql:"user_id"`

	// IP The ip address of the client making the request
	IP net.IP `json:"ip" sql:"ip"`

	// Host The host where the request originated from
	Host string `json:"host" sql:"host"`

	// Event The event type that triggered this usage
	Event WebTrackingEvent `json:"event" sql:"event"`

	// Timestamp Time at which the usage was recorded
	Timestamp time.Time `json:"timestamp" sql:"timestamp"`

	// TimeSpent The amount of time spent on the page
	TimeSpent *time.Duration `json:"timespent" sql:"timespent"`

	// Path The path of the page visited
	Path string `json:"path" sql:"path"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata" sql:"metadata"`
}

type WebTrackingSQL struct {
	ID        int64            `json:"_id" sql:"_id"`
	UserID    *int64           `json:"user_id" sql:"user_id"`
	IP        int64            `json:"ip" sql:"ip"`
	Host      string           `json:"host" sql:"host"`
	Event     WebTrackingEvent `json:"event" sql:"event"`
	Timestamp time.Time        `json:"timestamp" sql:"timestamp"`
	TimeSpent *time.Duration   `json:"timespent" sql:"timespent"`
	Path      string           `json:"path" sql:"path"`
	Metadata  []byte           `json:"metadata" sql:"metadata"`
}

func CreateWebTracking(_id int64, userId *int64, ip net.IP, host string, event WebTrackingEvent, timestamp time.Time, timespent *time.Duration, path string, metadata map[string]interface{}) *WebTracking {
	return &WebTracking{
		ID:        _id,
		UserID:    userId,
		IP:        ip,
		Host:      host,
		Event:     event,
		Timestamp: timestamp,
		TimeSpent: timespent,
		Path:      path,
		Metadata:  metadata,
	}
}

func WebTrackingFromSqlNative(rows *sql.Rows) (*WebTracking, error) {
	usage := &WebTrackingSQL{}
	err := sqlstruct.Scan(usage, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan usage: %v", err)
	}

	// parse the json from bytes
	var metadata map[string]interface{}
	err = json.Unmarshal(usage.Metadata, &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %v", err)
	}

	wt := &WebTracking{
		ID:        usage.ID,
		UserID:    usage.UserID,
		Host:      usage.Host,
		Event:     usage.Event,
		Timestamp: usage.Timestamp,
		TimeSpent: usage.TimeSpent,
		Path:      usage.Path,
		Metadata:  metadata,
	}

	// convert ip from int64 to net.IP
	binary.BigEndian.PutUint32(wt.IP, uint32(usage.IP))

	return wt, nil
}

func (w *WebTracking) ToSqlNative() ([]SQLInsertStatement, error) {
	// serialize the metadata as JSON
	bytes, err := json.Marshal(w.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %v", err)
	}

	// convert the ip to an integer
	ip := binary.BigEndian.Uint32(w.IP.To4())

	return []SQLInsertStatement{
		{
			Statement: "insert into web_tracking (_id, user_id, ip, host, event, timestamp, timespent, path, metadata) values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			Values:    []interface{}{w.ID, w.UserID, ip, w.Host, w.Event, w.Timestamp, w.TimeSpent, w.Path, bytes},
		},
	}, nil
}
