package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// These are predefined layouts for use in Time.Format and Time.Parse.
//
// The RFC3339Milli format always keep zeros milliseconds.
const (
	RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

	JSONDate      = "2006-01-02"
	JSONShortDate = "20060102"

	JSONMonth      = "2006-01"
	JSONShortMonth = "200601"
)

// UnixZero is time at January 1, 1970 (midnight UTC/GMT).
var UnixZero = Time{Time: time.Unix(0, 0)}

// Getter provides function to get time.Time object.
type Getter interface {
	GetTime() time.Time
}

// Date represents an instant in date by embeds a time.Time.
// Date use time.Time to store a date in UTC, so before get
// year, month and day, Date.Time must be converted to UTC timezone.
//
// Date can be marshal to text and unmarshal from text in format
// JSONDate in JSON, path param, query param and plain text.
//
// Date implements TextMarshaler and TextUnmarshaler to encode to text
// and decode from text in format "2016-01-02".
//
// Date implements json.Marshaler and json.Unmarshaler to encode to
// JSON string and decode from JSON string in format "2016-01-02"
//
// Date implements json.BindUnmarshaler to decode from query and
// path param in format "2016-01-02"
//
// Date implements bson.Getter and bson.Setter to encode and decode to/from
// mongodb's Date type.
//
// Programs using dates should typically store and pass them as values,
// not pointers. That is, time variables and struct fields should be of
// type types.Date, not *types.Date.
type Date struct {
	time.Time
}

var _ Getter = Date{}
var _ json.Marshaler = Date{}
var _ json.Unmarshaler = &Date{}

// DateString parses string in format "2016-01-02" and returns a new Date instance.
func DateString(str string) (Date, error) {
	d := Date{}
	t, err := time.Parse(JSONDate, str)
	if err != nil {
		return d, err
	}

	d.Time = t
	return d, err
}

// DateShortString parses string in format "20160102" and returns a new Date instance.
func DateShortString(str string) (Date, error) {
	d := Date{}
	t, err := time.Parse(JSONShortDate, str)
	if err != nil {
		return d, err
	}

	d.Time = t
	return d, err
}

// DateLocation creates a new date from time and location.
func DateLocation(t time.Time, loc *time.Location) Date {
	y, m, d := t.In(loc).Date()
	date := Date{
		Time: time.Date(y, m, d, 0, 0, 0, 0, time.UTC),
	}

	return date
}

// DateTimezoneOffset creates a new date from time and timezone offset.
func DateTimezoneOffset(t time.Time, tzOffset int) Date {
	t = t.Add(time.Duration(tzOffset) * time.Second)
	y, m, d := t.In(time.UTC).Date()
	date := Date{
		Time: time.Date(y, m, d, 0, 0, 0, 0, time.UTC),
	}

	return date
}

// MarshalJSON encodes Date to JSON string in JSONDate format.
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// UnmarshalJSON decodes a JSON string in format JSONDate to
// a Date.
func (d *Date) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(strconv.Quote(JSONDate), string(b))
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The date is formatted in JSONDate format.
func (d Date) MarshalText() ([]byte, error) {
	d.Time.MarshalText()
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The date is expected to be in JSONDate format.
func (d *Date) UnmarshalText(data []byte) error {
	t, err := time.Parse(JSONDate, string(data))
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

// Value returns a driver Value in time.Time type.
func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}

// Scan assigns a value in time.Time type from a database driver.
func (d *Date) Scan(value interface{}) error {
	*d = Date{Time: value.(time.Time)}
	return nil
}

// GetTime returns time.Time value that is embedded in.
func (d Date) GetTime() time.Time {
	return d.Time
}

// String returns a string representing the date in JSONDate format
func (d Date) String() string {
	return d.In(time.UTC).Format(JSONDate)
}

// ShortString returns a string representing the date in JSONShortDate format.
func (d Date) ShortString() string {
	return d.In(time.UTC).Format(JSONShortDate)
}

// Equal reports whether d and u represent the same date instant.
// Two dates can be equal even if they are in different locations.
func (d Date) Equal(u Date) bool {
	dY, dM, dD := d.Date()
	uT, uM, uD := u.Date()

	return dY == uT && dM == uM && dD == uD
}

// Before reports whether the date instant d is before u.
func (d Date) Before(u Date) bool {
	dY, dM, dD := d.Date()
	uT, uM, uD := u.Date()

	return dY < uT || dY == uT && dM < uM || dY == uT && dM == uM && dD < uD
}

// After reports whether the date instant d is after u.
func (d Date) After(u Date) bool {
	dY, dM, dD := d.Date()
	uT, uM, uD := u.Date()

	return dY > uT || dY == uT && dM > uM || dY == uT && dM == uM && dD > uD
}

// AddDate returns the Date corresponding to adding the
// given number of years, months, and days to t.
// For example, AddDate(-1, 2, 3) applied to January 1, 2011
// returns March 4, 2010.
//
// AddDate normalizes its result in the same way that Date does,
// so, for example, adding one month to October 31 yields
// December 1, the normalized form for November 31.
func (d Date) AddDate(years int, months int, days int) Date {
	return Date{Time: d.Time.AddDate(years, months, days)}
}

// IsZero reports whether t represents the zero date instant,
// January 1, year 1.
func (d Date) IsZero() bool {
	return d.Time.IsZero()
}

// A Time represents an instant in time with nanosecond precision.
//
// Programs using times should typically store and pass them as values,
// not pointers. That is, time variables and struct fields should be of
// type types.Time, not *types.Time.
type Time struct {
	time.Time
}

var _ Getter = Time{}
var _ json.Marshaler = Time{}
var _ json.Unmarshaler = &Time{}

// MarshalJSON encodes Time to JSON DateTime in format RFC3339Milli.
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.Time.UTC().Format(RFC3339Milli))), nil
}

