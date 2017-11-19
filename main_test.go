package main

import "reflect"
import "testing"

func TestBasic(t *testing.T) {
	tape := []byte("aaaaaa")

	// prepare machine
	q0 := state{Name: "q0", Transitions: nil, Final: true}
	t1 := transition{Destination: &q0,
		CurrentSymbol: []byte("a")[0],
		NewSymbol:     []byte("A")[0],
		Action:        right}
	q0.attach_transition(&t1)

	// execute machine
	result := run(&q0, tape)
	expected := []byte("AAAAAA")
	if !reflect.DeepEqual(result.Tape, expected) {
		t.Error("Expected result does not match")
	}
}

func TestInfiniteLoop(t *testing.T) {
	tape := []byte("aaaaba")

	// prepare machine
	q0 := state{Name: "q0", Transitions: nil, Final: true}
	t1 := transition{Destination: &q0,
		CurrentSymbol: []byte("a")[0],
		NewSymbol:     []byte("A")[0],
		Action:        right}
	q0.attach_transition(&t1)

	// execute machine
	result := run(&q0, tape)
	expected := []byte{}
	if !reflect.DeepEqual(result.Tape, expected) {
		t.Error("Expected result does not match")
	}
}

func TestLeftAndRight(t *testing.T) {
	tape := []byte("ababaa")

	// prepare machine
	q0 := state{Name: "q0", Transitions: nil, Final: true}
	t1 := transition{Destination: &q0,
		CurrentSymbol: []byte("a")[0],
		NewSymbol:     []byte("A")[0],
		Action:        right}
	q0.attach_transition(&t1)

	t2 := transition{Destination: &q0,
		CurrentSymbol: []byte("b")[0],
		NewSymbol:     []byte("B")[0],
		Action:        left}
	q0.attach_transition(&t2)

	t3 := transition{Destination: &q0,
		CurrentSymbol: []byte("A")[0],
		NewSymbol:     []byte("a")[0],
		Action:        right}
	q0.attach_transition(&t3)

	t4 := transition{Destination: &q0,
		CurrentSymbol: []byte("B")[0],
		NewSymbol:     []byte("B")[0],
		Action:        right}
	q0.attach_transition(&t4)

	// execute machine
	result := run(&q0, tape)
	expected := []byte("aBaBAA")
	if !reflect.DeepEqual(result.Tape, expected) {
		t.Error("Expected result does not match")
	}
}

func TestTwoStateMachine(t *testing.T) {
	tape := []byte("ababaa")

	// prepare machine
	// define states
	q0 := state{Name: "q1", Transitions: nil, Final: true}
	q1 := state{Name: "q1", Transitions: nil, Final: true}

	// define trasitions
	t0 := transition{Destination: &q1,
		CurrentSymbol: []byte("a")[0],
		NewSymbol:     []byte("a")[0],
		Action:        right}
	q0.attach_transition(&t0)

	t1 := transition{Destination: &q1,
		CurrentSymbol: []byte("a")[0],
		NewSymbol:     []byte("A")[0],
		Action:        right}
	q1.attach_transition(&t1)

	t2 := transition{Destination: &q1,
		CurrentSymbol: []byte("b")[0],
		NewSymbol:     []byte("B")[0],
		Action:        left}
	q1.attach_transition(&t2)

	t3 := transition{Destination: &q1,
		CurrentSymbol: []byte("A")[0],
		NewSymbol:     []byte("a")[0],
		Action:        right}
	q1.attach_transition(&t3)

	t4 := transition{Destination: &q1,
		CurrentSymbol: []byte("B")[0],
		NewSymbol:     []byte("B")[0],
		Action:        right}
	q1.attach_transition(&t4)

	// execute machine
	result := run(&q0, tape)
	expected := []byte("ABaBAA")
	if !reflect.DeepEqual(result.Tape, expected) {
		t.Error("Expected result does not match")
	}
}

func TestAmbiguousTransition(t *testing.T) {
	q0 := state{Name: "q1", Transitions: nil, Final: true}

	// define trasitions
	t0 := transition{Destination: &q0,
		CurrentSymbol: []byte("a")[0],
		NewSymbol:     []byte("a")[0],
		Action:        right}
	err := q0.attach_transition(&t0)
	if err != nil {
		t.Error("Error not expected")
	}

	t1 := transition{Destination: &q0,
		CurrentSymbol: []byte("a")[0],
		NewSymbol:     []byte("A")[0],
		Action:        right}
	err = q0.attach_transition(&t1)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
