package handlers

import (
	"log"

	"github.com/opencars/bot/internal/bot"
)

func somethingWentWrong(msg *bot.Event) {
	if err := msg.Send("Вибач. Щось пішло не так 😢"); err != nil {
		log.Printf("send error: %s", err.Error())
	}
}

func (h OpenCarsHandler) PhotoHandler(msg *bot.Event) {
	photos := *msg.Message.Photo

	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		log.Printf("action error: %s", err.Error())
	}

	if len(photos) < 1 {
		somethingWentWrong(msg)
		return
	}

	// TODO: Think about this code snippet.
	url, err := msg.API.GetFileDirectURL(photos[len(photos)-1].FileID)
	if err != nil {

		somethingWentWrong(msg)
		log.Println(err.Error())
		return
	}
	log.Printf("Photo: %s\n", url)

	results, err := h.client.ALPR().Recognize(url)
	if err != nil {
		somethingWentWrong(msg)
		log.Println(err.Error())
		return
	}

	// If nothing was found, send user notification.
	if len(results) == 0 {
		if err := msg.Send("Номер не знайдено 🤔"); err != nil {
			log.Printf("send error: %s\n", err.Error())
		}
		return
	}

	plates := make([]string, len(results))
	for i := range results {
		plates[i] = results[i].Plate
	}

	// Send number to user.
	if err := msg.Send(plates[0]); err != nil {
		log.Printf("send error: %s\n", err.Error())
	}

	text, err := h.getInfoByNumber(plates[0])
	if err != nil {
		log.Println(err.Error())
	}

	if err := msg.SendHTML(text); err != nil {
		log.Printf("send error: %s\n", err.Error())
	}
}
