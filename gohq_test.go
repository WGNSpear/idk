package gohq

import (
	"testing"
	"log"
	"fmt"
	"strconv"
	"net/http"
	"io/ioutil"
)

func TestHQ(t *testing.T) {
	acc := Account{AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEzOTYzNzkxLCJ1c2VybmFtZSI6IkRpc2NvbGkiLCJhdmF0YXJVcmwiOiJzMzovL2h5cGVzcGFjZS1xdWl6L2EvYzEvMTM5NjM3OTEtM0ZxSGVVLmpwZyIsInRva2VuIjoiQ3dldVB1Iiwicm9sZXMiOltdLCJjbGllbnQiOiIiLCJndWVzdElkIjpudWxsLCJ2IjoxLCJpYXQiOjE1Mjc3MTU3NDEsImV4cCI6MTUzNTQ5MTc0MSwiaXNzIjoiaHlwZXF1aXovMSJ9.fxrUu6cc_Iv4H2i4_lJV6lxPwgFaMxeJxlFFtFJGqng"}
	aws, _ := acc.RequestAWS()
	fmt.Println("Session Token: " + aws.SessionToken)

	resp, _ := http.DefaultClient.Get("https://picsum.photos/200?random")
	bytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	aws.UploadAvatarAWS("discoli.gif", bytes)
	acc.ChangeAvatarAWS("discoli.gif")

	log.Fatal("")

	game, err := Debug()
	if err != nil {
		log.Fatal(err)
	}

	for {
		bytes, err := game.Read()
		if err != nil {
			log.Fatal(err)
		}

		if stats := game.ParseBroadcastStats(bytes); stats != nil {
			continue
		} else if message := game.ParseChatMessage(bytes); message != nil {
			continue
		} else if gameStatus := game.ParseGameStatus(bytes); gameStatus != nil {
			fmt.Println("You have joined, the prize is", gameStatus.Prize)
		} else if question := game.ParseQuestion(bytes); question != nil {
			fmt.Println("Question Incoming, you have", strconv.Itoa(question.TimeLeftMs)+"ms to answer it!")
		} else if questionClosed := game.ParseQuestionClosed(bytes); questionClosed != nil {
			fmt.Println("The question", questionClosed.QuestionID, "is over!")
		} else if questionSummary := game.ParseQuestionSummary(bytes); questionSummary != nil {
			fmt.Println(questionSummary.EliminatedPlayersCount, "players have been eliminated,", questionSummary.AdvancingPlayersCount, "players remain.")
		} else if questionFinished := game.ParseQuestionFinished(bytes); questionFinished != nil {
			fmt.Println("The question", questionFinished.QuestionID, "has finished.")
		} else {
			fmt.Println("This is an unknown message:", string(bytes))
		}
	}
}
