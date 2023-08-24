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
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/spf13/cobra"
)

// elbCmd represents the elb command
var elbCmd = &cobra.Command{
	Use:   "elb",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if name == "" {
			fmt.Print("ロードバランサー名を入力してください: ")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)
		}

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config: aws.Config{
				Region: aws.String("ap-northeast-1"),
			},
		}))

		elbClient := elbv2.New(sess)

		input := &elbv2.DescribeLoadBalancersInput{
			Names: []*string{&name},
		}

		result, err := elbClient.DescribeLoadBalancers(input)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			os.Exit(1)
		}

		if len(result.LoadBalancers) > 0 {
			loadBalancerARN := *result.LoadBalancers[0].LoadBalancerArn
			fmt.Println(loadBalancerARN)
		} else {
			fmt.Println("ロードバランサーが見つかりませんでした。")
		}
	},
}

func init() {
	rootCmd.AddCommand(elbCmd)

	elbCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
	elbCmd.Flags().StringVarP(&name, "name", "n", "", "ELB instance name")
}
