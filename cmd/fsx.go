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
	"github.com/aws/aws-sdk-go/service/fsx"
	"github.com/spf13/cobra"
)


// fsxCmd represents the fsx command
var fsxCmd = &cobra.Command{
	Use:   "fsx",
	Short: "Show FSx File System ID",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if name == "" {
			fmt.Print("ファイルシステム名を入力してください: ")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)
		}

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config: aws.Config{
				Region: aws.String("ap-northeast-1"),
			},
		}))

		fsxClient := fsx.New(sess)

		result, err := fsxClient.DescribeFileSystems(&fsx.DescribeFileSystemsInput{
			FileSystemIds: []*string{},
		})
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			return
		}

		var filesystemID string
		for _, fs := range result.FileSystems {
			for _, tag := range fs.Tags {
				if aws.StringValue(tag.Key) == "Name" && aws.StringValue(tag.Value) == name {
					filesystemID = aws.StringValue(fs.FileSystemId)
					break
				}
			}
		}

		if filesystemID != "" {
			fmt.Println(filesystemID)
		} else {
			fmt.Println("ファイルシステムが見つかりません。")
		}
	},
}

func init() {
	rootCmd.AddCommand(fsxCmd)

	fsxCmd.Flags().StringVarP(&profile, "profile", "p", "", "Set AWS CLI's profile name")
	fsxCmd.Flags().StringVarP(&name, "name", "n", "", "Set Fsx name")
}