// UnmarshalJSON decodes JSON DateTime to Time.
func (t *Time) UnmarshalJSON(b []byte) error {
	gt := time.Time{}
	err := json.Unmarshal(b, &gt)
	if err != nil {
		return err
	}
	t.Time = gt

	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted in RFC 3339 Milli format.
func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in RFC 3339 Milli format.
func (t *Time) UnmarshalText(data []byte) error {
	gt, err := time.Parse(RFC3339Milli, string(data))
	if err != nil {
		return err
	}

	t.Time = gt
	return nil
}

// String returns a string representing the time in RFC3339Milli format.
func (t Time) String() string {
	return fmt.Sprintf("%s", t.Time.UTC().Format(RFC3339Milli))
}

// Value returns a driver Value in time.Time type.
func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}

// Scan assigns a value in time.Time type from a database driver.
func (t *Time) Scan(value interface{}) error {
	*t = Time{Time: value.(time.Time)}
	return nil
}

// GetTime returns time.Time value that is embedded in.
func (t Time) GetTime() time.Time {
	return t.Time
}

// Equal reports whether ts and u represent the same time instant.
// Two times can be equal even if they are in different locations.
func (t Time) Equal(u Time) bool {
	return t.Time.Equal(u.Time)
}

// Before reports whether the time instant t is before u.
func (t Time) Before(u Time) bool {
	return t.Time.Before(u.Time)
}

// After reports whether the time instant t is after u.
func (t Time) After(u Time) bool {
	return t.Time.After(u.Time)
}

// Add returns the time t+d.
func (t Time) Add(d time.Duration) Time {
	return Time{Time: t.Time.Add(d)}
}

// AddDate returns the time corresponding to adding the
// given number of years, months, and days to t.
// For example, AddDate(-1, 2, 3) applied to January 1, 2011
// returns March 4, 2010.
//
// AddDate normalizes its result in the same way that Date does,
// so, for example, adding one month to October 31 yields
// December 1, the normalized form for November 31.
func (t Time) AddDate(years int, months int, days int) Time {
	return Time{Time: t.Time.AddDate(years, months, days)}
}

// IsZero reports whether t represents the zero date instant,
// January 1, year 1, 00:00:00 UTC.
func (t Time) IsZero() bool {
	return t.Time.IsZero()
}

// Month represents an instant in date by embeds a time.Time.
// Month use time.Time to store a date in UTC, so before get
// year, month and day, Month.Time must be converted to UTC timezone.
//
// Month can be marshal to text and unmarshal from text in format
// "2006-01" in JSON, path param, query param and plain text.
//
// Month implements TextMarshaler and TextUnmarshaler to encode to text
// and decode from text in format "2016-01".
//
// Month implements json.Marshaler and json.Unmarshaler to encode to
// JSON string and decode from JSON string in format "2016-01"
//
// Month implements json.BindUnmarshaler to decode from query and
// path param in format "2016-01"
//
// Month implements bson.Getter and bson.Setter to encode and decode to/from
// mongodb's Date type.
//
// Programs using dates should typically store and pass them as values,
// not pointers. That is, time variables and struct fields should be of
// type types.Month, not *types.Month.
type Month struct {
	time.Time
}

