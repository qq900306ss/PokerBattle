package module

import (
	"reflect"
	"testing"
)

func TestComparePoker(t *testing.T) {
	type args struct {
		player1 Message
		player2 Message
	}
	tests := []struct {
		name string
		args args
		want Message
	}{
		{
			name: "TestComparePoker",
			args: args{
				player1: Message{
					Event: "PokerBattle",
					Name:  "player1",
					Card: Card{
						Value: 10,
						Suit:  "黑桃",
					},
				},
				player2: Message{
					Event: "PokerBattle",
					Name:  "player2",
					Card: Card{
						Value: 12,
						Suit:  "黑桃",
					},
				},
			},

			want: Message{
				Event: "PokerBattle",
				Name:  "player1",
				Card: Card{
					Value: 12,
					Suit:  "黑桃",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComparePoker(tt.args.player1, tt.args.player2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComparePoker() = %v, want %v", got, tt.want)
			}
		})
	}
}
