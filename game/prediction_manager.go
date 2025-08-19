package game

import (
	"fmt"
	"wizard/game/internals"
	"wizard/player"
)

func AskPredictions(turn *Turn, dealerPos int, players player.Players) {
	predictions := make([]Prediction, len(players))

	for i := 0; i < len(players); i++ {
		currPlayer := players[internals.GetPlayerPos(i, dealerPos, len(players))]

		if !currPlayer.IsAI {
			fmt.Printf("\n%s, here's your hand:\n", currPlayer.Name)
			currPlayer.ShowHand()
			fmt.Print("\nMake a prediction of tricks: ")
			var prediction int
			fmt.Scan(&prediction)
			predictions[i] = Prediction{Player: currPlayer, PredictedTricks: prediction}
		} else {

		}

	}

	turn.predictions = predictions

}
