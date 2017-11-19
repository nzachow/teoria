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
	if !reflect.DeepEqual(result, expected) {
		t.Error("Expected result does not match")
	}
}
