package gohq

import (
	"github.com/gorilla/websocket"
	"time"
)

type Account struct {
	UserID      int    `json:"userId"`
	Username    string `json:"username"`
	Admin       bool   `json:"admin"`
	Tester      bool   `json:"tester"`
	Guest       bool   `json:"guest"`
	AvatarURL   string `json:"avatarUrl"`
	LoginToken  string `json:"loginToken"`
	AccessToken string `json:"accessToken"`
	AuthToken   string `json:"authToken"`
}

type Auth struct {
	Auth Account `json:"auth"`
}

type FriendRequest struct {
	RequestedUser struct {
		AvatarURL string    `json:"avatarUrl"`
		Created   time.Time `json:"created"`
		UserID    int       `json:"userId"`
		Username  string    `json:"username"`
	} `json:"requestedUser"`
	RequestingUser struct {
		AvatarURL string    `json:"avatarUrl"`
		Created   time.Time `json:"created"`
		UserID    int       `json:"userId"`
		Username  string    `json:"username"`
	} `json:"requestingUser"`
	Status string `json:"status"`
}

type Users struct {
	Data []struct {
		UserID          int         `json:"userId"`
		Username        string      `json:"username"`
		AvatarURL       string      `json:"avatarUrl,omitempty"`
		Created         time.Time   `json:"created"`
		Live            bool        `json:"live,omitempty"`
		SubscriberCount int         `json:"subscriberCount,omitempty"`
		LastLive        interface{} `json:"lastLive,omitempty"`
		Featured        bool        `json:"featured,omitempty"`
	} `json:"data"`
	Links struct {
		Next string `json:"next"`
		Prev string `json:"prev"`
		Self string `json:"self"`
	} `json:"links"`
}

type Verification struct {
	CallsEnabled   bool      `json:"callsEnabled"`
	Expires        time.Time `json:"expires"`
	Phone          string    `json:"phone"`
	RetrySeconds   int       `json:"retrySeconds"`
	VerificationID string    `json:"verificationId"`
}

