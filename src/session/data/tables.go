package data

type (
	UserSession struct {
		UserID       string `diinamo:"type:string;hash"`
		SessionID    string `diinamo:"type:string;gsi:SessionIndex;keyPairs:SessionID=CreatedAt"`
		CreatedAt    int64  `diinamo:"type:int;range"`
		ExpiresIn    int64  `diinamo:"type:int;gsi:ExpiresInIndex;keyPairs:UserID=ExpiresIn"`
		SessionToken string `diinamo:"type:string"`
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

	PassTokens struct {
		Email     string `diinamo:"type:string;hash"`
		OTP       int    `diinamo:"type:int;range;gsi:OtpIndex;keyPairs:OTP=ExpiresIn"`
		ExpiresIn int64  `diinamo:"type:int;"`
	}
)
