package registry

import (
	"fmt"
	"github.com/jdcloud-api/jdcloud-sdk-go/core"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/containerregistry/apis"
	regClient "github.com/jdcloud-api/jdcloud-sdk-go/services/containerregistry/client"
	"github.com/l10178/eea-cli/config"
	"log"
)

func QueryRegistryToken() {
	cloud := config.GetConfig().Cloud

	crd := core.NewCredentials(cloud.AccessKeyId, cloud.SecretAccessKey)
	client := regClient.NewContainerregistryClient(crd)
	client.SetLogger(core.NewDefaultLogger(core.LogWarn))

	req := apis.NewDescribeAuthorizationTokensRequest(cloud.Region, cloud.Registry)

	resp, err := client.DescribeAuthorizationTokens(req)
	if err != nil {
		log.Fatalf("Get token error: %v", err)
		return
	}
	tokens := resp.Result.AuthorizationTokens
	if len(tokens) > 0 {
		fmt.Println(tokens[len(tokens)-1].ExpiresAt)
		fmt.Println(tokens[len(tokens)-1].LoginCmdLine)
	}
}
