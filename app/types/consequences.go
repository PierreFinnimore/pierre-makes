package types

type PoetFormValues struct {
	PoetName string `form:"poetName"`
}

type RoomFormValues struct {
	RoomCode            string `form:"roomCode"`
	PoemCount           int    `form:"poemCount"`
	MinimumLineDistance int    `form:"minimumLineDistance"`
}

type SubmissionFormTwoLineValues struct {
	Line0      string `form:"line0"`
	Line1      string `form:"line1"`
	IsLastLine bool   `form:"lastLine"`
}

type Poet struct {
	PoetID int    `bun:"poet_id,pk,autoincrement"`
	Token  string `bun:"token,notnull"`
	Name   string `bun:"name,notnull"`
}

type Room struct {
	RoomID               int    `bun:"room_id,pk,autoincrement"`
	Code                 string `bun:"code,notnull"`
	LinesPerSubmission   int    `bun:"lines_per_submission,notnull"`
	LinesVisible         int    `bun:"lines_visible,notnull"`
	SecondsPerSubmission int    `bun:"seconds_per_submission,notnull"`
	MinimumLineDistance  int    `bun:"minimum_line_distance,notnull"`
}

type Poem struct {
	PoemID                 int           `bun:"poem_id,pk,autoincrement"`
	ReservedPoetID         *int          `bun:"reserved_poet_id"`
	ReservedPoet           *Poet         `bun:"rel:has-one,join:reserved_poet_id=poet_id"`
	RoomID                 int           `bun:"room_id,notnull"` // Non-pointer int to ensure RoomID is not null
	Room                   *Room         `bun:"rel:has-one,join:room_id=room_id"`
	ReservedUntilTimestamp int64         `bun:"reserved_until_timestamp"`
	IsComplete             bool          `bun:"is_complete,default:false"`
	Submissions            []*Submission `bun:"rel:has-many,join:poem_id=poem_id"`
}

type Submission struct {
	SubmissionID int     `bun:"submission_id,pk,autoincrement"` // Primary key with auto-increment
	PoemID       int     `bun:"poem_id,notnull"`
	Poem         *Poem   `bun:"rel:belongs-to,join:poem_id=poem_id"`
	PoetID       int     `bun:"poet_id,notnull"`                     // Foreign key referencing the Poet
	Poet         *Poet   `bun:"rel:belongs-to,join:poet_id=poet_id"` // Relationship to Poet
	Position     int     `bun:"position,notnull"`                    // Position field
	Lines        []*Line `bun:"rel:has-many,join:submission_id=submission_id"`
}

type Line struct {
	LineID       int         `bun:"line_id,pk,autoincrement"`                        // Primary key with auto-increment
	SubmissionID int         `bun:"submission_id,notnull"`                           // Foreign key referencing the Submission
	Submission   *Submission `bun:"rel:belongs-to,join:submission_id=submission_id"` // Relationship to Submission
	Position     int         `bun:"position,notnull"`                                // Position field
	Text         string      `bun:"text,notnull"`                                    // Text of the line
}
