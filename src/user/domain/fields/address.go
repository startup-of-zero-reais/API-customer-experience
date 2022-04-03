package fields

import "encoding/json"

type (
	// Address struct represents a address
	Address struct {
		ID           string
		Street       string
		City         string
		State        string
		ZipCode      string
		Neighborhood string
	}

	// AddressRepository is a interface to access the address
	Addresses []Address
)

func NewAddress() *Address {
	return &Address{}
}

func (a *Addresses) String() string {
	bytes, err := json.Marshal(a)

	if err != nil {
		return ""
	}

	return string(bytes)
}
