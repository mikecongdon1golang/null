package zero

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Time is a nullable time.Time.
// JSON marshals to the zero value for time.Time if null.
// Considered to be null to SQL if zero.
type Time struct {
	mysql.NullTime
}

// Scan implements Scanner interface.
func (t *Time) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		t.Time = x
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.Time: %v", value, value)
	}
	t.Valid = err == nil
	return err
}

// Value implements the driver Valuer interface.
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// NewTime creates a new Time.
func NewTime(t time.Time, valid bool) Time {
	var ti Time
	ti.Time = t
	ti.Valid = valid
	return ti
}

// TimeFrom creates a new Time that will
// be null if t is the zero value.
func TimeFrom(t time.Time) Time {
	return NewTime(t, !t.IsZero())
}

// TimeFromPtr creates a new Time that will
// be null if t is nil or *t is the zero value.
func TimeFromPtr(t *time.Time) Time {
	if t == nil {
		return NewTime(time.Time{}, false)
	}
	return TimeFrom(*t)
}

// MarshalJSON implements json.Marshaler.
// It will encode the zero value of time.Time
// if this time is invalid.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return (time.Time{}).MarshalJSON()
	}
	return t.Time.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string, object (e.g. pq.NullTime and friends)
// and null input.
func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		var ti time.Time
		if err = ti.UnmarshalJSON(data); err != nil {
			return err
		}
		*t = TimeFrom(ti)
		return nil
	case map[string]interface{}:
		ti, tiOK := x["Time"].(string)
		valid, validOK := x["Valid"].(bool)
		if !tiOK || !validOK {
			return fmt.Errorf(`json: unmarshalling object into Go value of type null.Time requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["Time"], x["Valid"])
		}
		err = t.Time.UnmarshalText([]byte(ti))
		t.Valid = valid
		return err
	case nil:
		t.Valid = false
		return nil
	default:
		return fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Time", reflect.TypeOf(v).Name())
	}
}

func (t Time) MarshalText() ([]byte, error) {
	ti := t.Time
	if !t.Valid {
		ti = time.Time{}
	}
	return ti.MarshalText()
}

func (t *Time) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	if err := t.Time.UnmarshalText(text); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

// String changes this Time's value and
// sets it to be non-null.
func (t *Time) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Time.Format("2006-01-02 15:04:05")
}

// StringDate changes this Time's value and
// sets it to be non-null.
func (t *Time) StringDate() string {
	if t.IsZero() {
		return ""
	}
	return t.Time.Format("2006-01-02")
}

// StringTime changes this Time's value and
// sets it to be non-null.
func (t *Time) StringTime() string {
	if t.IsZero() {
		return ""
	}
	return t.Time.Format("15:04:05")
}

// SetValid changes this Time's value and
// sets it to be non-null.
func (t *Time) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

// Ptr returns a pointer to this Time's value,
// or a nil pointer if this Time is zero.
func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// IsZero returns true for null or zero Times, for potential future omitempty support.
func (t Time) IsZero() bool {
	return !t.Valid || t.Time.IsZero()
}

// Now returns true for null or zero Times, for potential future omitempty support.
func TimeNow() Time {
	return TimeFrom(time.Now())
}

// OverwriteWithIfValid returns nothing. Used for type conversion from sql.Nullstring to zero
func (s *Time) OverwriteWithIfValid(st time.Time, v bool) {
	if v {
		s.Time = st
		s.Valid = v
	}
}
