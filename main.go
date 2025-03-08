package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig"
)

var labelFilter string

func init() {
	flag.StringVar(&labelFilter, "labelFilter", "Common", "label to filter")
	flag.Parse()
}

func main() {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// APPCONFIGURATION_ENDPOINT=https://YOUR-APPCONFIG-appcs.azconfig.io
	connectionEndpoint := os.Getenv("APPCONFIGURATION_ENDPOINT")
	client, err := azappconfig.NewClient(connectionEndpoint, credential, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR creating azappconfig.NewClient: %s", err)
	}

	any := "*"
	settingsPager := client.NewListSettingsPager(
		azappconfig.SettingSelector{KeyFilter: &any, LabelFilter: &labelFilter, Fields: azappconfig.AllSettingFields()}, nil,
	)

	for settingsPager.More() {
		settingsPage, err := settingsPager.NextPage(context.TODO())

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR getting settingsPager: %s", err)
		}

		for _, setting := range settingsPage.Settings {
			fmt.Printf("%s=%s\n", *setting.Key, *setting.Value)
		}
	}
}
