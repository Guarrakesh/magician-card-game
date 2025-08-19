package deck

import (
	"fmt"
	"math/rand"
	"time"
	"wizard/card"
)

type Deck []card.Card

func InitDeck() Deck {
	var deck = make([]card.Card, 60)

	// 52 number cards
	// 4 jokers
	// 4 wizards

	for s := 0; s < 4; s++ {
		for n := 1; n <= 13; n++ {
			i := s*13 + (n - 1)
			deck[i].Number = n
			deck[i].Symbol = card.Symbols[s]
		}
	}

	for i := 0; i < 4; i++ {
		deck[52+i].IsJoker = true
		deck[52+i].Number = -1
		deck[52+i].Symbol = ""
	}
	for i := 0; i < 4; i++ {
		deck[56+i].IsWizard = true
		deck[56+i].Number = -1
		deck[56+i].Symbol = ""
	}

	return deck
}

func (deck Deck) Show() {
	for i := 0; i < len(deck); i++ {
		deck[i].Show()

		if i < len(deck)-1 {
			fmt.Printf(",")
		}
	}
}

func (deck *Deck) Shuffle() {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := len(*deck) - 1; i > 0; i-- {
		j := seed.Intn(i + 1)
		(*deck)[i], (*deck)[j] = (*deck)[j], (*deck)[i]
	}
}

func (deck *Deck) Draw(numberOfCards int) []card.Card {
	var draw = (*deck)[len(*deck)-numberOfCards : len(*deck)]
	*deck = (*deck)[:len(*deck)-numberOfCards]

	return draw
}

func (deck *Deck) Add(card *card.Card) {
	*deck = append(*deck, *card)
}
