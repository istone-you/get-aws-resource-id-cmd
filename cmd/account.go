/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
			if profile == "" {
				profile = "default"
			}

			sess := session.Must(session.NewSessionWithOptions(session.Options{
				Profile: profile,
				Config: aws.Config{
					Region: aws.String("ap-northeast-1"),
				},
			}))

			stsClient := sts.New(sess)

			result, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
			if err != nil {
				fmt.Println("データの取得に失敗しました:", err)
				os.Exit(1)
			}

			// アカウントIDを表示
			fmt.Println(*result.Account)
		},
}

func init() {
	rootCmd.AddCommand(accountCmd)

	accountCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
}
