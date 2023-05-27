package store

import (
	"GoMessenger/models"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Store struct {
	url string
	db  *sql.DB
}

func NewStore(user, pass, address, dbname, ssl string) *Store {
	url := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", user, pass, address, dbname, ssl)
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	return &Store{url, db}
}
func (s *Store) Close() {
	if err := s.db.Close(); err != nil {
		panic(err)
	}
}
func (s *Store) GetMessages(chatId int, limit, offset int) ([]models.Message, error) {
	msgs := make([]models.Message, 0, limit)
	rows, err := s.db.Query("select * from \"messages\" where \"chat_id\" = $1 order by \"stamp\" desc limit $2 offset $3", chatId, limit, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.Id, &msg.ChatId, &msg.From, &msg.Mess, &msg.Stamp); err != nil {
			return msgs, err
		}

		msgs = append(msgs, msg)
	}
	if err = rows.Err(); err != nil {
		return msgs, err
	}
	return msgs, nil
}
func (s *Store) GetUserChats(userId int, limit, offset int) ([]int, error) {
	users := make([]int, 0, limit)
	rows, err := s.db.Query("select \"chat_id\" from \"chat_participants\" where \"user_id\"=$1 limit $2 offset $3", userId, limit, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var chtid int
		var usr_id int
		chtid = 0
		usr_id = 0
		if err := rows.Scan(&chtid); err != nil {
			return users, err
		}
		err2 := s.db.QueryRow("select \"user_id\" from \"chat_participants\" where \"user_id\"!=$1 and \"chat_id\"=$2 ", userId, chtid).Scan(&usr_id)
		if err2 != nil {
			return users, err2
		}
		users = append(users, usr_id)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}
func (s *Store) GetUserById(id int) (*models.User, error) {
	var User models.User
	err := s.db.QueryRow("select * from \"users\" where \"id\"=$1", id).Scan(&User.Id, &User.Username, &User.PubKey, &User.Nonce, &User.NonceSign)
	if err != nil {
		return nil, err
	}
	return &User, nil
}
func (s *Store) GetUser(username string) (*models.User, error) {
	var User models.User
	err := s.db.QueryRow("select * from \"users\" where \"username\"=$1", username).Scan(&User.Id, &User.Username, &User.PubKey, &User.Nonce, &User.NonceSign)
	if err != nil {
		return nil, err
	}
	return &User, nil
}
func (s *Store) AddMessage(mess models.Message) error {
	var id int
	log.Println(mess.ChatId)
	err := s.db.QueryRow("insert into \"messages\" (\"chat_id\",\"from_u\",\"mess\",\"stamp\") values($1,$2,$3,$4) returning \"id\"", mess.ChatId, mess.From, mess.Mess, mess.Stamp).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) AddParticipants(chat int, ids []int) error {
	var id int
	for i := range ids {
		log.Println("INSISE ADD: ", ids, chat, ids[i])
		err := s.db.QueryRow("insert into \"chat_participants\" (chat_id,user_id) values($1,$2) returning \"id\"", chat, ids[i]).Scan(&id)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
	}
	return nil
}
func (s *Store) CreateChat() (int, error) {
	var Id int
	err := s.db.QueryRow("insert into \"chats\" default values returning \"id\"").Scan(&Id)
	log.Println("INSIDE SQL FUNC:    ", Id, err)
	if err != nil {
		return 0, err
	}
	log.Println("INSIDE SQL FUNC:    AFTER ERR     ", Id, err)
	return Id, nil
}
func (s *Store) GetChat(toId int, fromId int) (int, error) {
	var Id int
	err := s.db.QueryRow("select \"cht1\".\"chat_id\" from \"chat_participants\" as \"cht1\", \"chat_participants\" as \"cht2\" where \"cht1\".\"chat_id\"=\"cht2\".\"chat_id\" and \"cht1\".\"user_id\" in ($1) and \"cht2\".\"user_id\" in ($2)", fromId, toId).Scan(&Id)
	if err != nil {
		return 0, err
	}
	return Id, nil
}

func (s *Store) RegisterUser(username, pubKey string) (int, error) {
	var id int
	err := s.db.QueryRow("update \"users\" set \"pubkey\"=$1,\"logins\"=$2 where \"username\"=$3 ", pubKey, 1, username).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) UpdateLogins(username string, logins int) (bool, error) {
	var id int
	err := s.db.QueryRow("update \"users\" set \"logins\"=$1 where \"username\"=$2 returning \"id\"", logins, username).Scan(&id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) AddUser(username string, nonce string, logins int) (int, error) {
	var id int
	err := s.db.QueryRow("insert into \"users\" (\"username\",\"nonce\",\"logins\") values ($1,$2,$3) returning \"id\"", username, nonce, logins).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) DropEverything() {
	s.db.QueryRow("DELETE FROM \"chat_participants\"")
	s.db.QueryRow("DELETE FROM \"messages\"")
	s.db.QueryRow("DELETE FROM \"chats\"")
	s.db.QueryRow("DELETE FROM \"users\"")
	//	s.db.QueryRow("DELETE FROM \"chats\"")
}
