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
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/spf13/cobra"
)

// kmsCmd represents the kms command
var kmsCmd = &cobra.Command{
	Use:   "kms",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if name == "" {
			fmt.Print("キーエイリアスを入力してください: ")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)
		}

		name = fmt.Sprint("alias/", name)

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config: aws.Config{
				Region: aws.String("ap-northeast-1"),
			},
		}))

		kmsClient := kms.New(sess)

		describeInput := &kms.DescribeKeyInput{
			KeyId: &name,
		}
		describeOutput, err := kmsClient.DescribeKey(describeInput)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			return
		}

		fmt.Println(*describeOutput.KeyMetadata.KeyId)
	},
}

func init() {
	rootCmd.AddCommand(kmsCmd)

	kmsCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
	kmsCmd.Flags().StringVarP(&name, "name", "n", "", "KMS name")
}
