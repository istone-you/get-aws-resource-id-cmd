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
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

// snsCmd represents the sns command
var snsCmd = &cobra.Command{
	Use:   "sns",
	Short: "Show SNS Subscription ID",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if name == "" {
			fmt.Print("トピック名を入力してください: ")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)
		}

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config: aws.Config{
				Region: aws.String("ap-northeast-1"),
			},
		}))

		stsClient := sts.New(sess)
		snsClient := sns.New(sess)

		result, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			os.Exit(1)
		}

		accountId := *result.Account

		topicArn := fmt.Sprintf("arn:aws:sns:ap-northeast-1:%s:%s", accountId, name)

		if showArn {
			fmt.Println("-----")
			fmt.Printf("TopicArn: %s\n", topicArn)
		}

		subscriptions, err := getSubscriptions(snsClient, topicArn)
		if err != nil {
			fmt.Println("サブスクリプションの取得に失敗しました:", err)
			return
		}

		fmt.Println("-----")
		for _, sub := range subscriptions {
			if showArn {
				fmt.Printf("SubscriptionArn: %s\n", *sub.SubscriptionArn)
			} else {
				fmt.Printf("SubscriptionID: %s\n", (*sub.SubscriptionArn)[strings.LastIndex(*sub.SubscriptionArn, ":")+1:])
			}
			fmt.Printf("Protocol: %s\n", *sub.Protocol)
			fmt.Printf("Endpoint: %s\n", *sub.Endpoint)
			fmt.Println("-----")
		}
	},
}

func getSubscriptions(svc *sns.SNS, topicARN string) ([]*sns.Subscription, error) {
	result, err := svc.ListSubscriptionsByTopic(&sns.ListSubscriptionsByTopicInput{
		TopicArn: &topicARN,
	})
	if err != nil {
		return nil, err
	}

	return result.Subscriptions, nil
}

func init() {
	rootCmd.AddCommand(snsCmd)

	snsCmd.Flags().StringVarP(&profile, "profile", "p", "", "Set AWS CLI's profile name")
	snsCmd.Flags().StringVarP(&name, "name", "n", "", "Set SNS Topic name")
	snsCmd.Flags().BoolVarP(&showArn, "arn", "a", false, "Show SNS Topic & Subscription Arn")
}
