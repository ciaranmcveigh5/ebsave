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
	"os"

	"github.com/ciaranmcveigh5/ebsave/pkg/ebsaveJson"
	"github.com/ciaranmcveigh5/ebsave/pkg/tables"
	"github.com/ciaranmcveigh5/ebsave/pkg/volumes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

// unattachedCmd represents the unattached command
var unattachedCmd = &cobra.Command{
	Use:   "unattached",
	Short: "Get a list of unattached volumes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		flagAwsProfile, err := cmd.Flags().GetString("profile")
		handleError(err)
		flagAwsRegion, err := cmd.Flags().GetString("region")
		handleError(err)

		awsProfile, awsRegion := GetAwsProfileAndRegion(flagAwsProfile, flagAwsRegion)

		sess, err := session.NewSessionWithOptions(session.Options{
			Profile: awsProfile,
			Config: aws.Config{
				Region: aws.String(awsRegion),
			},
		})

		handleError(err)

		svc := ec2.New(sess)

		var input *ec2.DescribeVolumesInput
		e := volumes.EC2Client{
			Client: svc,
		}

		result, err := e.GetVolumes(input)
		handleError(err)

		unattachedVolumes := volumes.GetUnattachedVolumes(result.Volumes)

		returnJson, err := cmd.Flags().GetBool("json")
		handleError(err)

		if returnJson {
			json, err := ebsaveJson.ParseVolumesToJson(unattachedVolumes)
			handleError(err)
			fmt.Println(json)
		} else {
			tableDetails, err := tables.ParseVolumesForTable(unattachedVolumes)
			handleError(err)
			tables.RenderAssetsTable(tableDetails)
		}
	},
}

func GetAwsProfileAndRegion(flagAwsProfile, flagAwsRegion string) (customAwsProfile, customAwsRegion string) {

	if flagAwsProfile != "" {
		customAwsProfile = flagAwsProfile
	} else if os.Getenv("AWS_PROFILE") != "" {
		customAwsProfile = os.Getenv("AWS_PROFILE")
	} else {
		customAwsProfile = "default"
		// fmt.Println("AWS Profile must be set via the command line -profile or via the env var AWS_PROFILE")
		// os.Exit(1)
	}

	if flagAwsRegion != "" {
		customAwsRegion = flagAwsRegion
	} else if os.Getenv("AWS_REGION") != "" {
		customAwsRegion = os.Getenv("AWS_REGION")
	} else {
		customAwsRegion = "eu-west-1"
		// fmt.Println("AWS Region must be set via the command line -region or via the env var AWS_REGION")
		// os.Exit(1)
	}

	return customAwsProfile, customAwsRegion
}

func init() {
	rootCmd.AddCommand(unattachedCmd)
	unattachedCmd.Flags().BoolP("json", "j", false, "returns data in json format")
	unattachedCmd.Flags().StringP("profile", "p", "", "set aws profile to use, will default to default profile if not set, can also be set via the env var AWS_PROFILE")
	unattachedCmd.Flags().StringP("region", "r", "", "set aws region to use, will default to eu-west-1 if not set, can also be set via the env var AWS_REGION")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unattachedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unattachedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
