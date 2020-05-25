package tables

import (
	"reflect"
	"testing"

	"github.com/ciaranmcveigh5/ebsave/pkg/volumes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestParseVolumesForTable(t *testing.T) {
	t.Run("no volumes", func(t *testing.T) {
		volumes := []*ec2.Volume{}
		returned, err := ParseVolumesForTable(volumes)
		expected := TableDetails{
			Assets:    []AssetDetails{},
			TotalCost: "",
		}

		reflectDeepEqual(t, expected, returned)
		assertError(t, nil, err)
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

		returned, err := ParseVolumesForTable(volumes)
		expected := TableDetails{
			Assets: []AssetDetails{
				{
					Id:           "vol-123",
					SizeInGB:     "600",
					CostPerMonth: "$66.00",
					TableInput:   []string{"vol-123", "600", "$66.00"},
				},
			},
			TotalCost: "$66.00",
			Header:    []string{"ID", "Size(GB)", "Cost/mo"},
			Footer:    []string{"", "Total", "$66.00"},
		}

		reflectDeepEqual(t, expected, returned)
		assertError(t, nil, err)
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

		returned, err := ParseVolumesForTable(volumes)
		expected := TableDetails{
			Assets: []AssetDetails{
				{
					Id:           "vol-123",
					SizeInGB:     "600",
					CostPerMonth: "$66.00",
					TableInput:   []string{"vol-123", "600", "$66.00"},
				},
				{
					Id:           "vol-abc",
					SizeInGB:     "200",
					CostPerMonth: "$22.00",
					TableInput:   []string{"vol-abc", "200", "$22.00"},
				},
			},
			TotalCost: "$88.00",
			Header:    []string{"ID", "Size(GB)", "Cost/mo"},
			Footer:    []string{"", "Total", "$88.00"},
		}

		reflectDeepEqual(t, expected, returned)
		assertError(t, nil, err)
	})
}

func TestParseSnapshotsForTable(t *testing.T) {
	t.Run("no snapshots", func(t *testing.T) {
		snapshots := []*ec2.Snapshot{}
		amis := []string{}
		returned := ParseSnapshotsForTable(snapshots, amis)
		expected := TableDetails{
			Assets:    []AssetDetails{},
			TotalCost: "",
		}

		reflectDeepEqual(t, expected, returned)
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

		returned := ParseSnapshotsForTable(snapshots, amis)
		expected := TableDetails{
			Assets: []AssetDetails{
				{
					Id:           "snap-123",
					SizeInGB:     "200",
					CostPerMonth: "$10.00",
					TableInput:   []string{"snap-123", "200", "$10.00"},
				},
			},
			TotalCost: "$10.00",
			Header:    []string{"ID", "Size(GB)", "Cost/mo"},
			Footer:    []string{"", "Total", "$10.00"},
		}

		reflectDeepEqual(t, expected, returned)

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

		returned := ParseSnapshotsForTable(snapshots, amis)
		expected := TableDetails{
			Assets:    []AssetDetails{},
			TotalCost: "$0.00",
			Header:    []string{"ID", "Size(GB)", "Cost/mo"},
			Footer:    []string{"", "Total", "$0.00"},
		}

		reflectDeepEqual(t, expected, returned)

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

		returned := ParseSnapshotsForTable(snapshots, amis)
		expected := TableDetails{
			Assets: []AssetDetails{
				{
					Id:           "snap-abc",
					SizeInGB:     "200",
					CostPerMonth: "$10.00",
					TableInput:   []string{"snap-abc", "200", "$10.00"},
				},
				{
					Id:           "snap-xyz",
					SizeInGB:     "400",
					CostPerMonth: "$20.00",
					TableInput:   []string{"snap-xyz", "400", "$20.00"},
				},
			},
			TotalCost: "$30.00",
			Header:    []string{"ID", "Size(GB)", "Cost/mo"},
			Footer:    []string{"", "Total", "$30.00"},
		}

		reflectDeepEqual(t, expected, returned)

	})

}

func TestParseDuplicateSnapshotsForTable(t *testing.T) {
	t.Run("no snapshots/volumes", func(t *testing.T) {
		volumes := make(map[string]volumes.VolumeSnapshotData)
		returned := ParseDuplicateSnapshotsForTable(volumes, 2)
		expected := TableDetails{
			Assets:    []AssetDetails{},
			TotalCost: "",
		}

		reflectDeepEqual(t, expected, returned)
	})

	t.Run("duplicate volumes", func(t *testing.T) {
		volumes := map[string]volumes.VolumeSnapshotData{
			"vol-123": {
				VolumeSize:        200,
				NumberOfSnapshots: 3,
			},
		}
		returned := ParseDuplicateSnapshotsForTable(volumes, 2)
		expected := TableDetails{
			Assets: []AssetDetails{
				{
					Id:           "vol-123",
					SizeInGB:     "200",
					CostPerMonth: "$30.00",
					TableInput:   []string{"vol-123", "200", "3", "$30.00"},
				},
			},
			TotalCost: "$30.00",
			Header:    []string{"ID", "Size(GB)", "No of Snapshots", "Cost/mo"},
			Footer:    []string{"", "", "Total", "$30.00"},
		}

		reflectDeepEqual(t, expected, returned)
	})

	t.Run("not enough duplicate snapshots", func(t *testing.T) {
		volumes := map[string]volumes.VolumeSnapshotData{
			"vol-123": {
				VolumeSize:        200,
				NumberOfSnapshots: 3,
			},
		}
		returned := ParseDuplicateSnapshotsForTable(volumes, 4)
		expected := TableDetails{
			Assets:    []AssetDetails{},
			TotalCost: "$0.00",
			Header:    []string{"ID", "Size(GB)", "No of Snapshots", "Cost/mo"},
			Footer:    []string{"", "", "Total", "$0.00"},
		}

		reflectDeepEqual(t, expected, returned)
	})

	t.Run("multiple duplicate volumes", func(t *testing.T) {
		volumes := map[string]volumes.VolumeSnapshotData{
			"vol-123": {
				VolumeSize:        200,
				NumberOfSnapshots: 3,
			},
			"vol-abc": {
				VolumeSize:        100,
				NumberOfSnapshots: 1,
			},
			"vol-xyz": {
				VolumeSize:        400,
				NumberOfSnapshots: 4,
			},
		}
		returned := ParseDuplicateSnapshotsForTable(volumes, 2)
		expected := TableDetails{
			Assets: []AssetDetails{
				{
					Id:           "vol-123",
					SizeInGB:     "200",
					CostPerMonth: "$30.00",
					TableInput:   []string{"vol-123", "200", "3", "$30.00"},
				},
				{
					Id:           "vol-xyz",
					SizeInGB:     "400",
					CostPerMonth: "$80.00",
					TableInput:   []string{"vol-xyz", "400", "4", "$80.00"},
				},
			},
			TotalCost: "$110.00",
			Header:    []string{"ID", "Size(GB)", "No of Snapshots", "Cost/mo"},
			Footer:    []string{"", "", "Total", "$110.00"},
		}

		reflectDeepEqual(t, expected, returned)
	})

}

func reflectDeepEqual(t *testing.T, expected, returned TableDetails) {
	t.Helper()
	if !reflect.DeepEqual(expected, returned) {
		t.Errorf("Expected %v but returned %v", expected, returned)
	}
}

func assertError(t *testing.T, expected, returned error) {
	t.Helper()
	if expected != returned {
		t.Errorf("expected %q but returned %q", expected, returned)
	}
	if returned == nil && expected != nil {
		t.Fatal("expected to get an error.")
	}
}
