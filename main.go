package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig"
	"github.com/google/uuid"
)

var labelFilter string
var appconfigEndpoint string

func init() {
	flag.StringVar(&labelFilter, "labelFilter", "Common", "label to filter")
	flag.StringVar(&appconfigEndpoint, "appconfigEndpoint", "", "App Configuration endpoint")
	flag.Parse()
}

// for example:
// go run main.go -labelFilter IT-Dev-VSJohnKnutson,Common  # limited to 5 comme separated items

func main() {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}

	if appconfigEndpoint == "" {
		log.Fatalf("ERROR: APPCONFIGURATION_ENDPOINT is not set")
	}

	client, err := azappconfig.NewClient(appconfigEndpoint, credential, nil)

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
