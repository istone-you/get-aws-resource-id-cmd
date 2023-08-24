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

// eipCmd represents the eip command
var eipCmd = &cobra.Command{
	Use:   "eip",
	Short: "Show Elastic IP Allocation ID",
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

		input := &ec2.DescribeAddressesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("instance-id"),
					Values: []*string{aws.String(id)},
				},
			},
		}

		resp, err := ec2Client.DescribeAddresses(input)
		if err != nil {
			fmt.Println("Error describing addresses:", err)
		}

		if len(resp.Addresses) == 0 {
			fmt.Println("紐づいているElastic IPはありません。")
		} else if len(resp.Addresses) > 0 {
			fmt.Println("-----")
			for _, address := range resp.Addresses {
				nameTagFound := false
				for _, tag := range address.Tags {
					if *tag.Key == "Name" {
						fmt.Printf("NameTag: %s\n", *tag.Value)
						nameTagFound = true
						break
					}
				}
				if !nameTagFound {
					fmt.Println("NameTag: -")
				}
				fmt.Println("ElasticIPAddress:", *address.PublicIp)
				fmt.Println("ElasticIPAllocationID:", *address.AllocationId)
				fmt.Println("-----")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(eipCmd)

	eipCmd.Flags().StringVarP(&profile, "profile", "p", "", "Set AWS CLI's profile name")
	eipCmd.Flags().StringVarP(&id, "id", "i", "", "Set EC2 Instance ID")
}
