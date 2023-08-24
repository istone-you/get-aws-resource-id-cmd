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

// igwCmd represents the igw command
var igwCmd = &cobra.Command{
	Use:   "igw",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if id == "" {
			fmt.Print("VPC IDを入力してください: ")
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

		igwInput := &ec2.DescribeInternetGatewaysInput{}

		igwOutput, err := ec2Client.DescribeInternetGateways(igwInput)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			os.Exit(1)
		}

		for _, igw := range igwOutput.InternetGateways {
			for _, attachment := range igw.Attachments {
				if attachment.VpcId != nil && *attachment.VpcId == id {
					fmt.Println(*igw.InternetGatewayId)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(igwCmd)

	igwCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
	igwCmd.Flags().StringVarP(&id, "id", "i", "", "VPC ID")
}
