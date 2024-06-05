package service

import (
	"Server-go/model"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

func NewDeck() []model.Card {
	rank := []int{4, 3, 2, 1}
	values := []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	deck := make([]model.Card, 0)
	var color string
	var suitStr string
	for _, suit := range rank {
		for _, value := range values {
			if suit == 4 || suit == 3 {
				if suit == 4 {
					suitStr = "♥️"
				} else {
					suitStr = "♦️"
				}
				color = `red`
			} else {
				if suit == 2 {
					suitStr = "♣️"
				} else {
					suitStr = "♠️"
				}
				color = `black`
			}
			card := model.Card{Rank: suit, Value: value, Suit: suitStr, Color: color}
			deck = append(deck, card)
		}
	}

	return deck
}

func ShuffleDeck(deck []model.Card) []model.Card {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}

// DealCards deals cards to players.
func DealCards() map[int][]model.Card {
	deck := ShuffleDeck(NewDeck())
	hands := make(map[int][]model.Card)
	numPlayers := 4
	cardsPerHand := 13

	for i := 0; i < cardsPerHand; i++ {
		for j := 0; j < numPlayers; j++ {
			cardIndex := i*numPlayers + j
			hands[j] = append(hands[j], deck[cardIndex])
		}
	}

	var wg sync.WaitGroup

	for i := 0; i < numPlayers; i++ {
		wg.Add(1)
		go func(hand []model.Card) {
			defer wg.Done()
			SortHand(hand)
		}(hands[i])
	}

	wg.Wait()
	fmt.Println("All hands sorted.")
	return hands
}

func SortHand(hands []model.Card) []model.Card {
	straightRankMapValue := make(map[int][]model.Card)
	var strongValue []model.Card
	var anotherValue []model.Card
	var straight []model.Card
	var result []model.Card
	sort.Slice(hands, func(i, j int) bool {
		return hands[i].Value <= hands[j].Value
	})
	for _, hand := range hands {
		if hand.Value == 15 {
			strongValue = append(strongValue, hand)
			continue
		}
		if len(straightRankMapValue[hand.Rank]) >= 1 {
			//Kiểm tra trong sảnh có giá trị nào bằng lá bài tiếp theo k
			//Lá bài tiếp theo có index = length của mảng đấy - 1
			if straightRankMapValue[hand.Rank][len(straightRankMapValue[hand.Rank])-1].Value == hand.Value-1 {
				straightRankMapValue[hand.Rank] = append(straightRankMapValue[hand.Rank], hand)
				continue
				// Không có thì xóa lá bài cũ và thêm lá bài hiện tại vào
			} else if len(straightRankMapValue[hand.Rank]) <= 2 {
				//Các lá bài không phải straight vào
				anotherValue = append(anotherValue, straightRankMapValue[hand.Rank]...)
				//Tạo mới straight
				delete(straightRankMapValue, hand.Rank)
				straightRankMapValue[hand.Rank] = append(straightRankMapValue[hand.Rank], hand)
				continue
			} else {
				//Trường hợp đã có straight , lá bài cùng chất sẽ rơi vào trường hợp này
				anotherValue = append(anotherValue, hand)
			}
		} else {
			//Khởi tạo straight
			straightRankMapValue[hand.Rank] = append(straightRankMapValue[hand.Rank], hand)
		}
	}
	//Tìm các cặp
	for _, cards := range straightRankMapValue {
		if len(cards) <= 2 {
			anotherValue = append(anotherValue, cards...)
		} else {
			straight = append(straight, cards...)
		}
	}
	pairsRankMapValue := make(map[int][]model.Card)

	for _, card := range anotherValue {
		pairsRankMapValue[card.Value] = append(pairsRankMapValue[card.Value], card)
	}

	var singleCard []model.Card
	var pairs []model.Card
	for _, cards := range pairsRankMapValue {
		if len(cards) == 1 {
			singleCard = append(singleCard, cards...)
		} else {
			if len(cards) == 2 {
				if cards[0].Color == cards[1].Color {
					pairs = append(pairs, cards...)
				} else {
					singleCard = append(singleCard, cards...)
				}
			} else {
				pairs = append(pairs, cards...)
			}
		}
	}
	sort.Slice(singleCard, func(i, j int) bool {
		return singleCard[i].Value <= singleCard[j].Value
	})
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value <= pairs[j].Value
	})
	result = append(append(append(append([]model.Card{}, singleCard...), pairs...), straight...), strongValue...)

	fmt.Println(result)
	fmt.Println("------------------------")
	return result

}
