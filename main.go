package main

import (
	"fmt"
	"os"
	"strings"

	twilio "github.com/twilio/twilio-go"
	twilio_client "github.com/twilio/twilio-go/client"
	twilio_api "github.com/twilio/twilio-go/rest/api/v2010"
)

func mustEnv(env string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		panic(fmt.Sprintf("missing %q environment variable", env))
	}
	return value
}

func main() {
	var (
		twilioAccount      = mustEnv("TWILIO_ACCOUNT_SID")
		twilioAPIKey       = mustEnv("TWILIO_API_KEY_SID")
		twilioAPIKeySecret = mustEnv("TWILIO_API_KEY_SECRET")
	)

	twilioClient := twilio_client.Client{
		Credentials: twilio_client.NewCredentials(
			twilioAPIKey,
			twilioAPIKeySecret,
		),
	}
	twilioClient.SetAccountSid(twilioAccount)

	client := twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Client: &twilioClient,
		},
	)

	maxMessagesToList := 1000

	keys, err := client.Api.ListKey(&twilio_api.ListKeyParams{
		Limit: &maxMessagesToList,
	})
	if err != nil {
		panic(err)
	}

	for _, k := range keys {
		fmt.Println("created @", *k.DateCreated, " => ", *k.FriendlyName, " sid ", *k.Sid)

		if strings.HasPrefix(*k.FriendlyName, "vault-secrets") {
			err := client.Api.DeleteKey(*k.Sid, &twilio_api.DeleteKeyParams{})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("deleted", *k.FriendlyName, "sid", *k.Sid)
			}
		}
	}
}
