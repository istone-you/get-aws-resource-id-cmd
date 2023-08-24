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
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)



// instanceCmd represents the instance command
var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "EC2 Instance ID",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if name == "" {
			fmt.Print("インスタンス名を入力してください: ")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)
		}

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config: aws.Config{
				Region: aws.String("ap-northeast-1"),
			},
		}))

		ec2Client := ec2.New(sess)

		input := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("tag:Name"),
					Values: []*string{
						aws.String(name),
					},
				},
				{
					Name: aws.String("instance-state-name"),
					Values: []*string{
						aws.String("running"),
					},
				},
			},
		}

		result, err := ec2Client.DescribeInstances(input)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			return
		}

		if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
			fmt.Println("指定した名前のインスタンスが見つかりませんでした。")
			return
		}

		for _, reservation := range result.Reservations {
			for _, instance := range reservation.Instances {
				instanceID := aws.StringValue(instance.InstanceId)
				fmt.Println(instanceID)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(ec2Cmd)

	ec2Cmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
	ec2Cmd.Flags().StringVarP(&name, "name", "n", "", "EC2 instance name")
}