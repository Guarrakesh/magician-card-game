package player

import (
	"fmt"
	"wizard/card"
	"wizard/deck"
)

type Player struct {
	Name  string
	Score int
	Hand  []card.Card
	IsAI  bool
}

type Players []*Player

func (p *Player) DrawCards(deck *deck.Deck, numberOfCards int) {
	draws := (*deck).Draw(numberOfCards)
	p.Hand = draws
}

func (p *Player) ShowHand() {
	for _, card := range p.Hand {
		card.Show()
		fmt.Printf(", ")
	}
}

func Register(playerNames []string, numberOfAI int) Players {
	numberOfPlayers := len(playerNames) + numberOfAI
	players := make(Players, numberOfPlayers)

	for index := range numberOfPlayers {
		if index < len(playerNames) {
			players[index] = &Player{
				Name:  playerNames[index],
				Score: 0,
				IsAI:  false,
			}
		} else {
			players[index] = &Player{
				Name:  fmt.Sprintf("Computer %d", index-len(playerNames)+1),
				Score: 0,
				IsAI:  true,
			}
		}

	}

	return players

}
