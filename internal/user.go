package internal

type UserLogPas struct {
	// Must be unique in the system
	Login string `json:"login,omitempty"`
	// Can't be empty and will not change
	Password string `json:"password,omitempty"`
}

type UserInfo struct {
	// The user's name
	Name string `json:"name,omitempty"`
	// The user's surname
	Surname string `json:"surname,omitempty"`
	// The user's date of birth
	DateOfBirth string `json:"date_of_birth,omitempty"`
	// The user's mail
	Mail string `json:"mail,omitempty"`
	// The user's phone number
	PhoneNumber string `json:"phone_number,omitempty"`
}
