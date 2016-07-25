package smile

type Smiley struct {
	Id int
	X int
	Y int
}

type SmileyManager struct {
	Updates chan SmileLocation
	Smileys []*Smiley
}

func NewSmileyManager() *SmileyManager {
	return &SmileyManager{

	}
}
