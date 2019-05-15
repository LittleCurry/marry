package vm

type RegisterRes struct {
	Event  string `json:"event"`
	State  int    `json:"state"`
	Reason string `json:"reason"`
}

type RecData struct {
	Event   string `json:"event"`
	PrintId string `json:"print_id"`
	Type    string `json:"type"`
}

type SentData struct {
	Event   string `json:"event"`
	PrintId string `json:"print_id"`
	Url     string `json:"url"`
}
