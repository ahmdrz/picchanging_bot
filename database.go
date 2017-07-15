package main

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tucnak/telebot"
)

var tables []interface{} = []interface{}{
	&User{},
	&Target{},
}

var db *gorm.DB

func ConnectToDatabase(dialect, connString string) error {
	var err error
	db, err = gorm.Open(dialect, connString)
	if err != nil {
		return err
	}
	for _, table := range tables {
		if !db.HasTable(table) {
			db.CreateTable(table)
		} else {
			db.AutoMigrate(table)
		}
	}
	return nil
}

type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type User struct {
	telebot.User
	Model
}

func (u *User) Create() {
	db.Model(u).Create(u)
}

func (u *User) Save() {
	db.Model(u).Save(u)
}

func (u *User) Exists() bool {
	return u.ID > 0
}

func (u *User) HasTarget(username string) bool {
	output := Target{}
	db.Model(&Target{}).Find(&output, "username = ? AND user_id = ?", username, u.ID)
	return output.ID > 0
}

func (u *User) HasTargets() uint {
	count := uint(0)
	db.Model(&Target{}).Where("user_id = ?", u.ID).Count(&count)
	return count
}

func (u *User) Targets() []Target {
	var result []Target = make([]Target, 0)
	db.Model(result).Order("updated_at desc").Where("user_id = ?", u.ID).Scan(&result)
	return result
}

func GetUser(u telebot.User) *User {
	var result *User = new(User)
	db.Model(result).Find(result, "id = ?", u.ID)
	return result
}

func GetAllUsers() []User {
	var result []User = make([]User, 0)
	db.Model(result).Order("updated_at desc").Scan(&result)
	return result
}

type Target struct {
	ID int `gorm:"primary_key,AUTO_INCREMENT"`
	Model
	UserID       int
	User         User
	Username     string
	ProfileImage string
	ProfileHash  string
	FullName     string
}

func (u *Target) Create() {
	db.Model(u).Create(u)
}

func (u *Target) Save() {
	db.Model(u).Save(u)
}

func (u *Target) Delete() {
	db.Model(u).Delete(u)
}

func (u *Target) Exists() bool {
	return u.ID > 0
}

func GetTarget(u int) *Target {
	var result *Target = new(Target)
	db.Model(result).Find(result, "id = ?", u)
	return result
}

func GetAllTargets() []Target {
	var result []Target = make([]Target, 0)
	db.Model(result).Order("updated_at desc").Scan(&result)
	return result
}
