package gohq

import (
	"strconv"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"errors"
	"strings"
	"time"
	"github.com/gorilla/websocket"
	"github.com/smartystreets/go-aws-auth"
	"bytes"
	"net/url"
)

func Verify(number string) (*Verification, error) {
	body := `{"method":"sms","phone":"` + number + `"}`
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/verifications", strings.NewReader(body))
	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", strconv.Itoa(len(body)))
	req.Header.Add("user-agent", "okhttp/3.8.0")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	bytes, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	var verification Verification

	if json.Unmarshal(bytes, &verification); verification.VerificationID == "" {
		var hqError HQError
		if json.Unmarshal(bytes, &hqError); hqError.Error != "" {
			return nil, errors.New(hqError.Error)
		} else {
			return nil, errors.New("unknown error: " + string(bytes))
		}
	} else {
		return &verification, nil
	}
}
func (verification *Verification) Confirm(code string) (*Auth, error) {
	body := `{"code":"` + code + `"}`
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/verifications/"+verification.VerificationID, strings.NewReader(body))
	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", strconv.Itoa(len(body)))
	req.Header.Add("user-agent", "okhttp/3.8.0")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	bytes, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	var auth Auth
	var hqError HQError

	if json.Unmarshal(bytes, &auth); auth.Auth.AccessToken == "" {
		if json.Unmarshal(bytes, &hqError); hqError.Error != "" {
			return nil, errors.New(hqError.Error)
		} else {
			if auth.Auth == (Account{}) {
				return nil, nil
			}
			return nil, errors.New("unknown error: " + string(bytes))
		}
	}

	return &auth, nil
}
func (verification *Verification) Create(username, referrer, region string, transport *http.Transport) (*Account, error) {
	if transport == nil {
		transport = &http.Transport{}
	}

	body := `{"country":"` + region + `","language":"en","referringUsername":"` + referrer + `","username":"` + username + `","verificationId":"` + verification.VerificationID + `"}`
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/users", strings.NewReader(body))

	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", strconv.Itoa(len(body)))
	req.Header.Add("user-agent", "okhttp/3.8.0")

	client := http.Client{Transport: transport, Timeout: time.Second * 10}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bytes, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	var account Account
	var hqError HQError

	if json.Unmarshal(bytes, &account); account.AccessToken == "" {
		if json.Unmarshal(bytes, &hqError); hqError.Error != "" {
			return nil, errors.New(hqError.Error)
		} else {
			return nil, errors.New("unknown error: " + string(bytes))
		}
	}

	return &account, nil
}
func (account *Account) Weekly() (error) {
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/easter-eggs/makeItRain", strings.NewReader("{}"))
	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("authorization", "Bearer "+account.AccessToken)
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", "2")
	req.Header.Add("user-agent", "okhttp/3.8.0")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	bytes, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	hqError := HQError{}
	json.Unmarshal(bytes, &hqError)

	if hqError.Error != "" {
		return errors.New(hqError.Error)
	}

	return nil
}

