// In A Game round, each apprentice throws a card
// - A "Trump" is placed at the beginning of the turn
// - The first apprentice plays a Suit
// - Other apprentices must follow the symbol of the Suit
// - if an Apprentice doesn't have any card of the Suit, he/she can play the Trump or any card,
// - The Trick is granted to the apprentice that placed the highest Suit or the Highest trump
// - When both Trump's symbol and Suit's symbol cards are on the field, the highest trump always wins.

// Special Cases:
// - If the first apprentices plays a Wizard or a Joker as Suit, the other players can play any card.
// - At the final turn, where all cards are distributed and there is no Trump, the first apprentices
// plays the Suit and the highest Suit's Symbol card will win.
package game

import (
	"wizard/card"
	"wizard/player"
)

type Round struct {
	Trump   *card.Card // The card been placed at the beginning of the turn
	Suit    *card.Card // The card the first apprentice plays in this turn
	Highest *card.Card
	Tricker *player.Player // The apprentice who made a trick (ie won the round)
}
