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

	"github.com/ciaranmcveigh5/ebsave/pkg/ebsaveJson"
	"github.com/ciaranmcveigh5/ebsave/pkg/tables"
	"github.com/ciaranmcveigh5/ebsave/pkg/volumes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

// stoppedCmd represents the stopped command
var stoppedCmd = &cobra.Command{
	Use:   "stopped",
	Short: "Get a list of volumes attached to stopped instances",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		svc := ec2.New(session.New(), aws.NewConfig().WithRegion("eu-west-1"))
		getStoppedInstancesinput := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("instance-state-name"),
					Values: []*string{
						aws.String("stopped"),
					},
				},
			},
		}
		e := volumes.EC2Client{
			Client: svc,
		}
		instanceResult, err := e.GetInstances(getStoppedInstancesinput)
		handleError(err)
		instanceIds, err := volumes.GetInstanceIds(instanceResult)
		handleError(err)
		volumesInput := &ec2.DescribeVolumesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("attachment.instance-id"),
					Values: instanceIds,
				},
			},
		}

		volumeResult, err := e.GetVolumes(volumesInput)
		handleError(err)
		returnJson, err := cmd.Flags().GetBool("json")
		handleError(err)
		if returnJson {
			json, err := ebsaveJson.ParseVolumesToJson(volumeResult.Volumes)
			handleError(err)
			fmt.Println(json)
		} else {
			tableDetails, err := tables.ParseVolumesForTable(volumeResult.Volumes)
			handleError(err)
			tables.RenderAssetsTable(tableDetails)
		}
	},
}

func init() {
	rootCmd.AddCommand(stoppedCmd)
	stoppedCmd.Flags().BoolP("json", "j", false, "returns data in json format")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stoppedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stoppedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
