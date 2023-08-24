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
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/spf13/cobra"
)

// efsCmd represents the efs command
var efsCmd = &cobra.Command{
	Use:   "efs",
	Short: "Show EFS File System ID",
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

		efsClient := efs.New(sess)

		// ファイルシステム ID を取得
		fileSystemID, err := getFileSystemID(efsClient, name)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			return
		}

		fmt.Println(fileSystemID)
	},
}

func getFileSystemID(efsClient *efs.EFS, fileSystemName string) (string, error) {
	result, err := efsClient.DescribeFileSystems(&efs.DescribeFileSystemsInput{})
	if err != nil {
		return "", err
	}

	for _, fileSystem := range result.FileSystems {
		if *fileSystem.Name == fileSystemName {
			return *fileSystem.FileSystemId, nil
		}
	}

	return "", fmt.Errorf("ファイルシステムが見つかりません。")
}

func init() {
	rootCmd.AddCommand(efsCmd)

	efsCmd.Flags().StringVarP(&profile, "profile", "p", "", "Set AWS CLI's profile name")
	efsCmd.Flags().StringVarP(&name, "name", "n", "", "Set EFS FileSystem name")
}
