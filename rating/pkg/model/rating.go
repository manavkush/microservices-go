package model

// RecordID defines a record id. Together with RecordType
// identifies unique records across all types.
type RecordID string

// RecordType defines a record type. Together with RecordID
// identifies unique records across all types.
type RecordType string

// Existing record types
const (
	RecordMovieType = RecordType("movie")
)

// UserID defines a user id
type UserID string

// RatingValue defines a value of rating record.
type RatingValue int

// Rating defines an individual rating created by a user for
// some record
type Rating struct {
	RecordID   RecordID    `json:"recordId"`
	RecordType RecordType  `json:"recordType"`
	UserID     UserID      `json:"userId"`
	Value      RatingValue `json:"ratingValue"`
}
