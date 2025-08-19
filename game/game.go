package game

import (
	"fmt"
	"math"
	"wizard/player"
)

type Game struct {
	Players player.Players
}

func InitGame(players player.Players) Game {
	game := Game{Players: players}

	return game
}

func (game *Game) Run(dealerPos int16) {

	numberOfTurns := int(math.Floor(float64(60) / float64(len(game.Players))))

	fmt.Printf("There will be %d turns.\n\n", numberOfTurns)
	// Step 0 - init Turn
	turn := Turn{}

	for i := 0; i < numberOfTurns; i++ {
		// delear goes around

		currDelaerPos := ((int(dealerPos) + i) % len(game.Players))
		fmt.Printf("\n\nCurrent Dealer: %s\n\n", game.Players[currDelaerPos].Name)

		turn.Run(game.Players, i+1, currDelaerPos)
	}

}