func (account *Account) SearchUser(username string) (*Users, error) {
	req, _ := http.NewRequest("GET", "https://api-quiz.hype.space/users?q=" + url.QueryEscape(username), nil)
	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("authorization", "Bearer "+account.AccessToken)
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("user-agent", "okhttp/3.8.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)

	var users Users
	var hqError HQError

	if json.Unmarshal(bytes, &users); users.Links.Self == "" {
		if json.Unmarshal(bytes, &hqError); hqError.Error != "" {
			return nil, errors.New(hqError.Error)
		} else {
			return nil, errors.New("unknown error: " + string(bytes))
		}
	} else {
		return &users, nil
	}
}
func (account *Account) AddFriend(userId int) (error) {
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/friends/" + strconv.Itoa(userId) + "/requests", nil)
	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("authorization", "Bearer "+account.AccessToken)
	req.Header.Add("user-agent", "okhttp/3.8.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)

	hqError := HQError{}

	if json.Unmarshal(bytes, &hqError); hqError.Error != "" {
		return errors.New(hqError.Error)
	} else {
		return nil
	}
}
func (account *Account) RemoveFriend(userId int) (error) {
	req, _ := http.NewRequest("DELETE", "https://api-quiz.hype.space/friends/" + strconv.Itoa(userId) + "/requests", nil)
	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("authorization", "Bearer "+account.AccessToken)
	req.Header.Add("user-agent", "okhttp/3.8.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)

	hqError := HQError{}

	if json.Unmarshal(bytes, &hqError); hqError.Error != "" {
		return errors.New(hqError.Error)
	} else {
		return nil
	}
}

func (account *Account) UpdateTokens() (*Tokens, error) {
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/tokens/", strings.NewReader(`{"token":"`+account.LoginToken+`"}`))
	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("authorization", "Bearer "+account.AccessToken)
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("user-agent", "okhttp/3.8.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)

	var tokens Tokens
	var hqError HQError

	if json.Unmarshal(bytes, &tokens); tokens.Username == "" {
		if json.Unmarshal(bytes, &hqError); hqError.Error != "" {
			return nil, errors.New(hqError.Error)
		} else {
			return nil, errors.New("unknown error: " + string(bytes))
		}
	} else {
		return &tokens, nil
	}
}
func (account *Account) Me() (*Me, error) {
	req, _ := http.NewRequest("GET", "https://api-quiz.hype.space/users/me", nil)
	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("authorization", "Bearer "+account.AccessToken)
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("user-agent", "okhttp/3.8.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)

	var me Me
	var hqError HQError

	if json.Unmarshal(bytes, &me); me.Username == "" {
		if json.Unmarshal(bytes, &hqError); hqError.Error != "" {
			return nil, errors.New(hqError.Error)
		} else {
			return nil, errors.New("unknown error: " + string(bytes))
		}
	} else {
		return &me, nil
	}
}
func (account *Account) Cashout(paypal string) error {
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/users/me/payouts", strings.NewReader(`{"email":"`+paypal+`"}`))
	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("authorization", "Bearer "+account.AccessToken)
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("user-agent", "okhttp/3.8.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)

	hqError := HQError{}

	if json.Unmarshal(bytes, &hqError); hqError.Error != "" {
		return errors.New(hqError.Error)
	} else {
		return nil
	}
}

func (account *Account) RequestAWS() (*AWSCredentials, error) {
	req, _ := http.NewRequest("GET", "https://api-quiz.hype.space/credentials/s3", nil)
	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("authorization", "Bearer "+account.AccessToken)
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("user-agent", "okhttp/3.8.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)

	var creds AWSCredentials
	var hqError HQError

	if json.Unmarshal(bytes, &creds); creds.AccessKeyID == "" {
		if json.Unmarshal(bytes, &hqError); hqError.Error != "" {
			return nil, errors.New(hqError.Error)
		} else {
			return nil, errors.New("unknown error: " + string(bytes))
		}
	} else {
		return &creds, nil
	}
}
func (creds *AWSCredentials) UploadAvatarAWS(filename string, image []byte) error {
	req, _ := http.NewRequest("PUT", "https://hypespace-quiz.s3.amazonaws.com/avatars/"+filename, bytes.NewReader(image))
	req.Header.Add("Content-Type", "image/jpeg")
	req.Header.Add("Host", "hypespace-quiz.s3.amazonaws.com")

	awsauth.Sign(req, awsauth.Credentials{Expiration: creds.Expiration, AccessKeyID: creds.AccessKeyID, SecretAccessKey: creds.SecretKey, SecurityToken: creds.SessionToken})

	_, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
func (account *Account) ChangeAvatarAWS(filename string) (*ProfileChange, error) {
	payload := strings.NewReader("{\"avatarUrl\":\"" + filename + "\"}")
	req, _ := http.NewRequest("PUT", "https://api-quiz.hype.space/users/me/avatarUrl", payload)

	req.Header.Add("Authorization", "Bearer "+account.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)

	var pChange ProfileChange
	var hqError HQError

	if json.Unmarshal(bytes, &pChange); pChange.AvatarURL == "" {
		if json.Unmarshal(bytes, &hqError); hqError.Error != "" {
			return nil, errors.New(hqError.Error)
		} else {
			return nil, errors.New("unknown error: " + string(bytes))
		}
	} else {
		return &pChange, nil
	}
}

func Schedule(bearer string) (HQSchedule, error) {
	req, _ := http.NewRequest("GET", "https://api-quiz.hype.space/shows/now?type=hq", nil)
	req.Header.Set("authorization", "Bearer " + bearer)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return HQSchedule{}, err
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var hqError HQError
	if json.Unmarshal(bytes, &hqError); hqError.Error != "" {
		return HQSchedule{}, errors.New(hqError.Error)
	}

	var schedule HQSchedule
	json.Unmarshal(bytes, &schedule)

	return schedule, nil
}

func Connect(broadcastId int, bearer string) (*Game, error) {
	headers := http.Header{}
	headers.Add("authorization", "Bearer " + bearer)

	c, _, err := websocket.DefaultDialer.Dial("wss://ws-quiz.hype.space/ws/"+strconv.Itoa(broadcastId), headers)
	if err != nil {
		return nil, err
	}

	return &Game{Conn: c}, nil
}
func Debug() (*Game, error) {
	c, _, err := websocket.DefaultDialer.Dial("wss://hqecho.herokuapp.com/", nil)
	if err != nil {
		return nil, err
	}

	return &Game{Conn: c}, nil
}

func (game *Game) SendPing() {
	game.Conn.WriteMessage(websocket.PingMessage, nil)
}

func (game *Game) SendSubscribe(broadcastID int) error {
	return game.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"subscribe","broadcastId":`+strconv.Itoa(broadcastID)+`}`))
}
func (game *Game) SendAnswer(broadcastID, questionID, answerID int) error {
	return game.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"answer","broadcastId":`+strconv.Itoa(broadcastID)+`,"questionId":`+strconv.Itoa(questionID)+`,"answerId":`+strconv.Itoa(answerID)+`}`))
}
func (game *Game) SendExtraLife(broadcastID, questionID int) error {
	return game.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"useExtraLife","broadcastId":`+strconv.Itoa(broadcastID)+`,"questionId":`+strconv.Itoa(questionID)+`}`))
}

func (game *Game) ParseBroadcastStats(bytes []byte) (*BroadcastStats) {
	var broadcastStats BroadcastStats
	json.Unmarshal(bytes, &broadcastStats)

	if broadcastStats.Type == "broadcastStats" {
		return &broadcastStats
	}

	return nil
}
func (game *Game) ParseChatMessage(bytes []byte) (*ChatMessage) {
	var chatMessage ChatMessage
	json.Unmarshal(bytes, &chatMessage)

	if chatMessage.Type == "interaction" && chatMessage.ItemID == "chat" {
		return &chatMessage
	}

	return nil
}
func (game *Game) ParseQuestion(bytes []byte) (*Question) {
	var question Question
	json.Unmarshal(bytes, &question)

	if question.Type == "question" && len(question.Answers) != 0 {
		return &question
	}

	return nil
}
func (game *Game) ParseQuestionSummary(bytes []byte) (*QuestionSummary) {
	var questionSummary QuestionSummary
	json.Unmarshal(bytes, &questionSummary)

	if questionSummary.Type == "questionSummary" {
		return &questionSummary
	}

	return nil
}
func (game *Game) ParseQuestionFinished(bytes []byte) (*QuestionFinished) {
	var questionFinished QuestionFinished
	json.Unmarshal(bytes, &questionFinished)

	if questionFinished.Type == "questionFinished" {
		return &questionFinished
	}

	return nil
}
func (game *Game) ParseQuestionClosed(bytes []byte) (*QuestionClosed) {
	var questionClosed QuestionClosed
	json.Unmarshal(bytes, &questionClosed)

	if questionClosed.Type == "questionClosed" {
		return &questionClosed
	}

	return nil
}
func (game *Game) ParseGameStatus(bytes []byte) (*GameStatus) {
	var gameStatus GameStatus
	json.Unmarshal(bytes, &gameStatus)

	if gameStatus.Type == "gameStatus" {
		return &gameStatus
	}

	return nil
}
