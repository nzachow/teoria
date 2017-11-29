package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
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

func (t *transition) String() string {
	return "< Tname=" + t.Name + ": Current " + string(t.CurrentSymbol) + " Dest " + t.Destination.Name + " > "
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
	// log.Fatal(http.ListenAndServe(":8080", router))
	cert := "/etc/letsencrypt/live/teoria.nicolas.eti.br/fullchain.pem"
	key := "/etc/letsencrypt/live/teoria.nicolas.eti.br/privkey.pem"
	log.Fatal(http.ListenAndServeTLS(":8080", cert, key, router))
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

		if data.Word == "*" {
			respondWithError(w, http.StatusBadRequest, "Dados incompletos")
			return
		}
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

			if t.Action != "" {
				if (t.Action == "R") || (t.Action == "L") {
					if (t.TransitionSymbol != "") && (t.WriteSymbol != "") {
						new_transition := transition{Destination: nil,
							CurrentSymbol: []byte(t.TransitionSymbol)[0], NewSymbol: []byte(t.WriteSymbol)[0],
							Action: f, TargetString: t.TargetState, Name: t.Name}
						log.Println("created transition: ", new_transition)
						transitions = append(transitions, new_transition)
					}
				} else {
					respondWithError(w, http.StatusBadRequest, "Dados incompletos")
					return

				}
			}
		}

		var states []state
		for _, s := range data.States {
			if s.Name != "" {
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
			} else {
				respondWithError(w, http.StatusBadRequest, "Dados incompletos")
				return
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
		// data.Word = strings.Replace(data.Word, "β", "", -1)
		var wdr [][]byte
		for _, c := range data.Word {
			wdr = append(wdr, []byte(string(c)))
		}
		// r := run(&states[0], []byte(data.Word))
		r := run(&states[0], wdr)

		json.NewEncoder(w).Encode(r)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.WriteHeader(code)
	w.Write(response)
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

func run(start_state *state, tape [][]byte) execution_result {
	start := time.Now()
	time_limit := 5 * time.Second
	steps := 0
	current_state := start_state
	head_location := 0
	for {
		if time.Now().Sub(start) < time_limit {
			// log.Println("Current state: ", current_state.Name, " has ", len(current_state.Transitions), " transitions")
			var flag bool
			flag = false
			if ((head_location < len(tape)) && (head_location >= 0)) &&
				(len(current_state.Transitions) != 0) {
				for _, t := range current_state.Transitions {
					if head_location < len(tape) && !flag {
						// bytes.Compare
						// if tape[head_location] == (t.CurrentSymbol) {
						if tape[head_location][:] == (t.CurrentSymbol[:]) {
							log.Println("On state: ", current_state.Name)
							log.Printf("Available transitions: %+v", current_state.Transitions)
							log.Println("Executing transition : ", t.Name)
							log.Printf("Tape has %v, transition on %v ", string(tape[head_location]), string(t.CurrentSymbol))
							log.Printf("tape: %s, %v, %T",
								tape, head_location, tape[head_location])
							log.Println("Going to state: ", t.Destination.Name)
							log.Println(" ----------")
							tape[head_location] = t.NewSymbol
							head_location = t.Action(head_location)
							current_state = t.Destination
							steps += 1
							flag = true
						}
					}
				}
				flag = false
			} else {
				log.Println("Execution finished", steps)
				log.Println("Steps:", steps)
				log.Println("Final state:", current_state)
				log.Println("Final state final?:", current_state.Final)

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
