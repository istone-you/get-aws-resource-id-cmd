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

// vgwCmd represents the vgw command
var vgwCmd = &cobra.Command{
	Use:   "vgw",
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

		input := &ec2.DescribeVpnGatewaysInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("attachment.vpc-id"),
					Values: []*string{aws.String(id)},
				},
			},
		}

		resp, err := ec2Client.DescribeVpnGateways(input)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			return
		}

		fmt.Println("-----")

		for _, gw := range resp.VpnGateways {
			nameTagFound := false
			for _, tag := range gw.Tags {
				if *tag.Key == "Name" {
					fmt.Printf("NameTag: %s\n", *tag.Value)
					nameTagFound = true
					break
				}
			}
			if !nameTagFound {
				fmt.Println("NameTag: -")
			}
			fmt.Println("VPNGatewayID:", *gw.VpnGatewayId)
			fmt.Println("-----")
		}
	},
}

func init() {
	rootCmd.AddCommand(vgwCmd)

	vgwCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
	vgwCmd.Flags().StringVarP(&id, "id", "i", "", "VPC ID")
}
