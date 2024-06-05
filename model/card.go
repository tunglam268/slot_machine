package model

const (
	PlayerWon = "PLAYER WON"
	DealerWon = "DEALER WON"
	Draw      = "DRAW"
	BlackJack = "BLACKJACK"
)

type Card struct {
	Rank  int
	Value int
	Suit  string
	Color string
}
