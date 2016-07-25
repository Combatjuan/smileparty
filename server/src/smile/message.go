package smile

import (
	"fmt"
)

type SmileLocation struct {
	Id string `json:"id"`
	X   int `json:"x"`
	Y   int `json:"x"`
}

func (self *SmileLocation) String() string {
	return fmt.Sprintf("Smile %d (%d, %d)", self.Id, self.X, self.Y)
}
