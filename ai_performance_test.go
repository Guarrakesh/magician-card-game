package main

import (
	"fmt"
	"math"
	"testing"
	"time"
	"wizard/game"
	"wizard/player"
)

// AIPerformanceMetrics tracks AI performance across multiple games
type AIPerformanceMetrics struct {
	PlayerName           string
	GamesPlayed          int
	GamesWon             int
	TotalScore           int
	BidAccuracy          float64 // Percentage of exact bids
	BidTendency          float64 // Positive = overbids, Negative = underbids
	AverageDecisionTime  time.Duration
	TricksWonVsPredicted map[int]int // [predicted] -> actual_won
	SpecialCardUsage     SpecialCardStats
	GameResults          []GameResult
}

// SpecialCardStats tracks Wizard and Joker usage patterns
type SpecialCardStats struct {
	WizardsPlayed     int
	WizardsWonTricks  int
	JokersPlayed      int
	JokersLostTricks  int
	WizardTiming      []float64 // When in round (0.0 = early, 1.0 = late)
	JokerTiming       []float64
}

// GameResult stores outcome of a single game
type GameResult struct {
	FinalScore      int
	FinalRank       int
	BidsSuccessful  int
	BidsTotal       int
	AvgBidError     float64
	GameDuration    time.Duration
}

// AITestSuite contains all AI performance tests
type AITestSuite struct {
	TestGames       int
	DifficultyLevel string
	Opponents       []string // AI difficulty levels to test against
	Metrics         map[string]*AIPerformanceMetrics
}

// TestAIBiddingAccuracy tests the AI's ability to predict tricks accurately
func TestAIBiddingAccuracy(t *testing.T) {
	suite := NewAITestSuite(100, "medium") // 100 test games
	
	// Test scenarios with known card distributions
	testScenarios := []struct {
		name        string
		handSize    int
		trumpSuit   string
		expectedBid int
		tolerance   int
	}{
		{"Strong Wizard Hand", 5, "hearts", 3, 1},
		{"Weak Joker Hand", 5, "spades", 1, 1},
		{"Mixed Trump Hand", 7, "clubs", 4, 2},
		{"No Trump Hand", 3, "diamonds", 1, 1},
	}
	
	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			accuracy := suite.TestBiddingScenario(scenario.handSize, scenario.trumpSuit)
			
			if accuracy < 0.7 { // 70% accuracy threshold
				t.Errorf("AI bidding accuracy too low: %.2f%% for scenario %s", 
					accuracy*100, scenario.name)
			}
		})
	}
}

// TestAICardPlayStrategy tests strategic card play decisions
func TestAICardPlayStrategy(t *testing.T) {
	suite := NewAITestSuite(50, "medium")
	
	strategicTests := []struct {
		name     string
		scenario string
		expected string
	}{
		{"Need More Tricks", "behind_on_bid", "play_winning_cards"},
		{"Avoid Extra Tricks", "ahead_on_bid", "play_losing_cards"},
		{"Exact Bid Met", "exact_bid", "defensive_play"},
		{"Last Round Critical", "final_trick", "optimal_card"},
	}
	
	for _, test := range strategicTests {
		t.Run(test.name, func(t *testing.T) {
			correctDecisions := suite.TestCardPlayStrategy(test.scenario)
			
			if correctDecisions < 0.8 { // 80% correct strategic decisions
				t.Errorf("AI strategic decisions too low: %.2f%% for %s", 
					correctDecisions*100, test.name)
			}
		})
	}
}

// TestAIDifficultyLevels ensures each difficulty performs as expected
func TestAIDifficultyLevels(t *testing.T) {
	difficulties := []string{"easy", "medium", "hard", "expert"}
	expectedWinRates := []float64{0.15, 0.35, 0.55, 0.75} // Against human-level play
	
	for i, difficulty := range difficulties {
		t.Run(fmt.Sprintf("Difficulty_%s", difficulty), func(t *testing.T) {
			suite := NewAITestSuite(200, difficulty)
			winRate := suite.TestAgainstBaseline()
			
			tolerance := 0.1 // ±10% tolerance
			if math.Abs(winRate-expectedWinRates[i]) > tolerance {
				t.Errorf("AI difficulty %s: expected win rate %.2f%%, got %.2f%%",
					difficulty, expectedWinRates[i]*100, winRate*100)
			}
		})
	}
}

