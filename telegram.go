package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tucnak/telebot"
	"strconv"
)

func telegram(bot *telebot.Bot) {
	bot.Callbacks = make(chan telebot.Callback)
	bot.Messages = make(chan telebot.Message)

	go func() {
		for callback := range bot.Callbacks {
			user := GetUser(callback.Sender)
			if !user.Exists() {
				continue
			} else {
				user.Save()
			}

			text := callback.Data
			text = strings.Replace(text, "user_", "", -1)
			id, err := strconv.Atoi(text)
			if err != nil {
				continue
			}
			target := GetTarget(id)
			if target.ID > 0 && target.UserID == callback.Sender.ID {
				target.Delete()
				bot.SendMessage(callback.Sender, "کاربر مورد نظر از لیست حذف شد", nil)
				continue
			}
			bot.SendMessage(callback.Sender, "این کاربر در لیست شما نیست", nil)
		}
	}()

	go func() {
		for message := range bot.Messages {
			user := GetUser(message.Sender)
			if !user.Exists() {
				user.User = message.Sender
				user.Create()
			} else {
				user.Save()
			}

			sender := message.Sender
			text := message.Text
			text = strings.ToLower(text)
			text = strings.TrimSpace(text)

			if strings.HasPrefix(text, "@") {
				text = text[1:]
			}

			if text == "/start" {
				bot.SendMessage(sender, "سلام ، به ربات نظاره گر تصاویر کاربران خوش آمدید ، با من میتونید هر کاربر رو توی تلگرام زیر نظر بگیرید ، کافیه اسمش رو بهم بدید تا من هر وقت عکس پروفایلش تغییر کرد به شما اطلاع بدم", nil)
			} else if text == "/mylist" {
				list := user.Targets()
				if len(list) == 0 {
					bot.SendMessage(sender, "شما کسی را زیر نظر نگرفته اید", nil)
					continue
				}
				keyboard := make([][]telebot.KeyboardButton, 0)
				for _, u := range list {
					row := make([]telebot.KeyboardButton, 2)
					row[0] = telebot.KeyboardButton{
						Data: fmt.Sprintf("user_%d", u.ID),
						Text: "حذف",
					}
					row[1] = telebot.KeyboardButton{
						URL:  fmt.Sprintf("https://t.me/%s", u.Username),
						Text: u.FullName,
					}
					keyboard = append(keyboard, row)
				}
				bot.SendMessage(sender, "لیست کابران شما به شرح زیر است ، برای مشاهده یا حذف یکی را انتخاب کنید", &telebot.SendOptions{
					ReplyMarkup: telebot.ReplyMarkup{
						InlineKeyboard: keyboard,
					},
				})
			} else {
				if user.HasTargets() == 5 {
					bot.SendMessage(sender, "بیشتر از ۵ نفر نمی توانید زیر نظر بگیرید"+"\n"+"مشاهده لیست زیر نظر گرفته ها /mylist", nil)
					continue
				}
				if user.HasTarget(text) {
					bot.SendMessage(sender, "شما قبلا این کاربر رو زیر نظر گرفتید ، با /mylist میتونید حذف کنید", nil)
				} else {
					c, err := searchUser(text)
					if err != nil {
						log.Printf("error on searchUser %v", err)
						bot.SendMessage(sender, "خطایی رخ داده است", nil)
						continue
					}

					if strings.HasPrefix(c.FullName, "Telegram: ") {
						bot.SendMessage(sender, "کاربر یافت نشد", nil)
						continue
					}

					err = downloadImage(c.ProfileLink, text)
					if err != nil {
						log.Printf("error on downloadImage %v", err)
						bot.SendMessage(sender, "خطایی رخ داده است", nil)
						continue
					}
					filePath := "./images/@" + text + ".jpg"

					hash, err := calculateMD5Hash(filePath)
					if err != nil {
						log.Printf("error on calculateMD5Hash %v", err)
						bot.SendMessage(sender, "خطایی رخ داده است", nil)
						continue
					}
					target := Target{}
					target.ProfileHash = hash
					target.ProfileImage = c.ProfileLink
					target.UserID = sender.ID
					target.FullName = c.FullName
					target.Username = text
					target.Create()

					file, _ := telebot.NewFile(filePath)
					photo := telebot.Photo{}
					photo.Caption = fmt.Sprintf("کاربر %s یافت شد"+"\n"+"این کاربر زیر نظر گرفته شد", c.FullName)
					photo.File = file
					err = bot.SendPhoto(sender, &photo, nil)
					if err != nil {
						log.Printf("error on SendPhoto %v", err)
						bot.SendMessage(sender, "خطایی رخ داده است", nil)
						continue
					}
				}
			}
		}
	}()

	bot.Start(1 * time.Second)
}

func downloadImage(url, username string) error {
	out, err := os.Create("./images/@" + username + ".jpg")
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func calculateMD5Hash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
