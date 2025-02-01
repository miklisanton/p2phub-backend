package bot

import (
	"encoding/json"
	"fmt"
	"p2pbot/internal/services"
	"p2pbot/internal/utils"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (bot *Bot) HandleNotification(msg amqp.Delivery) {
	if msg.ContentType == "application/json" {
		var n services.Notification
		if err := json.Unmarshal(msg.Body, &n); err != nil {
			utils.Logger.LogError().Msg(err.Error())
			return
		}
		q, minA, maxA := n.Data.GetQuantity()
		price := n.Data.GetPrice()
		name := n.Data.GetName()
		pms := strings.Join(n.Data.GetPaymentMethods(), ", ")

		template := `Your %s %s advertisement on %s was outbided by %s.
Payment methods: %s.
Quantity: %.2fUSDT.
Min. amount: %.1f%s | Max. amount: %.1f%s.
Price: %.2f%s`
		message := fmt.Sprintf(
			template,
			n.Currency,
			n.Side,
			n.Exchange,
			name,
			pms,
			q,
			minA,
			n.Currency,
			maxA,
			n.Currency,
			price,
			n.Currency)
		bot.SendMessage(n.ChatID, message)
	} else {
		utils.Logger.LogInfo().Msg(string(msg.Body))
	}
}
