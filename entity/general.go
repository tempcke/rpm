package entity

type Phone struct {
	Number string // format determined by user input
	Note   string // any user entered notes
}

func NewPhone(number string) Phone {
	return Phone{
		Number: number,
	}
}
func (p Phone) WithNotes(notes string) Phone {
	p.Note = notes
	return p
}
func (p Phone) AppendNote(note string) Phone {
	if p.Note != "" {
		return p.WithNotes(p.Note + "\n" + note)
	}
	return p.WithNotes(note)
}
func (p Phone) Equal(p2 Phone) bool {
	return p.Number == p2.Number && p.Note == p2.Note
}
