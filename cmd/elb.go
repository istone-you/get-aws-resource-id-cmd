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
	Short: "Show ELB LoadBalancer ID",
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

		loadBalancerInput := &elbv2.DescribeLoadBalancersInput{
			Names: []*string{&name},
		}

		loadBalancerResult, err := elbClient.DescribeLoadBalancers(loadBalancerInput)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			os.Exit(1)
		}

		if len(loadBalancerResult.LoadBalancers) == 0 {
			fmt.Println("ロードバランサーが見つかりません。")
			os.Exit(1)
		}

		loadBalancerARN := *loadBalancerResult.LoadBalancers[0].LoadBalancerArn

		fmt.Println()
		if showArn {
			fmt.Println("LoadBalancerArn:", loadBalancerARN)
		} else {
			fmt.Println("LoadBalancerID:", (loadBalancerARN)[strings.LastIndex(loadBalancerARN, "/")+1:])
		}
		fmt.Println()

		listenersInput := &elbv2.DescribeListenersInput{
			LoadBalancerArn: &loadBalancerARN,
		}

		listenersResult, err := elbClient.DescribeListeners(listenersInput)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			os.Exit(1)
		}

		fmt.Println("-----")
		for _, listener := range listenersResult.Listeners {
			fmt.Println("Port:", *listener.Port)
			fmt.Println("Protocol:", *listener.Protocol)
			if showArn {
				fmt.Println("ListenerArn:", *listener.ListenerArn)
			} else {
				fmt.Println("ListenerID:",(*listener.ListenerArn)[strings.LastIndex(*listener.ListenerArn, "/")+1:])
			}

			targetGroupsInput := &elbv2.DescribeTargetGroupsInput{
				LoadBalancerArn: &loadBalancerARN,
			}

			targetGroupsResult, err := elbClient.DescribeTargetGroups(targetGroupsInput)
			if err != nil {
				fmt.Println("データの取得に失敗しました:", err)
				os.Exit(1)
			}

			for _, tg := range targetGroupsResult.TargetGroups {
				for _, rule := range listener.DefaultActions {
					if rule.TargetGroupArn != nil && *rule.TargetGroupArn == *tg.TargetGroupArn {
						fmt.Println("TargetGroupName:", *tg.TargetGroupName)
						if showArn {
							fmt.Println("TargetGroupArn:", *tg.TargetGroupArn)
						} else {
							fmt.Println("TargetGroupID:",(*tg.TargetGroupArn)[strings.LastIndex(*tg.TargetGroupArn, "/")+1:])
						}
					}
				}
			}

			fmt.Println("-----")
		}
	},
}

func init() {
	rootCmd.AddCommand(elbCmd)

	elbCmd.Flags().StringVarP(&profile, "profile", "p", "", "Set AWS CLI's profile name")
	elbCmd.Flags().StringVarP(&name, "name", "n", "", "Set ELB LoadBalancer name")
	elbCmd.Flags().BoolVarP(&showArn, "arn", "a", false, "Show ELB LoadBalancer Arn")
}
