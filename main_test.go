package main

import "reflect"
import "testing"

func create_tape(s string) [][]byte {

	var tape [][]byte
	for _, c := range s {
		tape = append(tape, []byte(string(c)))
	}
	return tape

}

func TestBasic(t *testing.T) {
	// tape := []byte("aaaaaa")
	// var tape [][]byte
	// for _, c := range "aaaaaa" {
	// 	tape = append(tape, []byte(string(c)))
	// }
	tape := create_tape("aaaaaa")

	// prepare machine
	q0 := state{Name: "q0", Transitions: nil, Final: true}
	t1 := transition{Destination: &q0,
		CurrentSymbol: []byte("a"),
		NewSymbol:     []byte("A"),
		Action:        right}
	q0.attach_transition(&t1)

	// execute machine
	result := run(&q0, tape)
	// exp_tape := []byte("AAAAAA")
	exp_steps := 6
	var exp_tape [][]byte
	for _, c := range "AAAAAA" {
		exp_tape = append(exp_tape, []byte(string(c)))
	}
	if !reflect.DeepEqual(result.Tape, exp_tape) {
		t.Error("Expected result does not match")
	}

	if result.Steps != exp_steps {
		t.Error("Wrong number of steps")
	}
}

func TestNoTrasitions(t *testing.T) {
	tape := create_tape("aaaaba")

	// prepare machine
	q0 := state{Name: "q0", Transitions: nil, Final: true}
	t1 := transition{Destination: &q0,
		CurrentSymbol: []byte("a"),
		NewSymbol:     []byte("A"),
		Action:        right}
	q0.attach_transition(&t1)

	// execute machine
	result := run(&q0, tape)
	expected := [][]byte{}
	exp_steps := 3 // at least 3 steps
	if !reflect.DeepEqual(result.Tape, expected) {
		t.Error("Expected result does not match")
	}

	if result.Steps < exp_steps {
		t.Error("Wrong number of steps")
	}
}

func TestInfiniteLoop(t *testing.T) {
	tape := create_tape("aaaaba")

	// prepare machine
	q0 := state{Name: "q0", Transitions: nil, Final: true}
	t1 := transition{Destination: &q0,
		CurrentSymbol: []byte("a"),
		NewSymbol:     []byte("a"),
		Action:        right}
	q0.attach_transition(&t1)
	t2 := transition{Destination: &q0,
		CurrentSymbol: []byte("b"),
		NewSymbol:     []byte("b"),
		Action:        left}
	q0.attach_transition(&t2)

	// execute machine
	result := run(&q0, tape)
	expected := [][]byte{}
	exp_steps := 100 // at least 3 steps
	if !reflect.DeepEqual(result.Tape, expected) {
		t.Error("Expected result does not match")
	}

	if result.Steps < exp_steps {
		t.Error("Wrong number of steps")
	}
}

func TestLeftAndRight(t *testing.T) {
	tape := create_tape("ababaa")

	// prepare machine
	q0 := state{Name: "q0", Transitions: nil, Final: true}
	t1 := transition{Destination: &q0,
		CurrentSymbol: []byte("a"),
		NewSymbol:     []byte("A"),
		Action:        right}
	q0.attach_transition(&t1)

	t2 := transition{Destination: &q0,
		CurrentSymbol: []byte("b"),
		NewSymbol:     []byte("B"),
		Action:        left}
	q0.attach_transition(&t2)

	t3 := transition{Destination: &q0,
		CurrentSymbol: []byte("A"),
		NewSymbol:     []byte("a"),
		Action:        right}
	q0.attach_transition(&t3)

	t4 := transition{Destination: &q0,
		CurrentSymbol: []byte("B"),
		NewSymbol:     []byte("B"),
		Action:        right}
	q0.attach_transition(&t4)

	// execute machine
	result := run(&q0, tape)
	exp_tape := create_tape("aBaBAA")
	exp_steps := 10
	if !reflect.DeepEqual(result.Tape, exp_tape) {
		t.Error("Expected result does not match")
	}

	if result.Steps != exp_steps {
		t.Error("Wrong number of steps")
	}
}

func TestTwoStateMachine(t *testing.T) {
	tape := create_tape("ababaa")

	// prepare machine
	// define states
	q0 := state{Name: "q1", Transitions: nil, Final: true}
	q1 := state{Name: "q1", Transitions: nil, Final: true}

	// define trasitions
	t0 := transition{Destination: &q1,
		CurrentSymbol: []byte("a"),
		NewSymbol:     []byte("a"),
		Action:        right}
	q0.attach_transition(&t0)

	t1 := transition{Destination: &q1,
		CurrentSymbol: []byte("a"),
		NewSymbol:     []byte("A"),
		Action:        right}
	q1.attach_transition(&t1)

	t2 := transition{Destination: &q1,
		CurrentSymbol: []byte("b"),
		NewSymbol:     []byte("B"),
		Action:        left}
	q1.attach_transition(&t2)

	t3 := transition{Destination: &q1,
		CurrentSymbol: []byte("A"),
		NewSymbol:     []byte("a"),
		Action:        right}
	q1.attach_transition(&t3)

	t4 := transition{Destination: &q1,
		CurrentSymbol: []byte("B"),
		NewSymbol:     []byte("B"),
		Action:        right}
	q1.attach_transition(&t4)

	// execute machine
	result := run(&q0, tape)
	exp_tape := create_tape("ABaBAA")
	exp_steps := 10
	if !reflect.DeepEqual(result.Tape, exp_tape) {
		t.Error("Expected result does not match")
	}

	if result.Steps != exp_steps {
		t.Error("Wrong number of steps")
	}
}

func TestAmbiguousTransition(t *testing.T) {
	q0 := state{Name: "q1", Transitions: nil, Final: true}

	// define trasitions
	t0 := transition{Destination: &q0,
		CurrentSymbol: []byte("a"),
		NewSymbol:     []byte("a"),
		Action:        right}
	err := q0.attach_transition(&t0)
	if err != nil {
		t.Error("Error not expected")
	}

	t1 := transition{Destination: &q0,
		CurrentSymbol: []byte("a"),
		NewSymbol:     []byte("A"),
		Action:        right}
	err = q0.attach_transition(&t1)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
