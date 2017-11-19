package main

import (
	"errors"
	"log"
	"time"
)

type state struct {
	Name        string
	Transitions []*transition
	Final       bool
}

func (s *state) attach_transition(t *transition) error {
	// prevent ambiguous transitions
	for _, v := range s.Transitions {
		if t.CurrentSymbol == v.CurrentSymbol {
			return errors.New("Cannot add ambiguous transition")
		}
	}
	s.Transitions = append(s.Transitions, t)
	return nil
}

type execution_result struct {
	// finished on a final state ?
	FinalState bool
	Steps      int
	Tape       []byte
}

type transition struct {
	Destination   *state
	CurrentSymbol byte
	NewSymbol     byte
	Action        func(int) int
}

func right(counter int) int {
	return counter + 1
}

func left(counter int) int {
	return counter - 1
}

func main() {
}

func run(start_state *state, tape []byte) execution_result {
	start := time.Now()
	time_limit := 2 * time.Second
	steps := 0
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
							current_state = t.Destination
							steps += 1
						}
					}
				}
			} else {
				log.Println("Execution finished", steps)

				res := execution_result{current_state.Final, steps, tape}
				return res
			}
		} else {
			log.Println("Time exceeded, halting execution")
			res := execution_result{false, steps, []byte{}}
			return res
		}
	}
}
