package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/keighl/metabolize"
	"github.com/tucnak/telebot"
)

func crawler(bot *telebot.Bot) {
	for {
		targets := GetAllTargets()
		for _, target := range targets {
			c, err := searchUser(target.Username)
			if err != nil {
				log.Printf("error on crawler %v for %s", err, target.Username)
				continue
			}
			if c.ProfileLink != target.ProfileImage {
				err = downloadImage(c.ProfileLink, target.Username)
				if err != nil {
					log.Printf("error on downloadImage %v for %s", err, target.Username)
					continue
				}
				filePath := "./images/@" + target.Username + ".jpg"
				hash, err := calculateMD5Hash(filePath)
				if err != nil {
					log.Printf("error on downloadImage %v for %s", err, target.Username)
					continue
				}
				if hash != target.ProfileHash {
					file, _ := telebot.NewFile(filePath)
					photo := telebot.Photo{}
					photo.Caption = fmt.Sprintf("عکس پروفایل %s تغییر کرد", c.FullName)
					photo.File = file
					bot.SendPhoto(telebot.User{ID: target.UserID}, &photo, nil)
				}
				target.ProfileHash = hash
				target.ProfileImage = target.ProfileImage
			}
			target.Save()

			time.Sleep(500 * time.Millisecond)
		}

		time.Sleep(5 * time.Second)
	}
}

type CrawlerUser struct {
	ProfileLink string `meta:"og:image"`
	FullName    string `meta:"og:title"`
}

func searchUser(username string) (CrawlerUser, error) {
	res, _ := http.Get(fmt.Sprintf("https://t.me/%s", username))
	data := new(CrawlerUser)

	err := metabolize.Metabolize(res.Body, data)
	if err != nil {
		return CrawlerUser{}, err
	}

	return *data, nil
}
