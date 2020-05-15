package ebsaveJson

import (
	"ebsave/pkg/ebsavePricing"
	"ebsave/pkg/volumes"
	"encoding/json"
	"math"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/pricing"
)

type AssetJson struct {
	Assets    []AssetDetailsJson
	TotalCost ebsavePricing.AssetCost
}

type AssetDetailsJson struct {
	AssetId string
	Size    struct {
		Unit  string
		Value int
	}
	Cost      ebsavePricing.AssetCost
	Snapshots []string `json:"snapshots,omitempty"`
}

func ParseVolumesToJson(volumes []*ec2.Volume) string {

	if len(volumes) == 0 {
		return "{}"
	}

	var volumeDetails = []AssetDetailsJson{}
	pricingSvc := pricing.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))
	var totalCost = ebsavePricing.AssetCost{}
	totalCost.Currency = "USD"
	totalCost.Timeframe = "monthly"

	p := ebsavePricing.PricingClient{
		Client: pricingSvc,
	}

	for _, volume := range volumes {
		v := AssetDetailsJson{}

		v.AssetId = *volume.VolumeId

		v.Size.Unit = "GB"
		v.Size.Value = int(*volume.Size)

		volumeCost := p.GetVolumeCost(volume)
		v.Cost = volumeCost

		totalCost.Value = totalCost.Value + v.Cost.Value

		volumeDetails = append(volumeDetails, v)
	}

	totalCost.Value = math.Ceil(totalCost.Value*100) / 100

	var volumesJson = AssetJson{
		Assets:    volumeDetails,
		TotalCost: totalCost,
	}

	parsedJson, _ := json.Marshal(volumesJson)
	return string(parsedJson)
}

func ParseSnapshotsToJson(snapshots []*ec2.Snapshot, amis []string) string {

	if len(snapshots) == 0 {
		return "{}"
	}

	var snapshotDetails = []AssetDetailsJson{}
	var totalCost = ebsavePricing.AssetCost{}
	totalCost.Currency = "USD"
	totalCost.Timeframe = "monthly"

	for _, snapshot := range snapshots {
		words := strings.Fields(*snapshot.Description)
		if len(words) == 7 {
			if words[0] == "Created" && words[1] == "by" && words[2][0:11] == "CreateImage" {
				amiId := words[4]
				amiExists := stringInSlice(amiId, amis)
				if amiExists == false {
					cost := float64(*snapshot.VolumeSize) * 0.05
					totalCost.Value = totalCost.Value + cost
					s := AssetDetailsJson{}

					s.AssetId = *snapshot.SnapshotId

					s.Size.Unit = "GB"
					s.Size.Value = int(*snapshot.VolumeSize)

					s.Cost.Timeframe = "monthly"
					s.Cost.Currency = "USD"
					s.Cost.Value = math.Ceil(cost*100) / 100
					snapshotDetails = append(snapshotDetails, s)

				}
			}
		}
	}

	totalCost.Value = math.Ceil(totalCost.Value*100) / 100

	var snapshotsJson = AssetJson{
		Assets:    snapshotDetails,
		TotalCost: totalCost,
	}

	if len(snapshotsJson.Assets) == 0 {
		return "{}"
	}

	parsedJson, _ := json.Marshal(snapshotsJson)
	return string(parsedJson)
}

func ParseDuplicateSnapshotsToJson(details map[string]volumes.VolumeSnapshotData, limit int64) string {

	if len(details) == 0 {
		return "{}"
	}

	var volumeSnapshotDetails = []AssetDetailsJson{}
	var totalCost = ebsavePricing.AssetCost{}
	totalCost.Currency = "USD"
	totalCost.Timeframe = "monthly"

	for volumeId, volumeData := range details {
		if volumeData.NumberOfSnapshots > limit {
			cost := float64(volumeData.VolumeSize) * float64(volumeData.NumberOfSnapshots) * 0.05
			totalCost.Value = totalCost.Value + cost

			v := AssetDetailsJson{}

			v.AssetId = volumeId

			v.Size.Unit = "GB"
			v.Size.Value = int(volumeData.VolumeSize)

			v.Cost.Timeframe = "monthly"
			v.Cost.Currency = "USD"
			v.Cost.Value = math.Ceil(cost*100) / 100

			v.Snapshots = volumeData.SnapshotIds

			volumeSnapshotDetails = append(volumeSnapshotDetails, v)
		}
	}

	if len(volumeSnapshotDetails) == 0 {
		return "{}"
	}

	totalCost.Value = math.Ceil(totalCost.Value*100) / 100

	var volumeSnapshotsJson = AssetJson{
		Assets:    volumeSnapshotDetails,
		TotalCost: totalCost,
	}

	parsedJson, _ := json.Marshal(volumeSnapshotsJson)
	return string(parsedJson)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
