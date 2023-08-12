package entity

import (
	"time"
)

type Phone struct {
	Number  string // format determined by user input
	AddedAt time.Time
	Notes   string // any user entered notes
}

func NewPhone(number string) Phone {
	return Phone{
		Number:  number,
		AddedAt: time.Now(),
	}
}
func (p Phone) WithNotes(notes string) Phone {
	p.Notes = notes
	return p
}
func (p Phone) AppendNote(note string) Phone {
	if p.Notes != "" {
		return p.WithNotes(p.Notes + "\n" + note)
	}
	return p.WithNotes(note)
}
