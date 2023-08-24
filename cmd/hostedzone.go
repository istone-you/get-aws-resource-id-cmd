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
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/cobra"
)

var hostedzoneCmd = &cobra.Command{
	Use:   "hostedzone",
	Short: "Route53 Zone ID",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if name == "" {
			fmt.Print("ホストゾーン名を入力してください: ")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)
		}

		name = fmt.Sprintf("%s.", name)

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config: aws.Config{
				Region: aws.String("ap-northeast-1"),
			},
		}))

		route53Client := route53.New(sess)

		listHostedZonesInput := &route53.ListHostedZonesInput{}
		resp, err := route53Client.ListHostedZones(listHostedZonesInput)
		if err != nil {
			fmt.Println("データの取得に失敗しました:", err)
			os.Exit(1)
		}

		var hostedZoneID string
		for _, hostedZone := range resp.HostedZones {
			if aws.StringValue(hostedZone.Name) == name {
				hostedZoneID = aws.StringValue(hostedZone.Id)
				break
			}
		}

		if hostedZoneID == "" {
			fmt.Println("指定した名前のホストゾーンが見つかりませんでした。")
		} else {
			hostedZoneID = hostedZoneID[12:]
			fmt.Println(hostedZoneID)
		}
	},
}

func init() {
	rootCmd.AddCommand(hostedzoneCmd)

	hostedzoneCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
	hostedzoneCmd.Flags().StringVarP(&name, "name", "n", "", "Route53 zone name")
}
