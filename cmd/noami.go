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
	"fmt"

	"github.com/ciaranmcveigh5/ebsave/pkg/accountDetails"
	"github.com/ciaranmcveigh5/ebsave/pkg/ebsaveJson"
	"github.com/ciaranmcveigh5/ebsave/pkg/tables"
	"github.com/ciaranmcveigh5/ebsave/pkg/volumes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

// noamiCmd represents the noami command
var noamiCmd = &cobra.Command{
	Use:   "noami",
	Short: "Get a list of ami snapshots for ami's that no longer exist",
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

		snapshotResult, _ := e.GetSnapshots(snapshotInput)

		amiInput := &ec2.DescribeImagesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("owner-id"),
					Values: []*string{
						aws.String(accountId),
					},
				},
			},
		}
		amiResult, _ := e.GetImages(amiInput)

		var amiArray []string

		for _, ami := range amiResult.Images {
			amiArray = append(amiArray, *ami.ImageId)
		}

		returnJson, _ := cmd.Flags().GetBool("json")
		if returnJson {
			fmt.Println(ebsaveJson.ParseSnapshotsToJson(snapshotResult.Snapshots, amiArray))
		} else {
			tableDetails := tables.ParseSnapshotsForTable(snapshotResult.Snapshots, amiArray)
			tables.RenderAssetsTable(tableDetails)
		}
	},
}

func init() {
	rootCmd.AddCommand(noamiCmd)
	noamiCmd.Flags().BoolP("json", "j", false, "returns data in json format")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// noamiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// noamiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