type HQError struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"errorCode"`
}

type Game struct {
	Conn *websocket.Conn
}

func (game *Game) Read() ([]byte, error) {
	_, bytes, err := game.Conn.ReadMessage()
	return bytes, err
}

type HQSchedule struct {
	Active        bool      `json:"active"`
	AtCapacity    bool      `json:"atCapacity"`
	ShowID        int       `json:"showId"`
	ShowType      string    `json:"showType"`
	StartTime     time.Time `json:"startTime"`
	NextShowTime  time.Time `json:"nextShowTime"`
	NextShowPrize string    `json:"nextShowPrize"`
	Upcoming []struct {
		Time  time.Time `json:"time"`
		Prize string    `json:"prize"`
	} `json:"upcoming"`
	Prize int `json:"prize"`
	Broadcast struct {
		BroadcastID   int           `json:"broadcastId"`
		UserID        int           `json:"userId"`
		Title         string        `json:"title"`
		Status        int           `json:"status"`
		State         string        `json:"state"`
		ChannelID     int           `json:"channelId"`
		Created       time.Time     `json:"created"`
		Started       time.Time     `json:"started"`
		Ended         interface{}   `json:"ended"`
		Permalink     string        `json:"permalink"`
		ThumbnailData interface{}   `json:"thumbnailData"`
		Tags          []interface{} `json:"tags"`
		SocketURL     string        `json:"socketUrl"`
		Streams struct {
			Source      string `json:"source"`
			Passthrough string `json:"passthrough"`
			High        string `json:"high"`
			Medium      string `json:"medium"`
			Low         string `json:"low"`
		} `json:"streams"`
		StreamURL         string `json:"streamUrl"`
		StreamKey         string `json:"streamKey"`
		RelativeTimestamp int    `json:"relativeTimestamp"`
		Links struct {
			Self       string `json:"self"`
			Transcript string `json:"transcript"`
			Viewers    string `json:"viewers"`
		} `json:"links"`
	} `json:"broadcast"`
	GameKey       string `json:"gameKey"`
	BroadcastFull bool   `json:"broadcastFull"`
}

type BroadcastStats struct {
	Type          string `json:"type"`
	LikeCount     int    `json:"likeCount"`
	StatusMessage string `json:"statusMessage"`
	ViewerCounts struct {
		Connected int `json:"connected"`
		Playing   int `json:"playing"`
		Watching  int `json:"watching"`
	} `json:"viewerCounts"`
	BroadcastSubscribers []interface{} `json:"broadcastSubscribers"`
	Ts                   time.Time     `json:"ts"`
	Sent                 time.Time     `json:"sent"`
}

type ChatMessage struct {
	Type   string `json:"type"`
	ItemID string `json:"itemId"`
	UserID int    `json:"userId"`
	Metadata struct {
		UserID      int    `json:"userId"`
		Message     string `json:"message"`
		AvatarURL   string `json:"avatarUrl"`
		Interaction string `json:"interaction"`
		Username    string `json:"username"`
	} `json:"metadata"`
	Ts   time.Time `json:"ts"`
	Sent time.Time `json:"sent"`
}

type Question struct {
	Type        string `json:"type"`
	TotalTimeMs int    `json:"totalTimeMs"`
	TimeLeftMs  int    `json:"timeLeftMs"`
	QuestionID  int    `json:"questionId"`
	Question    string `json:"question"`
	Category    string `json:"category"`
	Answers []struct {
		AnswerID int    `json:"answerId"`
		Text     string `json:"text"`
	} `json:"answers"`
	QuestionNumber int       `json:"questionNumber"`
	QuestionCount  int       `json:"questionCount"`
	Ts             time.Time `json:"ts"`
	Sent           time.Time `json:"sent"`
}

type QuestionClosed struct {
	Type       string    `json:"type"`
	QuestionID int       `json:"questionId"`
	Ts         time.Time `json:"ts"`
	Sent       time.Time `json:"sent"`
}

type QuestionFinished struct {
	Type       string    `json:"type"`
	QuestionID int       `json:"questionId"`
	Ts         time.Time `json:"ts"`
	Sent       time.Time `json:"sent"`
}

type QuestionSummary struct {
	AdvancingPlayersCount int `json:"advancingPlayersCount"`
	AnswerCounts []struct {
		Answer   string `json:"answer"`
		AnswerID int    `json:"answerId"`
		Correct  bool   `json:"correct"`
		Count    int    `json:"count"`
	} `json:"answerCounts"`
	EliminatedPlayersCount int       `json:"eliminatedPlayersCount"`
	ExtraLivesRemaining    int       `json:"extraLivesRemaining"`
	ID                     string    `json:"id"`
	Question               string    `json:"question"`
	QuestionID             int       `json:"questionId"`
	SavedByExtraLife       bool      `json:"savedByExtraLife"`
	Sent                   time.Time `json:"sent"`
	Ts                     time.Time `json:"ts"`
	Type                   string    `json:"type"`
	YouGotItRight          bool      `json:"youGotItRight"`
	YourAnswerID           int       `json:"yourAnswerId"`
}

type GameStatus struct {
	CardPlaysRemaining  int         `json:"cardPlaysRemaining"`
	Kind                string      `json:"kind"`
	Prize               string      `json:"prize"`
	InTheGame           bool        `json:"inTheGame"`
	Type                string      `json:"type"`
	QuestionCount       int         `json:"questionCount"`
	ExtraLivesRemaining int         `json:"extraLivesRemaining"`
	CurrentState        interface{} `json:"currentState"`
	Cts                 time.Time   `json:"cts"`
	QuestionNumber      int         `json:"questionNumber"`
	ExtraLives          int         `json:"extraLives"`
	Ts                  time.Time   `json:"ts"`
	Sent                time.Time   `json:"sent"`
}

type AWSCredentials struct {
	AccessKeyID  string    `json:"accessKeyId"`
	SecretKey    string    `json:"secretKey"`
	SessionToken string    `json:"sessionToken"`
	Expiration   time.Time `json:"expiration"`
}

type Tokens struct {
	UserID      int    `json:"userId"`
	Username    string `json:"username"`
	Admin       bool   `json:"admin"`
	Tester      bool   `json:"tester"`
	Guest       bool   `json:"guest"`
	AvatarURL   string `json:"avatarUrl"`
	LoginToken  string `json:"loginToken"`
	AccessToken string `json:"accessToken"`
	AuthToken   string `json:"authToken"`
}

type ProfileChange struct {
	UserID    int       `json:"userId"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatarUrl"`
	Created   time.Time `json:"created"`
}

type Me struct {
	UserID    int       `json:"userId"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatarUrl"`
	Created   time.Time `json:"created"`
	Broadcasts struct {
		Data []interface{} `json:"data"`
	} `json:"broadcasts"`
	Featured        bool     `json:"featured"`
	Voip            bool     `json:"voip"`
	DeviceTokens    []string `json:"deviceTokens"`
	HasPhone        bool     `json:"hasPhone"`
	PhoneNumber     string   `json:"phoneNumber"`
	ReferralURL     string   `json:"referralUrl"`
	Referred        bool     `json:"referred"`
	ReferringUserID int      `json:"referringUserId"`
	HighScore       int      `json:"highScore"`
	GamesPlayed     int      `json:"gamesPlayed"`
	WinCount        int      `json:"winCount"`
	Blocked         bool     `json:"blocked"`
	BlocksMe        bool     `json:"blocksMe"`
	Preferences struct {
		SharingEnabled bool `json:"sharingEnabled"`
	} `json:"preferences"`
	FriendIds []int  `json:"friendIds"`
	Lives     string `json:"lives"`
	Stk       string `json:"stk"`
	Leaderboard struct {
		TotalCents int    `json:"totalCents"`
		Total      string `json:"total"`
		Unclaimed  string `json:"unclaimed"`
		Wins       int    `json:"wins"`
		Rank       int    `json:"rank"`
		Alltime struct {
			Total string `json:"total"`
			Wins  int    `json:"wins"`
			Rank  int    `json:"rank"`
		} `json:"alltime"`
		Weekly struct {
			Total string `json:"total"`
			Wins  int    `json:"wins"`
			Rank  int    `json:"rank"`
		} `json:"weekly"`
	} `json:"leaderboard"`
}
