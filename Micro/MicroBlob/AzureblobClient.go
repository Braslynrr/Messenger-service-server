package main

import (
	"MessengerService/utils"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureClient struct {
	connectionString string
	Url              string
}

// Connect connects to azure blob server
func (ac *AzureClient) Connect() (client *azblob.Client, err error) {

	client, err = azblob.NewClientFromConnectionString(ac.connectionString, nil)

	return
}

// UpLoadBlob uploads blobs to azure
func (ac *AzureClient) UpLoadBlob(name string, file []byte) (url string, err error) {
	var client *azblob.Client
	client, err = ac.Connect()
	if err == nil {
		ctx := context.Background()
		template := strings.Split(name, ".")
		// it has to rename
		var temp string = name
		_, err = ac.LoadBlob(name)
		for err == nil {

			splittedName := strings.Split(temp, ".")

			if solv := utils.FilterString(splittedName[0], utils.HasNumber); solv != "" {
				var num int = 0
				num, err = strconv.Atoi(solv)

				if err == nil {
					temp = fmt.Sprintf("%s(%d).%s", template[0], num+1, template[1])
				}
			} else {
				temp = fmt.Sprintf("%s(1).%s", template[0], template[1])

			}
			_, err = ac.LoadBlob(temp)
		}
		name = temp

		_, err = client.UploadBuffer(ctx, "images", name, file, nil)
		if err == nil {
			url = fmt.Sprintf("%s/images/%s", ac.Url, name)
		}

	}

	return
}

// LoadBlob loads blobs from azure
func (ac *AzureClient) LoadBlob(name string) (img azblob.DownloadStreamResponse, err error) {
	client, err := ac.Connect()
	if err == nil {
		ctx := context.Background()
		img, err = client.DownloadStream(ctx, "images", name, nil)
	}
	return
}
