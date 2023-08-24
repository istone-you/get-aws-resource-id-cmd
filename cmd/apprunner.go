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
	"github.com/aws/aws-sdk-go/service/apprunner"
	"github.com/spf13/cobra"
)

// apprunnerCmd represents the apprunner command
var apprunnerCmd = &cobra.Command{
	Use:   "apprunner",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if name == "" {
			fmt.Print("サービス名を入力してください: ")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)
		}

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config: aws.Config{
				Region: aws.String("ap-northeast-1"),
			},
		}))

		apprunnerClient := apprunner.New(sess)

		listServicesInput := &apprunner.ListServicesInput{}
		listServicesOutput, err := apprunnerClient.ListServices(listServicesInput)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			os.Exit(1)
		}

		// サービスARNを取得
		var serviceARN string
		for _, serviceSummary := range listServicesOutput.ServiceSummaryList {
			if *serviceSummary.ServiceName == name {
				serviceARN = *serviceSummary.ServiceArn
				break
			}
		}

		if serviceARN == "" {
			fmt.Println("サービスが見つかりません。")
			os.Exit(1)
		}

		fmt.Println(serviceARN)
	},
}

func init() {
	rootCmd.AddCommand(apprunnerCmd)

	apprunnerCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
	apprunnerCmd.Flags().StringVarP(&name, "name", "n", "", "App Runner Service name")
}