var _ Getter = Month{}
var _ json.Marshaler = Month{}
var _ json.Unmarshaler = &Month{}

// NewMonth creates and returns the Month
func NewMonth(year int, month time.Month) Month {
	return Month{Time: time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)}
}

// MarshalJSON encodes Month to JSON string in JSONMonth format.
func (m Month) MarshalJSON() ([]byte, error) {
	time.Unix(0, 0)
	return []byte(`"` + m.String() + `"`), nil
}

// UnmarshalJSON decodes a JSON string in JSONMonth format to
// a Month.
func (m *Month) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(strconv.Quote(JSONMonth), string(b))
	if err != nil {
		return err
	}

	m.Time = t
	return nil
}

// String returns a string representing the date in JSONMonth format.
func (m Month) String() string {
	return m.In(time.UTC).Format(JSONMonth)
}

// ShortString returns a string representing the date in  JSONShortMonth format.
func (m Month) ShortString() string {
	return m.In(time.UTC).Format(JSONShortMonth)
}

// MarshalText implements the encoding.TextMarshaler interface.
// The month is formatted in JSONMonth format.
func (m Month) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The month is expected to be in JSONMonth format.
func (m *Month) UnmarshalText(data []byte) error {
	t, err := time.Parse(strconv.Quote(JSONMonth), string(data))
	if err != nil {
		return err
	}

	m.Time = t
	return nil
}

// Value returns a driver Value in time.Time type.
func (m Month) Value() (driver.Value, error) {
	return m.Time, nil
}

// Scan assigns a value in time.Time type from a database driver.
func (m *Month) Scan(value interface{}) error {
	*m = Month{Time: value.(time.Time)}
	return nil
}

// GetTime returns time.Time value that is embedded in.
func (m Month) GetTime() time.Time {
	return m.Time
}

// Equal reports whether ts and u represent the same month instant.
// Two months can be equal even if they are in different locations.
func (m Month) Equal(u Month) bool {
	mY, mM, _ := m.Date()
	uY, uM, _ := u.Date()

	return mY == uY && mM == uM
}

// Before reports whether the time instant t is before u.
func (m Month) Before(u Month) bool {
	mY, mM, _ := m.Date()
	uY, uM, _ := u.Date()

	return mY < uY || mY == uY && mM < uM
}

// After reports whether the time instant t is after u.
func (m Month) After(u Month) bool {
	mY, mM, _ := m.Date()
	uY, uM, _ := u.Date()

	return mY > uY || mY == uY && mM > uM
}

// IsZero reports whether t represents the zero date instant,
// January, year 1.
func (m Month) IsZero() bool {
	return m.Time.IsZero()
}

// DateBSONString represents an instant in date by embeds a time.Time.
// DateBSONString use time.Time to store a date in UTC, so before get
// year, month and day, DateBSONString.Time must be converted to UTC timezone.
//
// DateBSONString can be marshal to text and unmarshal from text in format
// JSONDate in JSON, path param, query param and plain text.
//
// DateBSONString implements TextMarshaler and TextUnmarshaler to encode to
// text and decode from text in format "2016-01-02".
//
// DateBSONString implements json.Marshaler and json.Unmarshaler to encode to
// JSON string and decode from JSON string in format "2016-01-02"
//
// DateBSONString implements json.BindUnmarshaler to decode from query and
// path param in format "2016-01-02"
//
// DateBSONString implements bson.Getter and bson.Setter to encode and decode
// to/from mongodb's string.
//
// Programs using dates should typically store and pass them as values,
// not pointers. That is, time variables and struct fields should be of
// type types.DateBSONString, not *types.DateBSONString.
type DateBSONString struct {
	time.Time
}

var _ Getter = DateBSONString{}
var _ json.Marshaler = DateBSONString{}
var _ json.Unmarshaler = &DateBSONString{}

// DateBSONStringLocation creates a new DateBSONString from time and location.
func DateBSONStringLocation(t time.Time, loc *time.Location) DateBSONString {
	date := DateBSONString{}
	y, m, d := t.In(loc).Date()
	date.Time = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return date
}

// MarshalJSON encodes Date to JSON string in format JSONDate
func (d DateBSONString) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// UnmarshalJSON decodes a JSON string in format JSONDate to
// a Date.
func (d *DateBSONString) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(strconv.Quote(JSONDate), string(b))
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

