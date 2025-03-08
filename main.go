package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig"
	"github.com/google/uuid"
)

var labelFilter string

func init() {
	flag.StringVar(&labelFilter, "labelFilter", "Common", "label to filter")
	flag.Parse()
}

// for example:
// go run main.go -labelFilter IT-Dev-VSJohnKnutson,Common  # limited to 5 comme separated items

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
	// more info on label filters:
	// https://learn.microsoft.com/en-us/azure/azure-app-configuration/concept-key-value#query-key-values
	settingsPager := client.NewListSettingsPager(azappconfig.SettingSelector{
		KeyFilter:   &any,
		LabelFilter: &labelFilter,
		Fields:      azappconfig.AllSettingFields(),
	}, nil)

	for settingsPager.More() {
		settingsPage, err := settingsPager.NextPage(context.TODO())

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR getting settingsPager: %s", err)
		}

		for _, setting := range settingsPage.Settings {

			// need to handle multiline values differentely (with a heredoc)
			if strings.Contains(*setting.Value, "\n") {
				// TODO: use a guid for the heredoc
				eof := uuid.New().String()
				fmt.Printf("%s<<%s\n", *setting.Key, eof)
				for line := range strings.Split(*setting.Value, "\n") {
					fmt.Println(line)
				}
				fmt.Println(eof)
			} else {
				fmt.Printf("%s=%s\n", *setting.Key, *setting.Value)
			}
		}
	}
}
