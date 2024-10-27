package ws

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"fyne.io/fyne/v2"
	"github.com/gorilla/websocket"

	MainScreenGui "github.com/schoeneBiene/g0chat/gui/mainscreen"
	State "github.com/schoeneBiene/g0chat/state"
)

const (
	OpMessage      = 0
	OpLogin        = 1
	OpHeartbeat    = 2
	OpMessageHistory = 3
	OpMemberList   = 4
	OpError        = -1
)

var send func ([]byte) error;
var u = url.URL{Scheme: "wss", Host: "chatws.nin0.dev"}

type SocketMessage struct {
    Op int `json:"op"`
    D map[string]interface{} `json:"d"`
}

type ReceiveMessage struct {
    Op int `json:"op"`
    D struct {
        UserInfo struct {
            Username string `json:"username"`
            Roles int `json:"roles"`
            Id string `json:"id"`
            BridgeMetadata map[string]interface{} `json:"bridgeMetadata"`
        } `json:"userInfo"`
        Content string `json:"content"`
        Timestamp int64 `json:"timestamp"`
        Id string `json:"id"`
        Device string `json:"device"`
    } `json:"d"`
}

type MessageHistory struct {
    Op int `json:"op"`
    D struct {
        History []struct {
            UserInfo struct {
                Username string `json:"username"`
                Roles int `json:"roles"`
                Id string `json:"id"`
                BridgeMetadata map[string]interface{} `json:"bridgeMetadata"`
            } `json:"userInfo"`

            Content string `json:"content"`
            Timestamp int64 `json:"timestamp"`
            Id string `json:"id"`
            Device string `json:"device"`
        };
    } `json:"d"`
}

type MemberList struct {
    Op int `json:"op"`
    D struct {
        Users []struct {
            Id string `json:"id"`
            Username string `json:"username"`
        } `json:"users"`
    }
}

type TokenRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type TokenResponse struct {
    Id string `json:"id"`
    Token string `json:"token"`
}

func requestToken(email, password string) string {
    if(State.Login_Token != "") {
        return State.Login_Token;
    }

    requestJson, err := json.Marshal(&TokenRequest{
        Email: email,
        Password: password,
    });

    if err != nil {
        log.Fatal("Failed to get a token for authentication: ", err);
    }

    resp, err := http.Post("https://chatapi.nin0.dev/api/auth/login", "application/json", bytes.NewBuffer(requestJson));

    if err != nil {
        log.Fatal("Failed to get a token for authentication: ", err);
    }


    body, err := io.ReadAll(resp.Body)

    if err != nil {
        log.Fatal("Failed to get a token for authentication: ", err);
    }

    var tokenRes *TokenResponse = &TokenResponse{};
    err = json.Unmarshal(body, tokenRes);

    if err != nil {
        log.Fatal("Failed to get a token for authentication: ", err);
    }

    return tokenRes.Token;
}

func sendLogin() {
    var msg *SocketMessage;

    if(State.Login_Anon) {
        log.Println("Trying to log in as a guest")

        msg = &SocketMessage{
            Op: OpLogin,
            D: map[string]interface{}{
                "anon": true,
                "username": State.Login_Username,
                "device": "web",
            },
        }
    } else {
        log.Println("Trying to log in with email and password")
        token := requestToken(State.Login_Email, State.Login_Password);

        State.Login_Token = token;

        msg = &SocketMessage{
            Op: OpLogin,
            D: map[string]interface{}{
                "anon": false,
                "token": token,
                "device": "web",
            },
        }
    }

    loginJson, err := json.Marshal(msg);

    if err != nil {
        log.Fatal("Failed to encode login json: ", err);
    }

    if(State.Debug_WS) {
        log.Println("Sending login payload: ", string(loginJson));
    }

    err = send(loginJson);

    if err != nil {
        log.Fatal("Failed to log in: ", err);
    }

    if(State.Login_Anon) {
        fyne.CurrentApp().Preferences().SetString("username", State.Login_Username);
    } else {
        fyne.CurrentApp().Preferences().SetString("token", State.Login_Token);
    }
}

