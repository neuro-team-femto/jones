package ws

import (
	"encoding/json"

	"github.com/creamlab/revcor/xp"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// messages in

type joinData struct {
	ExperimentId  string `json:"experimentId"`
	ParticipantId string `json:"participantId"`
}

type resultData struct {
	Chosen    string `json:"chosen"`
	Dismissed string `json:"dismissed"`
}

type messageIn struct {
	Kind    string `json:"kind"`
	Payload string `json:"payload"`
}

// messages out
type outData map[string]interface{}

type messageOut struct {
	Kind    string  `json:"kind"`
	Payload outData `json:"payload"`
}

func RunServer(conn *websocket.Conn) {
	defer conn.Close()

	// first message must be a join request
	joinMsg := messageIn{}
	err := conn.ReadJSON(&joinMsg)
	if err != nil || joinMsg.Kind != "join" {
		sendAndLogError(conn, err, "error-no-join")
		return
	}

	joinPayload, err := readJoin(joinMsg)
	if err != nil || !xp.IsValid(joinPayload.ExperimentId, joinPayload.ParticipantId) {
		sendAndLogError(conn, err, "error-wrong-join")
		return
	}

	es, err := xp.GetExperimentSettings(joinPayload.ExperimentId)
	if err != nil {
		return
	}

	ew, err := xp.GetExperimentWordingString(joinPayload.ExperimentId)
	if err != nil {
		return
	}

	p, err := xp.LoadParticipant(es, joinPayload.ParticipantId)
	if err != nil {
		return
	}

	// caution: es/p are structs that will be automatically deserialized as JS objects client side
	// ew is a string that remains to be parsed client-side (done this way not to declare wordings.json structure)
	initPayload := outData{
		"settings":    es,
		"wording":     ew,
		"participant": p,
	}

	if err := sendWithPayload(conn, "init", initPayload); err != nil {
		return
	}

	for {
		msg := messageIn{}
		err := conn.ReadJSON(&msg)

		if err != nil {
			return
		}

		if msg.Kind == "result" {
			result, err := readResult(msg)
			if err != nil {
				sendAndLogError(conn, err, "error-read-result")
				return
			}
			xp.UpdateParticipantState(es, p, result.Chosen, result.Dismissed)
		}
	}
}

func sendAndLogError(conn *websocket.Conn, err error, errorMsg string) {
	send(conn, errorMsg)
	log.Error().Err(err).Msg("[ws] " + errorMsg)
}

func send(conn *websocket.Conn, kind string) (err error) {
	m := &messageOut{Kind: kind}

	if err := conn.WriteJSON(m); err != nil {
		log.Error().Err(err).Msg("[ws] can't write JSON")
	}
	return
}

func sendWithPayload(conn *websocket.Conn, kind string, payload outData) (err error) {
	m := &messageOut{Kind: kind, Payload: payload}

	if err := conn.WriteJSON(m); err != nil {
		log.Error().Err(err).Msg("[ws] can't write JSON")
	}
	return
}

func readJoin(msg messageIn) (join joinData, err error) {
	err = json.Unmarshal([]byte(msg.Payload), &join)
	return
}

func readResult(msg messageIn) (result resultData, err error) {
	err = json.Unmarshal([]byte(msg.Payload), &result)
	return
}
