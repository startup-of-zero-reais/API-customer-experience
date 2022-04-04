package data

type (
	Event string

	UserEvent struct {
		Key       string `diinamo:"type:string;hash"`
		EventID   int    `diinamo:"type:int;range"`
		EventType Event  `diinamo:"type:string;gsi:EventTypeIndex;keyPairs:EventType=EventID"`
		Data      []byte `diinamo:"type:string"`
		Timestamp int64  `diinamo:"type:int;gsi:TimestampIndex;keyPairs:Key=Timestamp"`
	}

	UserModel struct {
		ID        string `diinamo:"type:string;hash" json:"id"`
		Name      string `diinamo:"type:string" json:"name"`
		Lastname  string `diinamo:"type:string" json:"lastname"`
		Email     string `diinamo:"type:string;range;gsi:EmailIndex;keyPairs:Email=ID" json:"email"`
		Phone     string `diinamo:"type:string" json:"phone"`
		Password  string `diinamo:"type:string" json:"-"`
		Avatar    string `diinamo:"type:string" json:"avatar,omitempty"`
		Addresses string `diinamo:"type:string" json:"addresses,omitempty"`
	}
)