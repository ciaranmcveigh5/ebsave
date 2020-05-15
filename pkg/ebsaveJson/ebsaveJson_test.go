package ebsaveJson

import (
	"ebsave/pkg/volumes"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestParseVolumesToJson(t *testing.T) {

	t.Run("no volumes", func(t *testing.T) {
		volumes := []*ec2.Volume{}
		returned := ParseVolumesToJson(volumes)
		expected := "{}"
		assertEqual(t, expected, returned)

	})

	t.Run("single volume", func(t *testing.T) {
		volumes := []*ec2.Volume{
			{
				VolumeId:         aws.String("vol-123"),
				Size:             aws.Int64(600),
				AvailabilityZone: aws.String("eu-west-1a"),
				VolumeType:       aws.String("gp2"),
			},
		}
		returned := ParseVolumesToJson(volumes)
		expected := `{"Assets":[{"AssetId":"vol-123","Size":{"Unit":"GB","Value":600},"Cost":{"Timeframe":"monthly","Currency":"USD","Value":66}}],"TotalCost":{"Timeframe":"monthly","Currency":"USD","Value":66}}`
		assertEqual(t, expected, returned)

	})

	t.Run("multiple volumes", func(t *testing.T) {
		volumes := []*ec2.Volume{
			{
				VolumeId:         aws.String("vol-123"),
				Size:             aws.Int64(600),
				AvailabilityZone: aws.String("eu-west-1a"),
				VolumeType:       aws.String("gp2"),
			},
			{
				VolumeId:         aws.String("vol-abc"),
				Size:             aws.Int64(200),
				AvailabilityZone: aws.String("eu-west-1a"),
				VolumeType:       aws.String("gp2"),
			},
		}
		returned := ParseVolumesToJson(volumes)
		expected := `{"Assets":[{"AssetId":"vol-123","Size":{"Unit":"GB","Value":600},"Cost":{"Timeframe":"monthly","Currency":"USD","Value":66}},{"AssetId":"vol-abc","Size":{"Unit":"GB","Value":200},"Cost":{"Timeframe":"monthly","Currency":"USD","Value":22}}],"TotalCost":{"Timeframe":"monthly","Currency":"USD","Value":88}}`
		assertEqual(t, expected, returned)
	})
}

func TestParseSnapshotsToJson(t *testing.T) {
	t.Run("no snapshots, no amis", func(t *testing.T) {
		snapshots := []*ec2.Snapshot{}
		amis := []string{}
		returned := ParseSnapshotsToJson(snapshots, amis)
		expected := "{}"
		assertEqual(t, expected, returned)
	})

	t.Run("single snapshot, no amis", func(t *testing.T) {
		snapshots := []*ec2.Snapshot{
			{
				SnapshotId:  aws.String("snap-123"),
				VolumeSize:  aws.Int64(200),
				Description: aws.String("Created by CreateImage(i-04644c09efbefb22b) for ami-041bd1d580f0446e1 from vol-0c885534270c99703"),
			},
		}
		amis := []string{}

		returned := ParseSnapshotsToJson(snapshots, amis)
		expected := `{"Assets":[{"AssetId":"snap-123","Size":{"Unit":"GB","Value":200},"Cost":{"Timeframe":"monthly","Currency":"USD","Value":10}}],"TotalCost":{"Timeframe":"monthly","Currency":"USD","Value":10}}`
		assertEqual(t, expected, returned)
	})

	t.Run("single snapshot, ami exists", func(t *testing.T) {
		snapshots := []*ec2.Snapshot{
			{
				SnapshotId:  aws.String("snap-123"),
				VolumeSize:  aws.Int64(200),
				Description: aws.String("Created by CreateImage(i-04644c09efbefb22b) for ami-041bd1d580f0446e1 from vol-0c885534270c99703"),
			},
		}
		amis := []string{
			"ami-041bd1d580f0446e1",
		}

		returned := ParseSnapshotsToJson(snapshots, amis)
		expected := "{}"
		assertEqual(t, expected, returned)

	})

	t.Run("multiple snapshot, both ami exists and doesn't exist", func(t *testing.T) {
		snapshots := []*ec2.Snapshot{
			{
				SnapshotId:  aws.String("snap-123"),
				VolumeSize:  aws.Int64(200),
				Description: aws.String("Created by CreateImage(i-04644c09efbefb22b) for ami-041bd1d580f0446e1 from vol-0c885534270c99703"),
			},
			{
				SnapshotId:  aws.String("snap-abc"),
				VolumeSize:  aws.Int64(200),
				Description: aws.String("Created by CreateImage(i-04644c09efbefb22b) for ami-abc from vol-0c885534270c99703"),
			},
			{
				SnapshotId:  aws.String("snap-xyz"),
				VolumeSize:  aws.Int64(400),
				Description: aws.String("Created by CreateImage(i-04644c09efbefb22b) for ami-xyz from vol-0c885534270c99703"),
			},
		}
		amis := []string{
			"ami-041bd1d580f0446e1",
		}

		returned := ParseSnapshotsToJson(snapshots, amis)
		expected := `{"Assets":[{"AssetId":"snap-abc","Size":{"Unit":"GB","Value":200},"Cost":{"Timeframe":"monthly","Currency":"USD","Value":10}},{"AssetId":"snap-xyz","Size":{"Unit":"GB","Value":400},"Cost":{"Timeframe":"monthly","Currency":"USD","Value":20}}],"TotalCost":{"Timeframe":"monthly","Currency":"USD","Value":30}}`
		assertEqual(t, expected, returned)

	})

}

func TestParseDuplicateSnapshotsToJson(t *testing.T) {
	t.Run("no snapshots/volumes", func(t *testing.T) {
		volumes := make(map[string]volumes.VolumeSnapshotData)
		returned := ParseDuplicateSnapshotsToJson(volumes, 2)
		expected := "{}"

		assertEqual(t, expected, returned)
	})

	t.Run("duplicate volumes", func(t *testing.T) {
		volumes := map[string]volumes.VolumeSnapshotData{
			"vol-123": {
				VolumeSize:        200,
				NumberOfSnapshots: 3,
				SnapshotIds:       []string{"snap-123", "snap-456", "snap-789"},
			},
		}
		returned := ParseDuplicateSnapshotsToJson(volumes, 2)
		expected := `{"Assets":[{"AssetId":"vol-123","Size":{"Unit":"GB","Value":200},"Cost":{"Timeframe":"monthly","Currency":"USD","Value":30},"snapshots":["snap-123","snap-456","snap-789"]}],"TotalCost":{"Timeframe":"monthly","Currency":"USD","Value":30}}`

		assertEqual(t, expected, returned)
	})

	t.Run("not enough duplicate snapshots", func(t *testing.T) {
		volumes := map[string]volumes.VolumeSnapshotData{
			"vol-123": {
				VolumeSize:        200,
				NumberOfSnapshots: 3,
			},
		}
		returned := ParseDuplicateSnapshotsToJson(volumes, 4)
		expected := "{}"

		assertEqual(t, expected, returned)
	})

	t.Run("multiple duplicate volumes", func(t *testing.T) {
		volumes := map[string]volumes.VolumeSnapshotData{
			"vol-123": {
				VolumeSize:        200,
				NumberOfSnapshots: 3,
				SnapshotIds:       []string{"snap-123", "snap-456", "snap-789"},
			},
			"vol-abc": {
				VolumeSize:        100,
				NumberOfSnapshots: 1,
				SnapshotIds:       []string{"snap-abc"},
			},
			"vol-xyz": {
				VolumeSize:        400,
				NumberOfSnapshots: 4,
				SnapshotIds:       []string{"snap-xyz", "snap-123", "snap-456", "snap-789"},
			},
		}
		returned := ParseDuplicateSnapshotsToJson(volumes, 2)
		expected := `{"Assets":[{"AssetId":"vol-123","Size":{"Unit":"GB","Value":200},"Cost":{"Timeframe":"monthly","Currency":"USD","Value":30},"snapshots":["snap-123","snap-456","snap-789"]},{"AssetId":"vol-xyz","Size":{"Unit":"GB","Value":400},"Cost":{"Timeframe":"monthly","Currency":"USD","Value":80},"snapshots":["snap-xyz","snap-123","snap-456","snap-789"]}],"TotalCost":{"Timeframe":"monthly","Currency":"USD","Value":110}}`

		assertEqual(t, expected, returned)
	})

}

func assertEqual(t *testing.T, expected, returned string) {
	t.Helper()
	if expected != returned {
		t.Errorf("Expected %v but returned %v", expected, returned)
	}
}
