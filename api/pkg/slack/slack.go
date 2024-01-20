package slack

import (
    "os"
    "log"

    "github.com/slack-go/slack"
)

func SendSlackMessage(channel, message string) error {
    token := os.Getenv("SLACK_TOKEN")
    if token == "" {
        log.Println("SLACK_TOKEN is not set.")
        return errors.New("SLACK_TOKEN is not set")
    }
    
    api := slack.New(token)
    _, _, err := api.PostMessage(
        channel,
        slack.MsgOptionText(message, false),
    )
    
    if err != nil {
        log.Printf("Failed to send message: %s\n", err)
        return err
    }

    log.Printf("Message sent to channel %s", channel)
    return nil
}
