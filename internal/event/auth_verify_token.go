package event

// AuthVerifyToken is an event for verifying a token
type AuthVerifyToken struct {
	Token string `json:"token"`
}

// AuthVerifyTokenAck is an acknowledgement for verifying a token
type AuthVerifyTokenAck struct {
	Valid bool   `json:"valid"`
	ID    string `json:"id"`
}
