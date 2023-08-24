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

// sgCmd represents the sg command
var sgCmd = &cobra.Command{
	Use:   "sg",
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

		sgInput := &ec2.DescribeSecurityGroupsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{&id},
				},
			},
		}

		sgOutput, err := ec2Client.DescribeSecurityGroups(sgInput)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			os.Exit(1)
		}

		fmt.Println("-----")
		for _, sg := range sgOutput.SecurityGroups {
			fmt.Printf("SecurityGroupID: %s\n", *sg.GroupId)
			fmt.Printf("SecurityGroupName: %s\n", *sg.GroupName)
			fmt.Println("-----")
		}
	},
}

func init() {
	rootCmd.AddCommand(sgCmd)

	sgCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
	sgCmd.Flags().StringVarP(&id, "id", "i", "", "VPC ID")
}
