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

// ebsCmd represents the ebs command
var ebsCmd = &cobra.Command{
	Use:   "ebs",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if id == "" {
			fmt.Print("インスタンスIDを入力してください: ")
			id, _ = reader.ReadString('\n')
			id = strings.TrimSpace(id)
		}

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config: aws.Config{
				Region: aws.String("ap-northeast-1"),
			},
		}))

		ec2Client := ec2.New(sess)

		resp, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{&id},
		})
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			return
		}

		fmt.Println("--------")
		for _, reservation := range resp.Reservations {
			for _, instance := range reservation.Instances {
				for _, blockDevice := range instance.BlockDeviceMappings {
					if blockDevice.Ebs != nil {
						fmt.Println("EBS Volume ID:", *blockDevice.Ebs.VolumeId)
						fmt.Println("Device Name:", *blockDevice.DeviceName)
						fmt.Println("--------")
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(ebsCmd)

	ebsCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
	ebsCmd.Flags().StringVarP(&id, "id", "i", "", "EC2 Instance ID")
}