func getRoleText(role int, bridgeMetadata map[string]interface{}) string {
    if(len(bridgeMetadata) > 0) {
        return "BRIDGE";
    }

    switch r := role; r {
        case 1: return "GUEST"; 
        case 2: return "USER";
        case 6: return "BOT";
        case 18: return "MOD";
        default: return "ADMIN";
    }
}

func handleReceiveMessage(socketMsg []byte) {
    var decodedMsg *ReceiveMessage = &ReceiveMessage{};

    err := json.Unmarshal(socketMsg, decodedMsg);

    if err != nil {
        log.Fatal("Got error while decoding receive message json: ", err);
    }

    data := decodedMsg.D;
    user := data.UserInfo;

    go MainScreenGui.AddMessage(user.Username, user.Id, getRoleText(user.Roles, user.BridgeMetadata), data.Content, data.Timestamp);
}

func handleMessageHistory(socketMsg []byte) {
    var decodedMsg *MessageHistory = &MessageHistory{};

    err := json.Unmarshal(socketMsg, decodedMsg);

    if err != nil {
        log.Fatal("Got error while decoding message history json: ", err);
    }

    history := decodedMsg.D.History;

    go func() {
        for _, msg := range history {
            user := msg.UserInfo;

            MainScreenGui.AddMessage(user.Username, user.Id, getRoleText(user.Roles, user.BridgeMetadata), msg.Content, msg.Timestamp);
        } 

        sendLogin();
    }()
}

func handleMemberList(socketMsg []byte) {
    var decodedMsg *MemberList = &MemberList{};

    err := json.Unmarshal(socketMsg, decodedMsg);

    if err != nil {
        log.Fatal("Got error while decoding member list json: ", err);
    }

    users := decodedMsg.D.Users;

    usernames := []string{};

    for _, u := range users {
        usernames = append(usernames, u.Username);
    }

    MainScreenGui.UpdateMembers(usernames);
}

func handleMessage(msg []byte) {
    if(State.Debug_WS) {
        log.Printf("Received: %s", msg);
    }

    var decodedMsg SocketMessage;

    err := json.Unmarshal(msg, &decodedMsg);
    
    if err != nil {
        log.Fatal(err);
    }

    if(decodedMsg.Op == OpMessage) {
        handleReceiveMessage(msg);
    }

    if(decodedMsg.Op == OpHeartbeat) {
        if(State.Debug_WS) {
            log.Println("Sending heartbeat")
        }

        res, err := json.Marshal(SocketMessage{
            Op: 2,
            D: map[string]interface{}{},
        })

        if err != nil {
            log.Fatal(err);
        }

        send(res);
    }

    if(decodedMsg.Op == OpMessageHistory) {
        handleMessageHistory(msg);
    }

    if(decodedMsg.Op == OpMemberList) {
        handleMemberList(msg);
    }
}

func MakeSocketConnection() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatal("Read error:", err)
				return
			}

            go handleMessage(message);        
		}
	}()

    send = func(msg []byte) error {
        if(State.Debug_WS) {
            log.Println("Sending: ", string(msg));
        }

		return c.WriteMessage(websocket.TextMessage, msg)
	}

    State.SendMessage = func(content string) {
        toSendMessage, err := json.Marshal(SocketMessage{
            Op: OpMessage,
            D: map[string]interface{}{
                "content": content,
            },
        });

        if err != nil {
            log.Fatal("Failed to send message: ", err);
        }

        send(toSendMessage);
    }

	for {
		select {
		case <-done:
			log.Println("Read loop exited, closing connection")
			return
		case <-interrupt:
			log.Println("Interrupt signal received, closing connection")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Write close error:", err)
				return
			}

			select {
			case <-done:

			case <-time.After(time.Second):
                
			}
			return
		}
	}
}

