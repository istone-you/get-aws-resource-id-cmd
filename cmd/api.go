/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Show API Gateway ID",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if name == "" {
			fmt.Print("API名を入力してください: ")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)
		}

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config: aws.Config{
				Region: aws.String("ap-northeast-1"),
			},
		}))

		apiClient := apigateway.New(sess)

		apiGatewayID, err := getApiGatewayID(apiClient, name)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			os.Exit(1)
		}

		fmt.Println(apiGatewayID)
	},
}

func getApiGatewayID(svc *apigateway.APIGateway, apiGatewayName string) (string, error) {

	input := &apigateway.GetRestApisInput{}

	resp, err := svc.GetRestApis(input)
	if err != nil {
		return "", err
	}

	for _, api := range resp.Items {
		if aws.StringValue(api.Name) == apiGatewayName {
			return aws.StringValue(api.Id), nil
		}
	}

	return "", fmt.Errorf("指定した名前のAPI Gatewayが見つかりませんでした。")
}

func init() {
	rootCmd.AddCommand(apiCmd)

	apiCmd.Flags().StringVarP(&profile, "profile", "p", "", "Set AWS CLI's profile name")
	apiCmd.Flags().StringVarP(&name, "name", "n", "", "Set API Gateway name")
}
