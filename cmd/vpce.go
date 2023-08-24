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

// vpceCmd represents the vpce command
var vpceCmd = &cobra.Command{
	Use:   "vpce",
	Short: "Show VPC Endpoint ID",
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

		input := &ec2.DescribeVpcEndpointsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{aws.String(id)},
				},
			},
		}

		resp, err := ec2Client.DescribeVpcEndpoints(input)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			return
		}

		if len(resp.VpcEndpoints) == 0 {
			fmt.Println("VPCエンドポイントはありません")
		} else if len(resp.VpcEndpoints) > 0 {
			fmt.Println("-----")
			for _, endpoint := range resp.VpcEndpoints {
				nameTagFound := false
				for _, tag := range endpoint.Tags {
					if *tag.Key == "Name" {
						fmt.Printf("NameTag: %s\n", *tag.Value)
						nameTagFound = true
						break
					}
				}
				if !nameTagFound {
					fmt.Println("NameTag: -")
				}
				fmt.Println("EndpointID:", *endpoint.VpcEndpointId)
				fmt.Println("-----")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(vpceCmd)

	vpceCmd.Flags().StringVarP(&profile, "profile", "p", "", "Set AWS CLI's profile name")
	vpceCmd.Flags().StringVarP(&id, "id", "i", "", "Set VPC ID")
}
