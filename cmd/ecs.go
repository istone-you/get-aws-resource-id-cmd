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

func getSession(profile string) *session.Session {
	if profile == "" {
		profile = "default"
	}

	return session.Must(session.NewSessionWithOptions(session.Options{
		Profile: profile,
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
		},
	}))
}

func getServiceArns(ecsClient *ecs.ECS, name string) ([]*string, error) {
	listServicesInput := &ecs.ListServicesInput{
		Cluster: &name,
	}

	listServicesOutput, err := ecsClient.ListServices(listServicesInput)
	if err != nil {
		return nil, err
	}

	return listServicesOutput.ServiceArns, nil
}

func printClusterArn(ecsClient *ecs.ECS, name string) {
	describeClustersInput := &ecs.DescribeClustersInput{
		Clusters: []*string{&name},
	}

	describeClustersOutput, err := ecsClient.DescribeClusters(describeClustersInput)
	if err != nil {
		fmt.Println("データの取得に失敗しました:", err)
		os.Exit(1)
	}

	if len(describeClustersOutput.Clusters) > 0 {
		clusterArn := describeClustersOutput.Clusters[0].ClusterArn
		fmt.Println("-----")
		fmt.Println("ClusterARN:", *clusterArn)
	} else {
		fmt.Println("クラスターが見つかりません。")
	}
}

func processService(ecsClient *ecs.ECS, name string, serviceArn *string) {
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
		return
	}

	// インスタンス情報を取得
	describeTasksInput := &ecs.DescribeTasksInput{
		Cluster: &name,
		Tasks:   listTasksOutput.TaskArns,
	}

	describeTasksOutput, err := ecsClient.DescribeTasks(describeTasksInput)
	if err != nil {
		fmt.Println("データの取得に失敗しました:", err)
		return
	}

	var isFargate bool
	var taskDifinitionName string

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
					fmt.Printf("LaunchType: EC2 (%s) \n", *instance.Ec2InstanceId)
				}
			}
		} else {
			isFargate = true
		}
		taskDifinitionName = (*task.TaskDefinitionArn)[strings.LastIndex(*task.TaskDefinitionArn, "/")+1:]
	}

	if isFargate {
		fmt.Println("LaunchType: FARGATE")
	}

	fmt.Printf("TaskDefinition: %s\n", taskDifinitionName)

	if showArn {
		fmt.Println("ServiceARN:", *serviceArn)
	}

	fmt.Println("-----")
}

func runEcsCommand(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	if name == "" {
		fmt.Print("クラスター名を入力してください: ")
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)
	}

	sess := getSession(profile)
	ecsClient := ecs.New(sess)

	if showArn {
		printClusterArn(ecsClient, name)
	}

	serviceArns, err := getServiceArns(ecsClient, name)
	if err != nil {
		fmt.Println("データの取得に失敗しました:", err)
		os.Exit(1)
	}

	fmt.Println("-----")
	for _, serviceArn := range serviceArns {
		processService(ecsClient, name, serviceArn)
	}
}

var ecsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "Show ECS Container Instance ID",
	Run:   runEcsCommand,
}

func init() {
	rootCmd.AddCommand(ecsCmd)

	ecsCmd.Flags().StringVarP(&profile, "profile", "p", "", "Set AWS CLI's profile name")
	ecsCmd.Flags().StringVarP(&name, "name", "n", "", "Set ECS cluster name")
	ecsCmd.Flags().BoolVarP(&showArn, "arn", "a", false, "Show ECS Cluster & Service Arn")
}