// TestAIPerformanceRegression ensures new changes don't hurt performance
func TestAIPerformanceRegression(t *testing.T) {
	baseline := LoadBaselineMetrics("baseline_performance.json")
	current := NewAITestSuite(100, "medium")
	currentMetrics := current.RunFullPerformanceTest()
	
	// Check for significant performance degradation
	if currentMetrics.BidAccuracy < baseline.BidAccuracy-0.05 {
		t.Errorf("Bid accuracy regression: baseline %.2f%%, current %.2f%%",
			baseline.BidAccuracy*100, currentMetrics.BidAccuracy*100)
	}
	
	if currentMetrics.AverageDecisionTime > time.Duration(float64(baseline.AverageDecisionTime)*1.5) {
		t.Errorf("Decision time regression: baseline %v, current %v",
			baseline.AverageDecisionTime, currentMetrics.AverageDecisionTime)
	}
}

// BenchmarkAIDecisionSpeed measures AI response time
func BenchmarkAIDecisionSpeed(b *testing.B) {
	suite := NewAITestSuite(1, "medium")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate AI decision making
		suite.BenchmarkDecisionTime()
	}
}

// NewAITestSuite creates a new test suite with specified parameters
func NewAITestSuite(games int, difficulty string) *AITestSuite {
	return &AITestSuite{
		TestGames:       games,
		DifficultyLevel: difficulty,
		Opponents:       []string{"easy", "medium", "hard"},
		Metrics:         make(map[string]*AIPerformanceMetrics),
	}
}

// RunFullPerformanceTest executes comprehensive AI evaluation
func (suite *AITestSuite) RunFullPerformanceTest() *AIPerformanceMetrics {
	fmt.Printf("Running %d games to evaluate AI performance...\n", suite.TestGames)
	
	metrics := &AIPerformanceMetrics{
		PlayerName:           fmt.Sprintf("AI_%s", suite.DifficultyLevel),
		TricksWonVsPredicted: make(map[int]int),
		GameResults:          make([]GameResult, 0, suite.TestGames),
	}
	
	for gameNum := 0; gameNum < suite.TestGames; gameNum++ {
		result := suite.SimulateGame(gameNum)
		suite.updateMetrics(metrics, result)
		
		if gameNum%10 == 0 {
			fmt.Printf("Completed %d/%d games\n", gameNum, suite.TestGames)
		}
	}
	
	suite.calculateFinalMetrics(metrics)
	suite.printPerformanceReport(metrics)
	
	return metrics
}

// SimulateGame runs a single game and returns detailed results
func (suite *AITestSuite) SimulateGame(gameNum int) GameResult {
	startTime := time.Now()
	
	// Create players: 1 AI under test + 3 opponents
	playerNames := []string{} // No human players
	players := player.Register(playerNames, 4)
	
	// Set AI difficulty for the player under test
	players[0].Name = fmt.Sprintf("TestAI_%s", suite.DifficultyLevel)
	
	// Initialize and run game
	testGame := game.InitGame(players)
	testGame.Run(0) // Start with player 0 as dealer
	
	// Collect game results
	result := GameResult{
		FinalScore:   players[0].Score,
		FinalRank:    suite.calculateRank(players, 0),
		GameDuration: time.Since(startTime),
	}
	
	return result
}

// TestBiddingScenario tests AI bidding in specific card scenarios
func (suite *AITestSuite) TestBiddingScenario(handSize int, trumpSuit string) float64 {
	correctBids := 0
	totalTests := 50
	
	for i := 0; i < totalTests; i++ {
		// Create controlled hand scenario
		hand := suite.generateTestHand(handSize, trumpSuit)
		actualBid := suite.getAIBid(hand, trumpSuit)
		expectedTricks := suite.simulateHandOutcome(hand, trumpSuit)
		
		if actualBid == expectedTricks {
			correctBids++
		}
	}
	
	return float64(correctBids) / float64(totalTests)
}

