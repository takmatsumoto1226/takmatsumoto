package df

import "fmt"

type LMap map[string]int

func (m LMap) Presentation() string {
	msg := ""
	for k, v := range m {
		msg = msg + fmt.Sprintf("%s:%d\n", k, v)
	}
	return msg
}
