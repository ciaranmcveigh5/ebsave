package tables

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ciaranmcveigh5/ebsave/pkg/ebsavePricing"
	"github.com/ciaranmcveigh5/ebsave/pkg/volumes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/olekukonko/tablewriter"
)

type TableDetails struct {
	Assets    []AssetDetails
	TotalCost string
	Header    []string
	Footer    []string
}

type AssetDetails struct {
	Id           string
	SizeInGB     string
	CostPerMonth string
	TableInput   []string
}

func RenderAssetsTable(t TableDetails) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(t.Header) // ToDo add tag name need error check on if tag doesn't exist

	for _, asset := range t.Assets {
		table.Append(asset.TableInput)
	}

	table.SetFooter(t.Footer)
	table.Render()
}

func ParseVolumesForTable(volumes []*ec2.Volume) TableDetails {
	var volumeDetails = []AssetDetails{}
	var tableDetails = TableDetails{
		Assets:    volumeDetails,
		TotalCost: "",
	}

	if len(volumes) == 0 {
		return tableDetails
	}

	pricingSvc := pricing.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))
	p := ebsavePricing.PricingClient{
		Client: pricingSvc,
	}

	var totalCost float64

	for _, volume := range volumes {
		v := AssetDetails{}
		volumeCost := p.GetVolumeCost(volume)
		totalCost = totalCost + volumeCost.Value
		v.Id = *volume.VolumeId
		v.CostPerMonth = ("$" + fmt.Sprintf("%.2f", volumeCost.Value))
		v.SizeInGB = strconv.FormatInt(int64(*volume.Size), 10)
		v.TableInput = []string{v.Id, v.SizeInGB, v.CostPerMonth}
		volumeDetails = append(volumeDetails, v)
	}

	tableDetails.Assets = volumeDetails
	tableDetails.TotalCost = ("$" + fmt.Sprintf("%.2f", totalCost))
	tableDetails.Header = []string{"ID", "Size(GB)", "Cost/mo"}
	tableDetails.Footer = []string{"", "Total", tableDetails.TotalCost}

	return tableDetails
}

func ParseSnapshotsForTable(snapshots []*ec2.Snapshot, amis []string) TableDetails {
	var snapshotDetails = []AssetDetails{}
	var tableDetails = TableDetails{
		Assets:    snapshotDetails,
		TotalCost: "",
	}

	if len(snapshots) == 0 {
		return tableDetails
	}

	var totalCost float64

	for _, snapshot := range snapshots {
		words := strings.Fields(*snapshot.Description)
		if len(words) == 7 {
			if words[0] == "Created" && words[1] == "by" && words[2][0:11] == "CreateImage" {
				amiId := words[4]
				amiExists := stringInSlice(amiId, amis)
				if amiExists == false {
					cost := float64(*snapshot.VolumeSize) * 0.05
					totalCost = totalCost + cost
					s := AssetDetails{}
					s.Id = *snapshot.SnapshotId
					s.SizeInGB = strconv.FormatInt(*snapshot.VolumeSize, 10)
					s.CostPerMonth = ("$" + fmt.Sprintf("%.2f", cost))
					s.TableInput = []string{s.Id, s.SizeInGB, s.CostPerMonth}
					snapshotDetails = append(snapshotDetails, s)

				}
			}
		}
	}

	tableDetails.Assets = snapshotDetails
	tableDetails.TotalCost = ("$" + fmt.Sprintf("%.2f", totalCost))
	tableDetails.Header = []string{"ID", "Size(GB)", "Cost/mo"}
	tableDetails.Footer = []string{"", "Total", tableDetails.TotalCost}

	return tableDetails
}

func ParseDuplicateSnapshotsForTable(details map[string]volumes.VolumeSnapshotData, limit int64) TableDetails {
	var volumeSnapshotDetails = []AssetDetails{}
	var tableDetails = TableDetails{
		Assets:    volumeSnapshotDetails,
		TotalCost: "",
	}

	if len(details) == 0 {
		return tableDetails
	}

	var totalCost float64

	for volumeId, volumeData := range details {
		if volumeData.NumberOfSnapshots > limit {
			cost := float64(volumeData.VolumeSize) * float64(volumeData.NumberOfSnapshots) * 0.05
			totalCost = totalCost + cost
			v := AssetDetails{}
			v.Id = volumeId
			v.SizeInGB = strconv.FormatInt(volumeData.VolumeSize, 10)
			v.CostPerMonth = ("$" + fmt.Sprintf("%.2f", cost))
			v.TableInput = []string{v.Id, v.SizeInGB, strconv.FormatInt(volumeData.NumberOfSnapshots, 10), v.CostPerMonth}
			volumeSnapshotDetails = append(volumeSnapshotDetails, v)
		}
	}

	tableDetails.Assets = volumeSnapshotDetails
	tableDetails.TotalCost = ("$" + fmt.Sprintf("%.2f", totalCost))
	tableDetails.Header = []string{"ID", "Size(GB)", "No of Snapshots", "Cost/mo"}
	tableDetails.Footer = []string{"", "", "Total", tableDetails.TotalCost}

	return tableDetails
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
