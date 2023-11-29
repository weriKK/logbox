package common

type LogMessage struct {
	Id        int    `json:"i"`
	Timestamp int    `json:"t"`
	Message   string `json:"m"`
}
