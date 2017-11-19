package main

import "log"

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
	tape := []byte("aaaaaa")

	// prepare machine
	q0 := state{Name: "q0", Transitions: nil, Final: true}
	t1 := transition{Destination: &q0,
		CurrentSymbol: []byte("a")[0],
		NewSymbol:     []byte("A")[0],
		Action:        right}
	q0.attach_transition(&t1)

	// execute machine
	current_state := &q0
	head_location := 0
	for {
		if head_location < (len(tape) - 1) {
			for _, t := range current_state.Transitions {
				if tape[head_location] == (t.CurrentSymbol) {
					tape[head_location] = t.NewSymbol
					log.Printf("tape: %s, %v, %T", tape, head_location, tape[head_location])
					head_location = t.Action(head_location)
				}
			}
		} else {
			log.Println("Execution finished")
			log.Println("Final tape: ", tape)
			return
		}
	}
}