// String returns a string representing the date in format "yyyy-mm-dd"
func (d DateBSONString) String() string {
	return d.In(time.UTC).Format(JSONDate)
}

// MarshalText implements the encoding.TextMarshaler interface.
// The date is formatted in JSONDate format.
func (d DateBSONString) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The date is expected to be in JSONDate format.
func (d *DateBSONString) UnmarshalText(data []byte) error {
	t, err := time.Parse(JSONDate, string(data))
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

// GetTime returns time.Time value that is embedded in.
func (d DateBSONString) GetTime() time.Time {
	return d.Time
}

// Equal reports whether d and u represent the same date instant.
// Two dates can be equal even if they are in different locations.
func (d DateBSONString) Equal(u DateBSONString) bool {
	dY, dM, dD := d.Date()
	uT, uM, uD := u.Date()

	return dY == uT && dM == uM && dD == uD
}

// Before reports whether the date instant d is before u.
func (d DateBSONString) Before(u DateBSONString) bool {
	dY, dM, dD := d.Date()
	uT, uM, uD := u.Date()

	return dY < uT || dY == uT && dM < uM || dY == uT && dM == uM && dD < uD
}

// After reports whether the date instant d is after u.
func (d DateBSONString) After(u DateBSONString) bool {
	dY, dM, dD := d.Date()
	uT, uM, uD := u.Date()

	return dY > uT || dY == uT && dM > uM || dY == uT && dM == uM && dD > uD
}

// AddDate returns the Date corresponding to adding the
// given number of years, months, and days to t.
// For example, AddDate(-1, 2, 3) applied to January 1, 2011
// returns March 4, 2010.
//
// AddDate normalizes its result in the same way that Date does,
// so, for example, adding one month to October 31 yields
// December 1, the normalized form for November 31.
func (d DateBSONString) AddDate(years int, months int, days int) DateBSONString {
	return DateBSONString{Time: d.Time.AddDate(years, months, days)}
}

// IsZero reports whether t represents the zero date instant,
// January 1, year 1.
func (d DateBSONString) IsZero() bool {
	return d.Time.IsZero()
}

// Now returns the current local time.
func Now() Time {
	return Time{Time: time.Now()}
}

// Parse parses a formatted string and returns the time value it represents.
// The layout defines the format by showing how the reference time,
// defined to be
//	Mon Jan 2 15:04:05 -0700 MST 2006
// would be interpreted if it were the value; it serves as an example of
// the input format. The same interpretation will then be made to the
// input string.
//
// Predefined layouts ANSIC, UnixDate, RFC3339 and others describe standard
// and convenient representations of the reference time. For more information
// about the formats and the definition of the reference time, see the
// documentation for ANSIC and the other constants defined by this package.
// Also, the executable example for Time.Format demonstrates the working
// of the layout string in detail and is a good reference.
//
// Elements omitted from the value are assumed to be zero or, when
// zero is impossible, one, so parsing "3:04pm" returns the time
// corresponding to Jan 1, year 0, 15:04:00 UTC (note that because the year is
// 0, this time is before the zero Time).
// Years must be in the range 0000..9999. The day of the week is checked
// for syntax but it is otherwise ignored.
//
// In the absence of a time zone indicator, Parse returns a time in UTC.
//
// When parsing a time with a zone offset like -0700, if the offset corresponds
// to a time zone used by the current location (Local), then Parse uses that
// location and zone in the returned time. Otherwise it records the time as
// being in a fabricated location with time fixed at the given zone offset.
//
// When parsing a time with a zone abbreviation like MST, if the zone abbreviation
// has a defined offset in the current location, then that offset is used.
// The zone abbreviation "UTC" is recognized as UTC regardless of location.
// If the zone abbreviation is unknown, Parse records the time as being
// in a fabricated location with the given zone abbreviation and a zero offset.
// This choice means that such a time can be parsed and reformatted with the
// same layout losslessly, but the exact instant used in the representation will
// differ by the actual zone offset. To avoid such problems, prefer time layouts
// that use a numeric zone offset, or use ParseInLocation.
func Parse(layout, value string) (Time, error) {
	stdtime, err := time.Parse(layout, value)
	if err != nil {
		return Time{}, err
	}

	return Time{Time: stdtime}, nil
}

// MustParse likes Parse but panics if an error occurs.
func MustParse(layout, value string) Time {
	t, err := Parse(layout, value)
	if err != nil {
		panic(err)
	}

	return t
}
