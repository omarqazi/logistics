package datastore

import (
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"github.com/omarqazi/logistics/auth"
	"log"
	"time"
)

var getUserStatement *sql.Stmt
var insertUserStatement *sql.Stmt
var updateUserStatement *sql.Stmt
var deleteUserStatement *sql.Stmt

type User struct {
	Id        string
	PublicKey string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeviceId  string // for push notifications
	Latitude  float64
	Longitude float64
}

func GetUser(userId string) (u *User, err error) {
	query := `
	select id, public_key, createdat, updatedat, device_id,
	ST_Y(location) as longitude, ST_X(location) as latitude
	from users where id = $1
	`
	if getUserStatement == nil {
		getUserStatement, err = Postgres.Prepare(query)
		if err != nil {
			return nil, err
		}
	}

	u = &User{}
	var lat, lon sql.NullFloat64
	err = getUserStatement.QueryRow(userId).Scan(
		&u.Id,
		&u.PublicKey,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeviceId,
		&lat,
		&lon,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if lat.Valid {
		u.Latitude = lat.Float64
	} else {
		u.Latitude = 0.0
	}

	if lon.Valid {
		u.Longitude = lon.Float64
	} else {
		u.Longitude = 0.0
	}
	return
}

func (u User) Location() Point {
	return Point{
		Latitude:  u.Latitude,
		Longitude: u.Longitude,
	}
}

func (u *User) UpdateTimestamps(updateCreate bool) {
	now := time.Now()
	if updateCreate {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
	return
}

func (u User) RSAKey() (k *rsa.PublicKey, e error) {
	k, e = auth.PublicKeyFromString(u.PublicKey)
	return
}

func (u *User) EnsureId() {
	if u.Id == "" {
		u.Id = NewUUID()
	}
}

func (u User) Channel() string {
	return "user-" + u.Id
}

func (u User) ToJSON() (string, error) {
	bytes, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (u User) BroadcastUpdate() (err error) {
	jsonString, err := u.ToJSON()
	if err != nil {
		return err
	}

	err = PublishNotification(u.Channel(), jsonString)
	if err != nil {
		log.Println("Error broadcasting update:", err)
	}
	return
}

func (u *User) Create() (err error) {
	if insertUserStatement == nil {
		insertUserStatement, err = Postgres.Prepare(`
		insert into users (
			id, public_key, createdat, updatedat, device_id, location
		) VALUES (
			$1, $2, now(), now(), $3, ST_GeometryFromText($4, 4326)
		);
		`)

		if err != nil {
			return err
		}
	}

	u.EnsureId()
	u.UpdateTimestamps(true)

	_, err = insertUserStatement.Exec(u.Id, u.PublicKey, u.DeviceId, u.Location())
	u.BroadcastUpdate()
	return
}

func (u *User) Update() (err error) {
	if updateUserStatement == nil {
		updateUserStatement, err = Postgres.Prepare(`
			update users set public_key = $1, updatedat = now(), device_id = $2,
			location = ST_GeometryFromText($3, 4326) where id = $4;
		`)
		if err != nil {
			return err
		}
	}

	u.UpdateTimestamps(false)
	_, err = updateUserStatement.Exec(u.PublicKey, u.DeviceId, u.Location(), u.Id)
	u.BroadcastUpdate()
	return
}

func (u *User) Delete() (err error) {
	if deleteUserStatement == nil {
		deleteUserStatement, err = Postgres.Prepare(`
			delete from users where id = $1;
		`)
		if err != nil {
			return err
		}
	}

	_, err = deleteUserStatement.Exec(u.Id)
	u.BroadcastUpdate()
	return
}

var changeSubscribers = make(map[string][]chan string)

func (u User) ChangeJSON() chan string {
	changeChannel := make(chan string, 10)
	userChannel := "user-*" //u.Channel()
	if channelArray, ok := changeSubscribers[userChannel]; ok {
		changeSubscribers[userChannel] = append(channelArray, changeChannel)
	} else {
		newArray := make([]chan string, 0)
		changeSubscribers[userChannel] = append(newArray, changeChannel)
	}
	return changeChannel
}

func init() {
	go monitorUserChanges()
}

func monitorUserChanges() {
	pubsub, err := Redis.PSubscribe("user-*")
	if err != nil {
		log.Fatalln("Error: could not connect to redis for change notifications")
		return
	}

	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			log.Println(err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		channel := msg.Channel
		payload := msg.Payload

		subscriberKeysToNotify := []string{channel, "user-*"}
		for _, subscriberKey := range subscriberKeysToNotify {
			if subscribers, ok := changeSubscribers[subscriberKey]; ok {
				newSubscriberList := make([]chan string, 0)

				for _, payloadChan := range subscribers {
					select {
					case payloadChan <- payload:
						newSubscriberList = append(newSubscriberList, payloadChan)
					default:
					}
				}

				changeSubscribers[subscriberKey] = newSubscriberList
			}
		}
	}
}
