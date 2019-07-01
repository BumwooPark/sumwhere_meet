package middlewares

import (
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/labstack/gommon/log"
	"google.golang.org/api/option"
	"os"
	"time"
)

const (
	CHATALERT   = "ChatAlert"
	EVENTALERT  = "EventAlert"
	FRIENDALERT = "FriendAlert"
	MATCHALERT  = "MatchAlert"
)

type AppAdapterInterface interface {
	SendMessage(title, body, token string) error
	SetSubscribe(ctx context.Context, isSubscribe bool, token []string, topic string) (int, error)
	subscribe(ctx context.Context, token []string, topic string) (int, error)
	unSubscribe(ctx context.Context, token []string, topic string) (int, error)
	MustMakeChatRoom(ctx context.Context, userID, toUserID int64) error
}

type FireBaseAppAdapter struct {
	app *firebase.App
}

func NewFireBaseApp() (AppAdapterInterface, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var opt option.ClientOption
	if os.Getenv("RELEASE_SYSTEM") == "kubernetes" {
		opt = option.WithCredentialsFile("/config/galmal-8f900-firebase-adminsdk-zhjsl-f6d034ad3b.json")
	} else {
		opt = option.WithCredentialsFile(dir + "/galmal-8f900-firebase-adminsdk-zhjsl-f6d034ad3b.json")
	}
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	return &FireBaseAppAdapter{
		app: app,
	}, nil
}

func (f *FireBaseAppAdapter) SendMessage(title, body, token string) error {
	ctx := context.Background()
	oneHour := time.Duration(1) * time.Hour
	badge := 0
	m := messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Android: &messaging.AndroidConfig{
			TTL: &oneHour,
			Notification: &messaging.AndroidNotification{
				Title: "",
				Icon:  "",
				Color: "",
				Sound: "",
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Badge: &badge,
				},
			},
		},
		Token:     token,
		Topic:     "",
		Condition: "",
	}

	client, err := f.app.Messaging(ctx)
	if err != nil {
		return err
	}
	_, err = client.Send(ctx, &m)
	if err != nil {
		return err
	}
	return nil
}

func (f *FireBaseAppAdapter) SetSubscribe(ctx context.Context, isSubscribe bool, token []string, topic string) (int, error) {
	if isSubscribe {
		return f.subscribe(ctx, token, topic)
	} else {
		return f.unSubscribe(ctx, token, topic)
	}
}

func (f *FireBaseAppAdapter) subscribe(ctx context.Context, token []string, topic string) (int, error) {
	client, err := f.app.Messaging(ctx)
	if err != nil {
		return 0, err
	}

	res, err := client.SubscribeToTopic(ctx, token, topic)
	if err != nil {
		return res.SuccessCount, err
	}

	if len(res.Errors) != 0 {
		for _, err := range res.Errors {
			log.Error(err.Reason)
		}
	}

	return res.SuccessCount, nil
}

func (f *FireBaseAppAdapter) unSubscribe(ctx context.Context, token []string, topic string) (int, error) {
	client, err := f.app.Messaging(ctx)
	if err != nil {
		return 0, err
	}
	res, err := client.UnsubscribeFromTopic(ctx, token, topic)
	if err != nil {
		return res.SuccessCount, err
	}
	if len(res.Errors) != 0 {
		for _, err := range res.Errors {
			log.Error(err.Reason)
		}
	}
	return res.SuccessCount, nil
}

func (f *FireBaseAppAdapter) MustMakeChatRoom(ctx context.Context, userID, toUserID int64) error {
	client, err := f.app.DatabaseWithURL(ctx, "https://galmal-8f900.firebaseio.com/")
	if err != nil {
		return err
	}

	location := fmt.Sprintf("%d_%d_%s", userID, toUserID, time.Now().UTC().Format("2006-01-02T15:04:05-0700"))
	user := fmt.Sprintf("%d", userID)
	toUser := fmt.Sprintf("%d", toUserID)

	err = client.NewRef("conversations").Child(location).Set(ctx, "")
	if err != nil {
		return err
	}

	err = client.NewRef("users").Child(user).Child("conversations").Child(toUser).Child("location").Set(ctx, location)
	if err != nil {
		return err
	}

	err = client.NewRef("users").Child(toUser).Child("conversations").Child(user).Child("location").Set(ctx, location)
	if err != nil {
		return err
	}

	return nil
}
