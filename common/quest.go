package common

type Quest struct {
	Components []Component `json:"components"`
}

type Component struct {
	Names  string `json:"names"`
	Status int    `json:"status"`
}
