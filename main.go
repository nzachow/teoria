package main

import (
	"log"
	"time"
)

type state struct {
	Name        string
	Transitions []*transition
	Final       bool
}

func (s *state) attach_transition(t *transition) {
	s.Transitions = append(s.Transitions, t)
}

type transition struct {
	Destination   *state
	CurrentSymbol byte
	NewSymbol     byte
	Action        func(int) int
}

func right(counter int) int {
	log.Println("go right!")
	return counter + 1
}

func left(counter int) int {
	log.Println("go left!")
	return counter - 1
}

func main() {
	// prepare tape
	// tape := []byte("aaabaa")

	// // prepare machine
	// q0 := state{Name: "q0", Transitions: nil, Final: true}
	// t1 := transition{Destination: &q0,
	// 	CurrentSymbol: []byte("a")[0],
	// 	NewSymbol:     []byte("A")[0],
	// 	Action:        right}
	// q0.attach_transition(&t1)

	// t2 := transition{Destination: &q0,
	// 	CurrentSymbol: []byte("b")[0],
	// 	NewSymbol:     []byte("B")[0],
	// 	Action:        left}
	// q0.attach_transition(&t2)

	// t3 := transition{Destination: &q0,
	// 	CurrentSymbol: []byte("A")[0],
	// 	NewSymbol:     []byte("a")[0],
	// 	Action:        right}
	// q0.attach_transition(&t3)

	// t4 := transition{Destination: &q0,
	// 	CurrentSymbol: []byte("B")[0],
	// 	NewSymbol:     []byte("B")[0],
	// 	Action:        right}
	// q0.attach_transition(&t4)

	// // execute machine
	// result := run(&q0, tape)
	// log.Printf("Result on tape: %s", result)
}

func run(start_state *state, tape []byte) []byte {
	start := time.Now()
	time_limit := 2 * time.Second
	current_state := start_state
	head_location := 0
	for {
		if time.Now().Sub(start) < time_limit {
			if (head_location < len(tape)) && (head_location >= 0) {
				for _, t := range current_state.Transitions {
					if head_location < len(tape) {
						if tape[head_location] == (t.CurrentSymbol) {
							tape[head_location] = t.NewSymbol
							log.Printf("tape: %s, %v, %T",
								tape, head_location, tape[head_location])
							head_location = t.Action(head_location)
						}
					}
				}
			} else {
				log.Println("Execution finished")
				return tape
			}
		} else {
			log.Println("Time exceeded, halting execution")
			return []byte{}
		}
	}
}
