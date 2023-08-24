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
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// ecsCmd represents the ecs command
var ecsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if profile == "" {
			profile = "default"
		}

		if name == "" {
			fmt.Print("クラスター名を入力してください: ")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)
		}

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config: aws.Config{
				Region: aws.String("ap-northeast-1"),
			},
		}))

		ecsClient := ecs.New(sess)

		// ListServices APIを使用してクラスター内の全てのサービスを取得
		listServicesInput := &ecs.ListServicesInput{
			Cluster: &name,
		}

		listServicesOutput, err := ecsClient.ListServices(listServicesInput)
		if err != nil {
			fmt.Println("Error listing services:", err)
			os.Exit(1)
		}

		fmt.Println("-----")

		// サービスARNからサービス名を抽出して表示
		for _, serviceArn := range listServicesOutput.ServiceArns {
			// サービス名はARNの最後の部分（/で分割された最後の要素）
			serviceName := (*serviceArn)[strings.LastIndex(*serviceArn, "/")+1:]

			fmt.Println("Service:", serviceName)

			// サービスのインスタンス情報を取得
			listTasksInput := &ecs.ListTasksInput{
				Cluster:     &name,
				ServiceName: &serviceName,
			}

			listTasksOutput, err := ecsClient.ListTasks(listTasksInput)
			if err != nil {
				fmt.Println("データの取得に失敗しました:", err)
				continue
			}

			// インスタンス情報を取得
			describeTasksInput := &ecs.DescribeTasksInput{
				Cluster: &name,
				Tasks:   listTasksOutput.TaskArns,
			}

			describeTasksOutput, err := ecsClient.DescribeTasks(describeTasksInput)
			if err != nil {
				fmt.Println("データの取得に失敗しました:", err)
				continue
			}

			var isFargate bool

			// 各タスクのインスタンスIDを表示
			for _, task := range describeTasksOutput.Tasks {
				if *task.LaunchType == "EC2" {
					if task.ContainerInstanceArn != nil {
						describeContainerInstancesInput := &ecs.DescribeContainerInstancesInput{
							Cluster:            aws.String(name),
							ContainerInstances: []*string{task.ContainerInstanceArn},
						}
						describeContainerInstancesOutput, err := ecsClient.DescribeContainerInstances(describeContainerInstancesInput)
						if err != nil {
							fmt.Println("データの取得に失敗しました:", err)
							continue
						}

						for _, instance := range describeContainerInstancesOutput.ContainerInstances {
							fmt.Printf("LaunchType: EC2\nInstanceID: %s\n",*instance.Ec2InstanceId)
						}
					}
				} else {
					isFargate = true
				}
			}

			if isFargate {
				fmt.Println("LaunchType: FARGATE")
			}

			fmt.Println("-----")
		}
	},
}

func init() {
	rootCmd.AddCommand(ecsCmd)

	ecsCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI's profile name")
	ecsCmd.Flags().StringVarP(&name, "name", "n", "", "EC2 instance name")
}
