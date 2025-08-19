# Wizard AI Implementation Guide

This guide provides detailed implementation steps for creating AI players in the Wizard card game.

## Overview

The AI needs to handle two main phases:
1. **Bidding/Prediction Phase**: Accurately predict how many tricks it will win
2. **Card Playing Phase**: Take exactly the predicted number of tricks

## Step 1: Complete AI Prediction Logic

**Location**: `prediction_manager.go:22-23`

Bidding Heuristics:
```  // Pseudo-logic for bidding
  probableTricks := 0
  for card := range hand {
      if card.IsWizard() {
          probableTricks += 1.0
      } else if card.IsTrump() {
          probableTricks += trumpWinProbability(card, seenCards)
      } else {
          probableTricks += offSuitWinProbability(card, suit, seenCards)
      }
  }
  bid := round(probableTricks * conservatismFactor)
```
### Hand Evaluation Algorithm

```
expectedTricks = countWizards(hand) + 
                 evaluateTrumpCards(hand, trump) + 
                 evaluateOffSuitCards(hand, trump)
```
 func selectCard(trick []Card, hand []Card, tricksNeeded int) Card {
      if tricksNeeded > 0 {
          return playToWin(trick, hand)
      } else if tricksStaken >= bid {
          return playToLose(trick, hand)
      }
      return playOptimal(trick, hand, tricksNeeded)
  }


### Card Value Assessment

- **Wizards**: Guaranteed tricks (count = 1 each)
- **High trump cards**: Likely tricks (probability based on card rank)
- **High off-suit cards**: Possible tricks (lower probability)
- **Jokers**: Usually 0 value (unless all players play jokers)

### AI Personality Types

- **Conservative AI**: `bid = floor(expectedTricks * 0.9)`
- **Balanced AI**: `bid = round(expectedTricks)`
- **Aggressive AI**: `bid = ceil(expectedTricks * 1.1)`

## Step 2: Implement Round/Trick Playing System

### Missing Components Needed

1. **Round playing loop** in `game_round.go`
2. **Card comparison function** (determines which card wins)
3. **AI card selection logic**
4. **Trick winner determination**

### Card Hierarchy (Highest to Lowest)

1. **Wizards** (always win the trick)
2. **Trump suit cards** (ordered by number, high to low)
3. **Lead suit cards** (ordered by number, high to low)
4. **Other suits** (cannot win unless no trump/lead suit played)
5. **Jokers** (always lose, except in all-joker tricks where first joker wins)

## Step 3: AI Card Playing Strategy

### When AI Needs More Tricks

- Play **lowest winning card** if possible
- Play **Wizards** when necessary for guaranteed wins
- **Lead with high trump cards** to force opponents to use high cards
- **Save Wizards** for crucial moments

### When AI Wants to Avoid Tricks

- Play **highest losing card** to get rid of dangerous cards
- **Use Jokers** to safely lose unwanted tricks
- **Avoid playing Wizards** unless absolutely necessary
- **Follow suit with low cards** when possible

### When AI Has Exact Bid

- Play **defensively** to avoid accidental wins
- **Carefully manage remaining Wizards** (save or play strategically)
- **Block opponents** who still need tricks
- **Help opponents** who have overbid

## Step 4: Special Card Logic

### Wizard Handling

- **Strategic timing**: Save for crucial moments unless you need guaranteed tricks
- **Leading with Wizard**: Next player can play any card (no suit restriction)
- **Trump selection**: If Wizard is trump card, dealer chooses trump suit
- **Count remaining**: Track how many Wizards are still in play

### Joker Handling

- **Defensive tool**: Use to safely lose unwanted tricks
- **Leading with Joker**: Next real suit card determines the suit to follow
- **All-Joker tricks**: First Joker played wins the trick
- **Timing**: Play early when you need to lose, save when others might lead with Jokers

## Step 5: Game State Tracking

### Information to Track

- **Cards played each round** (by suit and number)
- **Each player's remaining trick needs** (bid vs current tricks)
- **High cards still in play** by suit
- **Trump cards remaining** in the game
- **Wizards and Jokers** still unplayed

### Implementation Locations

- **Player struct**: Add AI state fields for tracking
- **Turn struct**: Add game state tracking
- **Round results**: Track opponent patterns over multiple games

### Advanced AI Features

#### Opponent Modeling

- **Bidding patterns**: Conservative vs aggressive tendencies
- **Play style**: Risk-taking vs safe play
- **Card counting**: Remember what cards opponents have played
- **Adaptation**: Adjust strategy based on opponent behavior

#### Difficulty Levels

- **Beginner**: Basic probability, no opponent modeling
- **Intermediate**: Card counting, simple opponent tracking  
- **Advanced**: Full opponent modeling, situational awareness
- **Expert**: Monte Carlo simulations, multi-step planning

## Implementation Priority

1. **Start with basic bidding logic** (Step 1)
2. **Implement card comparison and winner determination** (Step 2)
3. **Add simple AI card selection** (Step 3)
4. **Handle special cards** (Step 4)
5. **Add game state tracking** (Step 5)
6. **Implement difficulty levels and opponent modeling** (Advanced)

## Key Design Principles

- **Exact bidding is crucial** - only exact bids score points
- **Conservative approach often wins** - overbidding and underbidding both lose points
- **Special cards are game-changers** - Wizards and Jokers require careful timing
- **Opponent tracking improves performance** - learning opponent patterns gives advantages
- **Adaptability matters** - AI should adjust strategy based on game state

## Testing Strategy

- **Unit tests** for hand evaluation functions
- **Integration tests** for complete games against AI
- **Performance tests** for different difficulty levels
- **Regression tests** to ensure AI improvements don't break existing behavior