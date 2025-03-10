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

var labelFilter,
	appconfigEndpoint string

func init() {
	var defaultLabelFilter,
		defaultAppconfigEndpoint string
	var ok bool

	// setup flag defaults from environment variables
	defaultLabelFilter, ok = os.LookupEnv("ACE_LABEL_FILTER")
	if !ok {
		defaultLabelFilter = "Common"
	}
	defaultAppconfigEndpoint, ok = os.LookupEnv("ACE_APPCONFIG_ENDPOINT")
	if !ok {
		defaultAppconfigEndpoint = ""
	}

	flag.StringVar(&labelFilter, "labelFilter", defaultLabelFilter, "label to filter")
	flag.StringVar(&appconfigEndpoint, "appconfigEndpoint", defaultAppconfigEndpoint, "App Configuration endpoint")
	flag.Parse()
}

func main() {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	const appConfigEndpointNotSetError = "appconfig endpoint must be set via `-appconfigEndpoint` flag, or the `ACE_APPCONFIG_ENDPOINT` environment variable."
	if appconfigEndpoint == "" {
		log.Fatalln("ERROR:", appConfigEndpointNotSetError)
	}

	client, err := azappconfig.NewClient(appconfigEndpoint, credential, nil)
	if err != nil {
		log.Fatalln("ERROR creating azappconfig.NewClient:", err)
	}

	// more info on label filters:
	// https://learn.microsoft.com/en-us/azure/azure-app-configuration/concept-key-value#query-key-values
	any := "*"
	settingsPager := client.NewListSettingsPager(azappconfig.SettingSelector{
		KeyFilter:   &any,
		LabelFilter: &labelFilter,
		Fields:      azappconfig.AllSettingFields(),
	}, nil)

	for settingsPager.More() {
		settingsPage, err := settingsPager.NextPage(context.TODO())
		if err != nil {
			log.Fatalln("ERROR getting settingsPager:", err)
		}

		for _, setting := range settingsPage.Settings {
			if strings.Contains(*setting.Value, "\n") {
				// if value is multiline, use heredoc
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
