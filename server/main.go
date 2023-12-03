package main

import (
	"GoMessenger/models"
	"GoMessenger/store"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gookit/config"
	"github.com/gookit/config/yaml"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var two64 = (big.NewInt(0)).Exp(big.NewInt(2), big.NewInt(64), nil)

// var myStore = store.NewStore("postgres", "admin", "localhost", "messengerdb", "disable")
var myStore = &store.Store{}
var key ecdsa.PrivateKey
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var conns = models.NewConnections()

func main() {

	config.AddDriver(yaml.Driver)
	err := config.LoadFiles("thesis_work/server/config.yaml")
	if err != nil {
		fmt.Println(os.Getwd())
		fmt.Println(err.Error())
		panic(err)
	}
	var url, _ = config.String("db")
	myStore = store.NewStoreURL(url)
	//setting keys, etc. Should be in the config probably :/
	//919a3e190688a262da21239706c80a4a5b11f897514bc9dcbae3aa627f72ade5
	//365e1e475d29cafbc09e72e2147633e535f61c6a6b2edf456b08c7d2abec8290
	// X then Y
	{
		configPriv, _ := config.String("server_private_key")
		keyD, _ := big.NewInt(0).SetString(configPriv, 16)
		keyX, keyY := elliptic.P256().ScalarBaseMult(keyD.Bytes())
		//	keyD, _ := (big.NewInt(0)).SetString("___", 16) // Enter server's private key here (a big number). curve - P256.
		//	keyX, _ := (big.NewInt(0)).SetString("___", 16) // Enter server's public key part (X coordinate) here. Should be derived from the private key.
		//	keyY, _ := (big.NewInt(0)).SetString("___", 16) // Enter server's public key part (Y coordinate) here. Should be derived from the private key (or X coordinate).
		key = ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: elliptic.P256(),
				X:     keyX,
				Y:     keyY,
			},
			D: keyD,
		}
	}

	ec := echo.New()
	defer myStore.Close()
	ec.POST("api/preregister/user", getUser)   // api/preregister/
	ec.POST("api/register/user", registerUser) // api/register/
	ec.POST("api/login/user", loginUser)       // api/login
	ec.GET("api/get/user/:name", findUser)     // api/get/user/:name
	ec.GET("api/ws/:name/:sig", openWs)        // api/ws/
	ec.GET("/ping", pong)
	ec.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	defer conns.CloseAllConns()
	ec.Logger.Fatal(ec.Start(":8081"))
}
func pong(c echo.Context) error {
	return c.String(http.StatusOK, "asd")
}
func findUser(c echo.Context) error {
	name := c.Param("name")
	user, err := myStore.GetUser(name)
	if err != nil {
		return err
	}
	type userReturn struct {
		Username string `json:"username"`
		PubKey   string `json:"key"`
	}
	var ret userReturn
	ret.Username = user.Username
	ret.PubKey = user.PubKey
	return c.JSON(http.StatusOK, ret)
}

func loginUser(c echo.Context) error {
	type ServerAuthRequest struct {
		Username  string `json:"Username"`
		UserNonce string `json:"UserNonce"`
	}
	var sar ServerAuthRequest
	if err := c.Bind(&sar); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	sha := sha256.New()
	sha.Write([]byte(sar.UserNonce))

	r, s, err := ecdsa.Sign(rand.Reader, &key, sha.Sum(nil))
	if err != nil {
		log.Println("dropped 2")
		return err
	}
	type ServerAuthAnswer struct {
		Username    string `json:"username"`
		SignatureR  string `json:"r"`
		SignatureS  string `json:"s"`
		RandomNonce string `json:"salt"`
		SignNonce   string `json:"nonce"`
	}
	var user *models.User
	if user, err = myStore.GetUser(sar.Username); err != nil {
		log.Println("1")
		log.Println(err)
		return err
	}

	err = nil
	var saa ServerAuthAnswer
	//bigSalt.SetString("1010101010101010101010111111111010101010101010101010101010101010", 2)
	saa.RandomNonce = user.Nonce
	saa.SignatureR = r.Text(16)
	saa.SignatureS = s.Text(16)
	saa.SignNonce = strconv.Itoa(user.NonceSign)
	saa.Username = sar.Username

	return c.JSON(http.StatusOK, saa)
}

