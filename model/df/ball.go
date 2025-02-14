package df

// Balls ...
type IBalls []IBall // base array class

// Ball base class
type IBall struct {
	Number   string `json:"number"`
	Position int    `json:"position"`
	Digit    int    `json:"digit"`
	Period   int    `json:"period"`
	Continue int    `json:"continue"`
}
