package data

type (
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
