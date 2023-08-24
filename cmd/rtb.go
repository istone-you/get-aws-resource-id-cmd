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

// rtbCmd represents the rtb command
var rtbCmd = &cobra.Command{
	Use:   "rtb",
	Short: "Show VPC Route Table ID",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if id == "" {
			fmt.Print("サブネットIDを入力してください: ")
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

		routeTablesInput := &ec2.DescribeRouteTablesInput{}

		routeTablesOutput, err := ec2Client.DescribeRouteTables(routeTablesInput)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			os.Exit(1)
		}

		for _, rt := range routeTablesOutput.RouteTables {
			for _, assoc := range rt.Associations {
				if assoc.SubnetId != nil && *assoc.SubnetId == id {
					nameTagFound := false
					for _, tag := range rt.Tags {
						if *tag.Key == "Name" {
							fmt.Printf("NameTag: %s\n", *tag.Value)
							nameTagFound = true
							break
						}
					}
					if !nameTagFound {
						fmt.Println("NameTag: -")
					}
					fmt.Printf("RouteTableID: %s\n", *rt.RouteTableId)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(rtbCmd)

	rtbCmd.Flags().StringVarP(&profile, "profile", "p", "", "Set AWS CLI's profile name")
	rtbCmd.Flags().StringVarP(&id, "id", "i", "", "Set Subnet ID")
}
