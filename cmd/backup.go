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
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if name == "" {
			fmt.Print("Backupプラン名を入力してください: ")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)
		}

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config: aws.Config{
				Region: aws.String("ap-northeast-1"),
			},
		}))

		backupClient := backup.New(sess)

		backupPlanID, err := getBackupPlanID(backupClient, name)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			return
		}

		fmt.Println(backupPlanID)
	},
}

func getBackupPlanID(backupClient *backup.Backup, backupPlanName string) (string, error) {
	result, err := backupClient.ListBackupPlans(&backup.ListBackupPlansInput{})
	if err != nil {
		return "", err
	}

	for _, backupPlan := range result.BackupPlansList {
		if *backupPlan.BackupPlanName == backupPlanName {
			return *backupPlan.BackupPlanId, nil
		}
	}

	return "", fmt.Errorf("Backupプランが見つかりません。")
}

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
	backupCmd.Flags().StringVarP(&name, "name", "n", "", "Backup plan name")
}
