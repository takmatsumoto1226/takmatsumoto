package df

type Frequency struct {
	Title string `json:"title"`
	Type  int    `json:"type"`
	Min   int    `json:"min"`
	Avg   int    `json:"avg"`
	Max   int    `json:"max"`
	Count int    `json:"count"`
	Sum   int    `json:"sum"`
}

func NewFrequency(t string, i int) *Frequency {
	return &Frequency{
		Title: t,
		Type:  i,
		Min:   0,
		Avg:   0,
		Max:   0,
		Count: 0,
		Sum:   0,
	}
}

func (f *Frequency) AddOne() {

}

type FrequencyStatic struct {
	FContinue2  Frequency `json:"fcontinue2"`
	FContinue3  Frequency `json:"fcontinue3"`
	FContinue4  Frequency `json:"fcontinue4"`
	FContinue32 Frequency `json:"fcontinue32"`
	FContinue22 Frequency `json:"fcontinue22"`
}
