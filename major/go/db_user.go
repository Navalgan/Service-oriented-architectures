package swagger

type DataBaseUser struct {
	SessionID int32       `json:"sessionID,omitempty"`
	User      User        `json:"user,omitempty"`
	Info      Information `json:"info,omitempty"`
}