func getUser(c echo.Context) error {
	type ServerAuthRequest struct {
		Username  string `json:"Username"`
		UserNonce string `json:"UserNonce"`
	}
	var sar ServerAuthRequest
	if err := c.Bind(&sar); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	sha := sha256.New()
	sha.Write([]byte(sar.UserNonce))

	r, s, err := ecdsa.Sign(rand.Reader, &key, sha.Sum(nil))
	if err != nil {
		log.Println("dropped 2")
		return err
	}
	type ServerAuthAnswer struct {
		Username    string `json:"username"`
		SignatureR  string `json:"r"`
		SignatureS  string `json:"s"`
		RandomNonce string `json:"salt"`
		SignNonce   string `json:"nonce"`
	}
	if user, err := myStore.GetUser(sar.Username); err != nil && err != sql.ErrNoRows {
		log.Println("ERROR        ", err.Error())
		log.Println(user.Id)
		if user.Id != 0 {
			log.Println("dropped 3      ", user)
			return c.String(http.StatusBadRequest, "User with such username already exists!")
		}
		log.Println(err)
		return err
	}
	err = nil
	var saa ServerAuthAnswer
	bigSalt, err := rand.Int(rand.Reader, two64)
	if err != nil {
		log.Println(err)
		return err
	}
	for bigSalt.BitLen() != 64 {
		bigSalt, err = rand.Int(rand.Reader, two64)
	}

	if err != nil {
		log.Println(err)
		return err
	}
	//bigSalt.SetString("1010101010101010101010111111111010101010101010101010101010101010", 2)
	saa.RandomNonce = bigSalt.Text(16)
	saa.SignatureR = r.Text(16)
	saa.SignatureS = s.Text(16)
	saa.SignNonce = "0"
	saa.Username = sar.Username

	_, err = myStore.AddUser(saa.Username, saa.RandomNonce, 0)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusBadRequest, "Such user already exists")
	}
	return c.JSON(http.StatusOK, saa)
}

func registerUser(c echo.Context) error {
	type UserRegistrationRequest struct {
		Username   string `json:"Username"`
		SignatureR string `json:"r"`
		SignatureS string `json:"s"`
		PubKeyX    string `json:"x"`
		PubKeyY    string `json:"y"`
	}
	var urr UserRegistrationRequest
	if err := c.Bind(&urr); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	UserX, _ := (big.NewInt(0)).SetString(urr.PubKeyX, 10)
	UserY, _ := (big.NewInt(0)).SetString(urr.PubKeyY, 10)
	//UserX.SetString(urr.PubKeyX, 10)
	//UserY.SetString(urr.PubKeyY, 10)
	UserPk := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     UserX,
		Y:     UserY,
	}
	UserR, _ := (big.NewInt(0)).SetString(urr.SignatureR, 10)
	UserS, _ := (big.NewInt(0)).SetString(urr.SignatureS, 10)
	//	UserR.SetString(urr.SignatureR, 10)
	//	UserS.SetString(urr.SignatureS, 10)
	log.Println("Signature")
	log.Println(urr.SignatureR)
	log.Println(urr.SignatureS)
	log.Println("PK")
	log.Println(urr.PubKeyX)
	log.Println(urr.PubKeyY)
	sha := sha256.New()
	sha.Write([]byte("0"))
	if !ecdsa.Verify(&UserPk, sha.Sum(nil), UserR, UserS) {
		log.Println("Signature")
		log.Println(urr.SignatureR)
		log.Println(urr.SignatureS)
		log.Println("PK")
		log.Println(urr.PubKeyX)
		log.Println(urr.PubKeyY)
		log.Println("wrong signature")
		return echo.NewHTTPError(http.StatusBadRequest, "Wrong Signature")
	}

	if _, err := myStore.RegisterUser(urr.Username, UserX.Text(16)+"_"+UserY.Text(16)); err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return err
	}

	return c.String(http.StatusOK, "")
}

