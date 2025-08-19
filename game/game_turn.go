// A game turn has N rounds, with each player starting with N cards
// - At the end of the turn, the players have no more cards
// - A the beginning of the turn, each player makes a guess (prediction) on how many rounds he will win

// Special Cases
// - If the Trump is a Wizard, the Dealer decides what Symbol is going to be the trump for the turn

package game

import (
	"fmt"
	"wizard/card"
	"wizard/deck"
	"wizard/game/internals"
	"wizard/player"
)

type Prediction struct {
	Player          *player.Player
	PredictedTricks int
	Outcome         int
}

type Turn struct {
	rounds      []Round
	predictions []Prediction
	trump       card.Card
}

// isHigherCard compares two cards and returns true if newCard beats currentHighest
/*
  1. Wizards always win (highest priority)
  2. Trump cards beat suit cards and
  off-suit cards
  3. Suit cards beat off-suit cards
  4. Higher numbers win within the same
  category
  5. Jokers never win
*/
func isHigherCard(newCard, currentHighest *card.Card, trumpSymbol, suitSymbol *card.Symbol) bool {
	// Trump cards beat suit cards and off-suit cards
	newIsTrump := newCard.Symbol == *trumpSymbol
	currentIsTrump := currentHighest.Symbol == *trumpSymbol

	//2. Trump cards beat suit cards and off-suit cards
	if newIsTrump && !currentIsTrump {
		return true
	}
	if !newIsTrump && currentIsTrump {
		return false
	}

	// If both are trump or both are not trump, compare by suit and number
	newIsSuit := newCard.Symbol == *suitSymbol
	currentIsSuit := currentHighest.Symbol == *suitSymbol

	// Suit cards beat off-suit cards
	if newIsSuit && !currentIsSuit {
		return true
	}
	if !newIsSuit && currentIsSuit {
		return false
	}

	// If same type (both trump, both suit, or both off-suit), higher number wins
	return newCard.Number > currentHighest.Number
}

func (turn *Turn) Run(players player.Players, numberOfRounds int, dealerPos int) {

	turn.rounds = make([]Round, numberOfRounds)

	turnDeck := deck.InitDeck()
	turnDeck.Shuffle()
	turnDeck.Shuffle()

	// Step 1 - each player gets N cards (N = number of rounds )
	for _, player := range players {
		player.DrawCards(&turnDeck, numberOfRounds)
	}

	// Step 2 - place the trump
	turn.trump = turnDeck.Draw(1)[0]
	fmt.Print("The trump is: ")
	turn.trump.Show()
	fmt.Println()

	// Step 3 - each player makes a prediction
	AskPredictions(turn, dealerPos, players)
	fmt.Println("Here are the predictions:")
	fmt.Printf("\n\n")
	for _, prediction := range turn.predictions {
		//	fmt.Print(prediction.Player.*Name)
		if prediction.Player == nil {
			continue
		}
		fmt.Printf("Player '%s' Predicts: %d tricks\n", (*prediction.Player).Name, prediction.PredictedTricks)
	}

	for i := 0; i < numberOfRounds; i++ {
		field := deck.Deck{}
		round := Round{Trump: &turn.trump}

		for i, _ := range players {
			fmt.Printf("\nStarting round %d", i+1)
			currPlayer := players[internals.GetPlayerPos(i, dealerPos, len(players))]
			var selected int = -1

			for selected == -1 {
				fmt.Printf("\n%s, it's your turn. Here is your hand\n", currPlayer.Name)
				for index, card := range currPlayer.Hand {
					fmt.Printf("[%d] - ", index)
					card.Show()
				}
				fmt.Printf("\n Type the number corresponding to the card you want to play: ")
				fmt.Scan(&selected)
				if selected < 0 || selected > len(currPlayer.Hand)-1 {
					fmt.Printf("\n The selected card '%d' is not valid. Please try again.", selected)
					selected = -1
				}

			}
			fmt.Printf("%s played: ", currPlayer.Name)
			selectedCard := currPlayer.Hand[selected]
			selectedCard.Show()

			if i == 0 || round.Suit == nil {
				round.Highest = &selectedCard
				round.Suit = &selectedCard
				round.Tricker = currPlayer
				fmt.Print("\t\t \n\nThe Suit Is: ")
				round.Suit.Show()
				fmt.Printf("\n\n")
			}

			field.Add(&selectedCard)

			// Handle card comparison and update highest card
			if selectedCard.IsWizard {
				// Wizard always wins
				round.Tricker = currPlayer
				round.Highest = &selectedCard
			} else if selectedCard.IsJoker {
				// Joker never wins, skip comparison
			} else if round.Highest.IsJoker {
				// Any non-joker beats a joker
				round.Tricker = currPlayer
				round.Highest = &selectedCard
			} else if !round.Highest.IsWizard {
				// Only compare if current highest is not a wizard
				if isHigherCard(&selectedCard, round.Highest, &turn.trump.Symbol, &round.Suit.Symbol) {
					round.Tricker = currPlayer
					round.Highest = &selectedCard
				}
			}
		}

		fmt.Printf("\n\nCards on the table: ")
		field.Show()

		fmt.Printf("\n\nThe tricker for this round is: %s with the card ", round.Tricker.Name)
		round.Highest.Show()
		fmt.Printf("\n\n")

	}

}
