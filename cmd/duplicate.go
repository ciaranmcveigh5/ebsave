/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"ebsave/pkg/accountDetails"
	"ebsave/pkg/ebsaveJson"
	"ebsave/pkg/tables"
	"ebsave/pkg/volumes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

// duplicateCmd represents the duplicate command
var duplicateCmd = &cobra.Command{
	Use:   "duplicate",
	Short: "Get a list of volumes with multiple snapshots",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String("eu-west-1"),
		})
		ec2Svc := ec2.New(sess)
		stsSvc := sts.New(sess)

		e := volumes.EC2Client{
			Client: ec2Svc,
		}

		s := accountDetails.STSClient{
			Client: stsSvc,
		}
		stsInput := &sts.GetCallerIdentityInput{}
		stsResult, _ := s.GetCallerIdentity(stsInput)
		accountId := *stsResult.Account

		snapshotInput := &ec2.DescribeSnapshotsInput{
			OwnerIds: []*string{
				aws.String(accountId),
			},
		}

		result, _ := e.GetSnapshots(snapshotInput)

		volumeSnapshots := volumes.GenerateVolumeSnapshotDetails(result.Snapshots)

		returnJson, _ := cmd.Flags().GetBool("json")
		limit, _ := cmd.Flags().GetInt64("limit")

		if returnJson {
			fmt.Println(ebsaveJson.ParseDuplicateSnapshotsToJson(volumeSnapshots, limit)) // TODO: should limit be passed into Parse or a separate function
		} else {
			tableDetails := tables.ParseDuplicateSnapshotsForTable(volumeSnapshots, limit)
			tables.RenderAssetsTable(tableDetails)
		}
	},
}

func init() {
	rootCmd.AddCommand(duplicateCmd)
	duplicateCmd.Flags().BoolP("json", "j", false, "returns data in json format")
	duplicateCmd.Flags().Int64P("limit", "l", 2, "Only show volume with number of snapshots above limit")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// duplicateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// duplicateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
