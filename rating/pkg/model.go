package model

// Record ID define a record id. Together with RecordType
// identifes unique record across all types

type RecordID string

type RecordType string

// Existing record types.
const (
	RecordTypeMovie = RecordType("movie")
)

// Validation map
var validRecordTypes = map[RecordType]bool{
	RecordTypeMovie: true,
}

// Validate a record type exists
func IsValidRecordType(rt RecordType) bool {
	_, exists := validRecordTypes[rt]
	return exists
}

//UserId defines  a user id.
type UserID string

// RatingValue  define a value of a rating record.
type RatingValue int

// Rating defnies a individual rating created by a user
// for some record

type Rating struct {
	RecordID   RecordID    `json: "recordId"`
	RecordType RecordType  `json: "recordType"`
	UserID     UserID      `json: UserId`
	Value      RatingValue `json: value`
}

//RatingEvent define a event containing a rating information.
type RatingEvent struct {
	UserID     string          `json: "userId"`
	RecordID   RecordID        `json: "recordId"`
	RecordType RecordType      `json: "recordType"`
	Value      RatingValue     `json:"value"`
	EventType  RatingEventType `json:"eventType"`
}

//RatingEventType defines a the type of a rating event.

type RatingEventType string

// Rating event types.
const (
	RatingEventTypePut    = "put"
	RatingEventTypeDelete = "delete"
)
