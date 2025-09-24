package handlers

import (
	"encoding/json"
	"fmt"
	"jobby/internals/cache"
	"jobby/internals/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

const NotificationChannel = "notifications"

func PublishNotification(userID uint, userType, message string) {
	notif := models.Notifications{
		UserID:   userID,
		UserType: userType,
		Message:  message,
	}
	data, _ := json.Marshal(notif)

	err := cache.Rdb.Publish(cache.Ctx, NotificationChannel, data).Err()
	if err != nil {
		log.Printf("Faild to publish notifications: %v", err)
	} else {
		log.Printf("Notifications published to channel: %s", message)
	}
}
func SubscribeToNotifications(c *fiber.Ctx) error {
	subscriber := cache.Rdb.Subscribe(cache.Ctx, NotificationChannel)

	ch := subscriber.Channel()

	go func() {
		for msg := range ch {
			fmt.Printf("New Notifications: %s \n", msg.Payload)
		}
	}()

	return c.JSON(fiber.Map{
		"status": "subscribed to notifications",
	})
}
