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
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/spf13/cobra"
)

// cognitoCmd represents the cognito command
var cognitoCmd = &cobra.Command{
	Use:   "cognito",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if name == "" {
			fmt.Print("ユーザープール名を入力してください: ")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)
		}

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config: aws.Config{
				Region: aws.String("ap-northeast-1"),
			},
		}))

		cognitoClient := cognitoidentityprovider.New(sess)

		userPoolID, err := getUserPoolID(cognitoClient, name)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			return
		}

		fmt.Println(userPoolID)
	},
}

func getUserPoolID(cognitoClient *cognitoidentityprovider.CognitoIdentityProvider, userPoolName string) (string, error) {
	result, err := cognitoClient.ListUserPools(&cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: aws.Int64(60),
	})
	if err != nil {
		return "", err
	}

	for _, userPool := range result.UserPools {
		if *userPool.Name == userPoolName {
			return *userPool.Id, nil
		}
	}

	return "", fmt.Errorf("ユーザープールが見つかりません。")
}

func init() {
	rootCmd.AddCommand(cognitoCmd)

	cognitoCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
	cognitoCmd.Flags().StringVarP(&name, "name", "n", "", "Fsx name")
}