func openWs(c echo.Context) error {
	name := c.Param("name")
	sig := c.Param("sig")
	sigRS := strings.Split(sig, "_")
	if len(sigRS) != 2 {
		log.Println("openWs: wrong sig =>      ", name, sig, sigRS)
		return errors.New("wrong signature")
	}
	user, err := myStore.GetUser(name)
	if err != nil {
		return err
	}
	UPubKey := strings.Split(user.PubKey, "_")
	if len(UPubKey) != 2 {
		log.Println("openWs: UPubKey Length problem =>    ", UPubKey, user)
	}
	UPubX, _ := (big.NewInt(0)).SetString(UPubKey[0], 16)
	UPubY, _ := (big.NewInt(0)).SetString(UPubKey[1], 16)
	UserPK := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     UPubX,
		Y:     UPubY,
	}
	USigR, _ := (big.NewInt(0)).SetString(sigRS[0], 10)
	USigS, _ := (big.NewInt(0)).SetString(sigRS[1], 10)
	sha := sha256.New()
	sha.Write([]byte(strconv.Itoa(user.NonceSign)))
	if !ecdsa.Verify(&UserPK, sha.Sum(nil), USigR, USigS) {
		log.Println("openWs: signature wrong =>     ", user)
		return errors.New("openWS: sig wrong")
	}
	err = nil
	_, err = myStore.UpdateLogins(name, user.NonceSign+1)
	if err != nil {
		log.Println("update wrong ")
		log.Println(err)
		return err
	}
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	conn := models.NewSocket(ws)
	conns.AddSocket(name, conn)
	log.Println("added ", name)
	go SocketReadLogic(name, conn)
	go SocketLogic(name, conn)
	return nil
}
func SocketReadLogic(username string, soc *models.Socket) {
	for {
		log.Println("Reading in ReadLogic     ", username)
		var packet models.Packet
		if err := soc.Conn.ReadJSON(&packet); err != nil {
			log.Println("error while reading message from the user ", err)
			return
		}
		soc.Out <- packet
	}
}
func SocketLogic(username string, soc *models.Socket) {
	myUser, err123 := myStore.GetUser(username)
	if err123 != nil {
		log.Println("myUser ", err123)
	}
	for {
		log.Println(username, " in for")
		select {
		case <-soc.Quit:
			return
		case packet := <-soc.In:
			log.Println("received message ", packet)
			//add to db here --- no
			err := soc.Conn.WriteJSON(packet)
			if err != nil {
				log.Println("username ws error on soc.In:  ", err)
			}
		case packet := <-soc.Out:
			log.Println("received pakcet from frontend user ", username, packet.To)
			log.Println("packet   ", packet)
			switch packet.Type {
			case 0:
				if packet.Message != "" {
					log.Println(username, " 1")
					if strings.Compare(username, packet.From) != 0 {
					} else {
						log.Println(username, " 2")
						toUser, err := myStore.GetUser(packet.To)
						if err != nil {
							log.Println("WS ToUser ", err, "   my user - ", username)
							continue
						}
						log.Println(username, " 3")
						chatId, err := myStore.GetChat(myUser.Id, toUser.Id)

						if err == sql.ErrNoRows {
							log.Println(username, " 4")
							chatId, err2 := myStore.CreateChat()
							if err2 != nil {

								log.Println("ws ", username, " creating chat ", chatId)
							}
							log.Println(username, " 5")
							err2 = nil
							if err2 = myStore.AddParticipants(chatId, []int{myUser.Id, toUser.Id}); err2 != nil {
								log.Println(username, " 6")
								log.Println("ws ", username, " adding participants ", chatId, " from ", toUser.Id, "err ", err)
							}
							log.Println(username, " 7")
						} else if err != nil && err != sql.ErrNoRows {
							log.Println(username, " 8")
							log.Println("ws chatid get ", err)
							continue
						}
						log.Println(username, " 0")
						chatId, _ = myStore.GetChat(myUser.Id, toUser.Id)
						mess := models.Message{
							Id:     0,
							ChatId: chatId,
							From:   myUser.Id,
							Mess:   packet.Message,
							Stamp:  time.Now().Unix(),
						}
						log.Println(username, " 10")
						err = nil
						err = myStore.AddMessage(mess)
						if err != nil {
							log.Println(username, " 11")
							log.Println("add mess ", err)
							continue
						}
						log.Println(username, " 12")
						log.Println("sending message")
						sendiblePacket := models.SendiblePacket{
							Type:    packet.Type,
							From:    packet.Message,
							To:      packet.To,
							Message: packet.Message,
							Stamp:   mess.Stamp,
						}
						conns.SendMessage(sendiblePacket)
						log.Println(username, " 13")
					}
				}
				break
			case 100:
				log.Println("packet in 100")
				i, _ := strconv.Atoi(packet.To)
				// magic 5 - how many chats to receive
				usr, err := myStore.GetUser(packet.From)
				if err != nil {
					log.Println(err)
					break
				}
				log.Println("packet in 100________1")
				type userReturn struct {
					Username string `json:"username"`
					PubKey   string `json:"key"`
				}
				chats, err := myStore.GetUserChats(usr.Id, 5, i)
				if err != nil {

					log.Println(err)
					break
				}
				log.Println("packet in 100________2")
				chatsReturn := make([]userReturn, 0, 5)
				for i := range chats {
					usr2, err2 := myStore.GetUserById(chats[i])
					if err2 != nil {
						log.Println(err2)
						break
					}
					var ret userReturn
					ret.Username = usr2.Username
					ret.PubKey = usr2.PubKey
					chatsReturn = append(chatsReturn, ret)
				}
				log.Println("packet in 100________3")
				js, err3 := json.Marshal(chatsReturn)
				if err3 != nil {
					log.Println("3   ", err3)
					break
				}
				log.Println("packet in 100________4")
				var myMess models.Packet
				myMess.To = packet.From
				myMess.Type = 101
				myMess.Message = string(js)
				if len(chatsReturn) < 5 {
					myMess.Type = 102
				}
				if len(chatsReturn) != 0 {
					err := soc.Conn.WriteJSON(myMess)
					if err != nil {
						log.Println("username ws error on soc.In:  ", err)
					}
				}
				log.Println("packet in 100________5")
				break
			case 110:
				log.Println("packet in 110")
				// magic 5 - how many chats to receive
				usr, err := myStore.GetUser(packet.To)
				i, _ := strconv.Atoi(packet.Message)
				if err != nil {
					log.Println(err)
					break
				}
				me, err := myStore.GetUser(packet.From)
				if err != nil {
					log.Println(err)
					break
				}
				log.Println("packet in 110________1")

				chat, err := myStore.GetChat(usr.Id, me.Id)
				if err != nil {

					log.Println(err)
					break
				}
				log.Println("packet in 110________2")
				type messReturn struct {
					From    string `json:"username"`
					Message string `json:"mess"`
					Stamp   int64  `json:"stamp"`
				}
				storyReturn := make([]messReturn, 0, 5)
				messages, err2 := myStore.GetMessages(chat, 5, i)
				if err2 != nil {
					log.Println(err2)
					break
				}
				for i := range messages {
					var mr messReturn
					if messages[i].From == usr.Id {
						mr.From = usr.Username
					} else {
						mr.From = packet.From
					}
					mr.Message = messages[i].Mess
					mr.Stamp = messages[i].Stamp
					storyReturn = append(storyReturn, mr)
				}
				log.Println("packet in 110________3")
				js, err3 := json.Marshal(storyReturn)
				if err3 != nil {
					log.Println("3   ", err3)
					break
				}
				log.Println("packet in 110________4")
				var myMess models.Packet
				myMess.To = packet.From
				myMess.Type = 111
				myMess.Message = string(js)
				if len(storyReturn) < 5 {
					myMess.Type = 112
				}
				if len(storyReturn) != 0 {
					err := soc.Conn.WriteJSON(myMess)
					if err != nil {
						log.Println("username ws error on soc.In:  ", err)
					}
				}
				log.Println("packet in 110________5")
				break
			case 999:
				usr, err := myStore.GetUser(packet.From)
				if err != nil {
					log.Println("del acc get usr, case 999:      ", err)
				}
				err = myStore.DeleteAccount(usr.Id)
				if err != nil {
					log.Println("del acc, case 999:      ", err)
				}

			}
			packet.Message = ""
			log.Println("we are here ", username)
		}
	}

}

/*
username -> loginned
user -> server :
	give me my last 5 chats (select * from chats, where username = me max = 5)
							(select participants from ch

*/