// TestCardPlayStrategy evaluates AI's strategic card play
func (suite *AITestSuite) TestCardPlayStrategy(scenario string) float64 {
	correctDecisions := 0
	totalDecisions := 100
	
	for i := 0; i < totalDecisions; i++ {
		gameState := suite.generateGameState(scenario)
		aiChoice := suite.getAICardChoice(gameState)
		optimalChoice := suite.getOptimalChoice(gameState)
		
		if suite.isStrategicallyCorrect(aiChoice, optimalChoice, scenario) {
			correctDecisions++
		}
	}
	
	return float64(correctDecisions) / float64(totalDecisions)
}

// TestAgainstBaseline tests AI against a baseline difficulty
func (suite *AITestSuite) TestAgainstBaseline() float64 {
	wins := 0
	
	for i := 0; i < suite.TestGames; i++ {
		if suite.playAgainstBaseline() {
			wins++
		}
	}
	
	return float64(wins) / float64(suite.TestGames)
}

// Helper methods (implement based on actual game logic)

func (suite *AITestSuite) updateMetrics(metrics *AIPerformanceMetrics, result GameResult) {
	metrics.GamesPlayed++
	metrics.TotalScore += result.FinalScore
	metrics.GameResults = append(metrics.GameResults, result)
	
	if result.FinalRank == 1 {
		metrics.GamesWon++
	}
}

func (suite *AITestSuite) calculateFinalMetrics(metrics *AIPerformanceMetrics) {
	if metrics.GamesPlayed > 0 {
		// Calculate bid accuracy, tendency, etc.
		totalBidError := 0.0
		exactBids := 0
		
		for _, result := range metrics.GameResults {
			if result.BidsTotal > 0 {
				exactBids += result.BidsSuccessful
				totalBidError += result.AvgBidError
			}
		}
		
		metrics.BidAccuracy = float64(exactBids) / float64(metrics.GamesPlayed)
		metrics.BidTendency = totalBidError / float64(metrics.GamesPlayed)
	}
}

func (suite *AITestSuite) printPerformanceReport(metrics *AIPerformanceMetrics) {
	fmt.Printf("\n=== AI Performance Report ===\n")
	fmt.Printf("Player: %s\n", metrics.PlayerName)
	fmt.Printf("Games Played: %d\n", metrics.GamesPlayed)
	fmt.Printf("Win Rate: %.2f%%\n", float64(metrics.GamesWon)/float64(metrics.GamesPlayed)*100)
	fmt.Printf("Average Score: %.1f\n", float64(metrics.TotalScore)/float64(metrics.GamesPlayed))
	fmt.Printf("Bid Accuracy: %.2f%%\n", metrics.BidAccuracy*100)
	fmt.Printf("Bid Tendency: %.2f (positive = overbids)\n", metrics.BidTendency)
	fmt.Printf("Average Decision Time: %v\n", metrics.AverageDecisionTime)
	fmt.Printf("===========================\n")
}

// Placeholder implementations - replace with actual game logic
func (suite *AITestSuite) generateTestHand(handSize int, trumpSuit string) []interface{} {
	return make([]interface{}, handSize) // Implement based on card.Card
}

func (suite *AITestSuite) getAIBid(hand []interface{}, trumpSuit string) int {
	return 0 // Implement AI bidding logic
}

func (suite *AITestSuite) simulateHandOutcome(hand []interface{}, trumpSuit string) int {
	return 0 // Implement hand simulation
}

func (suite *AITestSuite) generateGameState(scenario string) interface{} {
	return nil // Implement game state generation
}

func (suite *AITestSuite) getAICardChoice(gameState interface{}) interface{} {
	return nil // Implement AI card selection
}

func (suite *AITestSuite) getOptimalChoice(gameState interface{}) interface{} {
	return nil // Implement optimal play calculation
}

func (suite *AITestSuite) isStrategicallyCorrect(aiChoice, optimal interface{}, scenario string) bool {
	return true // Implement strategic correctness check
}

func (suite *AITestSuite) playAgainstBaseline() bool {
	return true // Implement baseline comparison
}

func (suite *AITestSuite) calculateRank(players player.Players, playerIndex int) int {
	return 1 // Implement ranking calculation
}

func (suite *AITestSuite) BenchmarkDecisionTime() {
	// Implement decision timing benchmark
}

func LoadBaselineMetrics(filename string) *AIPerformanceMetrics {
	// Implement baseline loading from file
	return &AIPerformanceMetrics{}
}