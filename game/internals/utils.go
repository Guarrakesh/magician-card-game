package internals

func GetPlayerPos(currPlayerIndex int, dealerPos int, playersLen int) int {
	return (dealerPos + currPlayerIndex + 1) % playersLen
}
