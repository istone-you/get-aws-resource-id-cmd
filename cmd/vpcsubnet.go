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

// vpcsubnetCmd represents the vpcsubnet command
var vpcsubnetCmd = &cobra.Command{
	Use:   "vpcsubnet",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		input := &ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{aws.String(id)},
				},
			},
		}

		result, err := ec2Client.DescribeSubnets(input)
		if err != nil {
			fmt.Println("Error describing subnets:", err)
			return
		}

		fmt.Println("-----")

		for _, subnet := range result.Subnets {
			fmt.Printf("NameTag: %s\n", getTagName(subnet.Tags, "Name"))
			fmt.Printf("SubnetID: %s\n", *subnet.SubnetId)
			fmt.Println("-----")
		}
	},
}

func getTagName(tags []*ec2.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return "-"
}

func init() {
	rootCmd.AddCommand(vpcsubnetCmd)

	vpcsubnetCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
	vpcsubnetCmd.Flags().StringVarP(&id, "id", "i", "", "VPC ID")
}
