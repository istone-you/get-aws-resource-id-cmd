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
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

// rdsCmd represents the rds command
var rdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "Show RDS DB Instance ID",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if name == "" {
			fmt.Print("データベース名を入力してください: ")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)
		}

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config: aws.Config{
				Region: aws.String("ap-northeast-1"),
			},
		}))

		rdsClient := rds.New(sess)

		// DescribeDBInstances操作でDBインスタンスの情報を取得
		input := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: &name,
		}

		result, err := rdsClient.DescribeDBInstances(input)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			return
		}

		// リソースIDを取得
		if len(result.DBInstances) > 0 {
			resourceID := aws.StringValue(result.DBInstances[0].DbiResourceId)
			fmt.Println(resourceID)
		} else {
			fmt.Println("指定した名前のデータベースが見つかりませんでした。",)
		}
	},
}

func init() {
	rootCmd.AddCommand(rdsCmd)

	rdsCmd.Flags().StringVarP(&profile, "profile", "p", "", "Set AWS CLI's profile name")
	rdsCmd.Flags().StringVarP(&name, "name", "n", "", "Set RDS Database name")
}
