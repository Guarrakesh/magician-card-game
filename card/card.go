package card

import (
	"github.com/fatih/color"
)

type Symbol string

const (
	Blue   Symbol = "Human"
	Green  Symbol = "Elf"
	Red    Symbol = "Dwarf"
	Yellow Symbol = "Giant"
)

var Symbols = []Symbol{Blue, Green, Red, Yellow}

type Card struct {
	IsJoker  bool
	IsWizard bool
	Number   int
	Symbol   Symbol
}

func (c Card) Show() {
	if c.IsJoker {
		format := color.New(color.FgBlack)
		format = format.Add(color.BgWhite)
		format.Print("[Joker]")
		return
	}
	if c.IsWizard {
		format := color.New(color.FgHiWhite)
		format = format.Add(color.BgHiMagenta)
		format.Print("[Wizard]")
		return
	}

	format := color.New(color.FgHiWhite)
	switch c.Symbol {
	case Blue:
		format = format.Add(color.BgBlue)
	case Green:
		format = format.Add(color.BgGreen)
	case Yellow:
		format = format.Add(color.BgYellow)
	case Red:
		format = format.Add(color.BgRed)
	}

	format.Printf("[%d, %s]", c.Number, c.Symbol)

}
