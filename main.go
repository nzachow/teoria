package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
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

//     {"name":"transition11","targetState":"q1","transitionSymbol":"β","writeSymbol":"X","action":"R"},

type transition struct {
	Destination   *state
	CurrentSymbol byte
	NewSymbol     byte
	Action        func(int) int
	TargetString  string
	Name          string
}

func (t *transition) set_destination(d *state) {
	t.Destination = d
}

func right(counter int) int {
	return counter + 1
}

func left(counter int) int {
	return counter - 1
}

func main() {

	// networking code
	router := mux.NewRouter()
	router.HandleFunc("/send", handleWrapper(receive_machine)).Methods("POST", "OPTIONS")
	log.Fatal(http.ListenAndServe(":8080", router))
}

type receive_data struct {
	Machine_name string               `json:"machineName"`
	Word         string               `json:"word"`
	Transtions   []receive_transition `json:"transitions"`
	States       []receive_state      `json:"states"`
}

type receive_transition struct {
	Name             string `json:"name"`
	TargetState      string `json:"targetState"`
	TransitionSymbol string `json:"transitionSymbol"`
	WriteSymbol      string `json:"writeSymbol"`
	Action           string `json:"action"`
}

type receive_state struct {
	Name        string   `json:"name"`
	Transitions []string `json:"transitions"`
	IsFinal     bool     `json:"isFinal"`
}

func receive_machine(w http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		return
	} else {
		decoder := json.NewDecoder(req.Body)
		var data receive_data
		err := decoder.Decode(&data)
		if err != nil {
			// panic(err)
			log.Println("error :", err)
			json.NewEncoder(w).Encode(err)
			return
		}
		defer req.Body.Close()
		log.Println("machine: ", data.Machine_name)
		log.Println("struct", data)

		// {"machineName":"Minha máquina",
		// "word":"",
		// "transitions":[
		//  {"name":"transition11","targetState":"q1","transitionSymbol":"β","writeSymbol":"X","action":"R"},
		//  {"name":"transition21","targetState":"q3","transitionSymbol":"*","writeSymbol":"Y","action":"R"},
		//  {"name":"transition12","targetState":"q1","transitionSymbol":"β","writeSymbol":"X","action":"R"},
		//  {"name":"transition22","targetState":"q0","transitionSymbol":"*","writeSymbol":"X","action":"L"}],
		// "states":[{"name":"q0","transitions":["transition11","transition21"],"isFinal":false},
		//     {"name":"q1","transitions":["transition12","transition22"],"isFinal":true}]}
		var transitions []transition
		for _, t := range data.Transtions {
			var f func(int) int
			if t.Action == "R" {
				f = right
			} else {
				f = left
			}

			new_transition := transition{Destination: nil,
				CurrentSymbol: []byte(t.TransitionSymbol)[0], NewSymbol: []byte(t.WriteSymbol)[0],
				Action: f, TargetString: t.TargetState, Name: t.Name}
			log.Println("created transition: ", new_transition)
			transitions = append(transitions, new_transition)
		}

		var states []state
		for _, s := range data.States {
			new_state := state{Name: s.Name, Transitions: nil, Final: s.IsFinal}
			states = append(states, new_state)
			for _, trn := range s.Transitions {
				for i, v := range transitions {
					if v.Name == trn {
						log.Println("adding :", transitions[i].Name, "to : ", s.Name)
						states[len(states)-1].attach_transition(&transitions[i])
					}
				}
			}
		}

		log.Println("states : ", states)
		for ii, t := range transitions {
			if t.Destination == nil {
				for i, s := range states {
					if s.Name == t.TargetString {
						log.Println("destination of: ", t.Name, "should be: ", s.Name)
						log.Println("setting destination to: ", &s)
						log.Println("setting destination to: ", &states[i])
						transitions[ii].set_destination(&states[i])
					}
				}
			}
		}

		log.Println("states: ", states)
		log.Println("transitions: ", transitions)
		// log.Println("w: ",
		data.Word = strings.Replace(data.Word, "β", "", -1)
		r := run(&states[0], []byte(data.Word))

		json.NewEncoder(w).Encode(r)
	}
}

func handleWrapper(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received: ", r.Header)
		log.Println("Adding headers: ")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Cache-Control, Pragma, Origin, Authorization,   Content-Type, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		f(w, r)
	}
}

func run(start_state *state, tape []byte) execution_result {
	start := time.Now()
	time_limit := 5 * time.Second
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
