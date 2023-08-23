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

// vpcCmd represents the vpc command
var vpcCmd = &cobra.Command{
	Use:   "vpc",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if name == "" {
			fmt.Print("VPC名を入力してください: ")
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

		// DescribeVpcsリクエストの作成
		describeVpcsInput := &ec2.DescribeVpcsInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("tag:Name"),
					Values: []*string{
						aws.String(name),
					},
				},
			},
		}

		// VPCの検索
		describeVpcsOutput, err := ec2Client.DescribeVpcs(describeVpcsInput)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			os.Exit(1)
		}

		// VPCが見つかった場合、そのIDを表示
		if len(describeVpcsOutput.Vpcs) > 0 {
			for _, vpc := range describeVpcsOutput.Vpcs {
				vpcID := aws.StringValue(vpc.VpcId)
				fmt.Println(vpcID)
			}
		} else {
			fmt.Println("指定した名前のVPCが見つかりませんでした。")
		}
	},
}

func init() {
	rootCmd.AddCommand(vpcCmd)

	vpcCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
	vpcCmd.Flags().StringVarP(&name, "name", "n", "", "VPC name")
}
