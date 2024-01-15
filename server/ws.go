package server

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/neuro-team-femto/revcor/xp"
)

type participantConn struct {
	conn *websocket.Conn
	p    xp.Participant
	es   xp.ExperimentSettings
}

// messages in
type messageIn struct {
	Kind    string `json:"kind"`
	Payload string `json:"payload"`
}

// messages in payloads
type joinData struct {
	ExperimentId  string `json:"experimentId"`
	ParticipantId string `json:"participantId"`
}

type infoData struct {
	Age string `json:"age"`
	Sex string `json:"sex"`
}

type trialData struct {
	Result1 xp.Result `json:"result1"`
	Result2 xp.Result `json:"result2"`
}

// messages out
type messageOut struct {
	Kind    string  `json:"kind"`
	Payload outData `json:"payload"`
}

// messages out payloads => untyped
type outData map[string]interface{}

func sendAndLogError(conn *websocket.Conn, err error, errorMsg string) {
	send(conn, errorMsg)
	log.Printf("[error] ws %v: %+v\n", errorMsg, err)
}

func send(conn *websocket.Conn, kind string) (err error) {
	m := &messageOut{Kind: kind}

	if err := conn.WriteJSON(m); err != nil {
		log.Println("[error] ws: can't write JSON")
	}
	return
}

func sendWithPayload(conn *websocket.Conn, kind string, payload outData) (err error) {
	m := &messageOut{Kind: kind, Payload: payload}

	if err := conn.WriteJSON(m); err != nil {
		log.Println("[error] ws: can't write JSON")
	}
	return
}

func (pc participantConn) loop() {
	for {
		msg := messageIn{}
		err := pc.conn.ReadJSON(&msg)

		if err != nil {
			return
		}

		if msg.Kind == "trial" {
			trial := trialData{}
			err = json.Unmarshal([]byte(msg.Payload), &trial)
			if err != nil {
				sendAndLogError(pc.conn, err, "error-trial-read")
				return
			}

			if !trial.Result1.IsValid() || !trial.Result2.IsValid() {
				sendAndLogError(pc.conn, err, "error-trial-invalid")
				return
			}

			err = xp.WriteToCSV(pc.es, pc.p, trial.Result1, trial.Result2)
			if err != nil {
				sendAndLogError(pc.conn, err, "error-trial-write")
				return
			}

			err = pc.p.UpdateTodo(trial.Result1.Stimulus, trial.Result2.Stimulus)
			if err != nil {
				sendAndLogError(pc.conn, err, "error-todo-update")
				return
			}
		}
	}
}

// API

func runWs(conn *websocket.Conn) {
	defer conn.Close()

	// there is an implied protocol to be followed:
	// 1. first received message *must* be a "join"
	joinMsg := messageIn{}
	err := conn.ReadJSON(&joinMsg)
	if err != nil || joinMsg.Kind != "join" {
		sendAndLogError(conn, err, "error-join-missing")
		return
	}

	join := joinData{}
	err = json.Unmarshal([]byte(joinMsg.Payload), &join)
	if err != nil || !xp.IsParticipantValid(join.ExperimentId, join.ParticipantId) {
		sendAndLogError(conn, err, "error-join-invalid")
		return
	}

	es, err := xp.GetSanitizedExperimentSettings(join.ExperimentId)
	if err != nil {
		return
	}

	ew, err := xp.GetExperimentWordingRunString(join.ExperimentId)
	if err != nil {
		return
	}

	p, err := xp.LoadParticipant(es, join.ParticipantId)
	if err != nil {
		return
	}

	// caution: es/p are structs that will be automatically deserialized as JS objects client side
	// ew is a string that remains to be parsed client-side
	// done this way on purpose, not to type/declare wording json evolving structures
	initPayload := outData{
		"settings":    es,
		"participant": p,
		"wording":     ew,
	}

	// 2. first sent message is an "init" containing the data needed to initialized the client state
	if err := sendWithPayload(conn, "init", initPayload); err != nil {
		return
	}

	// 3. if participant info is empty, the next received message *must* be a "info"
	if p.Age == "" || p.Sex == "" {
		infoMsg := messageIn{}
		err := conn.ReadJSON(&infoMsg)
		if err != nil || infoMsg.Kind != "info" {
			sendAndLogError(conn, err, "error-info-missing")
			return
		}

		info := infoData{}
		err = json.Unmarshal([]byte(infoMsg.Payload), &info)
		if err != nil {
			sendAndLogError(conn, err, "error-info-invalid")
			return
		}

		err = p.UpdateInfo(info.Age, info.Sex)
		if err != nil {
			sendAndLogError(conn, err, "error-info-save")
			return
		}
	}

	// 4. client/server initialization is over, now we loop on "result" messages
	pc := participantConn{conn, p, es}
	pc.loop()
}
