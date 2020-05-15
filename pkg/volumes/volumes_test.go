package volumes

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestGetVolumes(t *testing.T) {
	t.Run("Get all volumes", func(t *testing.T) {

		var input *ec2.DescribeVolumesInput

		ec2 := &EC2Client{
			describeVolumes: func(*ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
				return mockDescribeVolumes, nil
			},
		}

		returned, _ := ec2.GetVolumes(input)

		if returned != mockDescribeVolumes {
			t.Errorf("expect %v but returned %v", mockDescribeVolumes, returned)
		}
	})
}

func TestGetInstances(t *testing.T) {
	t.Run("Get instance state stopped", func(t *testing.T) {

		var input *ec2.DescribeInstancesInput

		ec2 := &EC2Client{
			describeInstances: func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
				return mockDescribeInstances, nil
			},
		}

		returned, _ := ec2.GetInstances(input)

		if returned != mockDescribeInstances {
			t.Errorf("expect %v but returned %v", mockDescribeInstances, returned)
		}
	})
}

func TestGetSnapshots(t *testing.T) {
	t.Run("no snapshots", func(t *testing.T) {
		var input *ec2.DescribeSnapshotsInput

		ec2 := &EC2Client{
			describeSnapshots: func(*ec2.DescribeSnapshotsInput) (*ec2.DescribeSnapshotsOutput, error) {
				return mockNoSnapshots, nil
			},
		}

		returned, _ := ec2.GetSnapshots(input)
		if returned != mockNoSnapshots {
			t.Errorf("expect %v but returned %v", mockNoSnapshots, returned)
		}
	})
}

func TestGetImages(t *testing.T) {
	t.Run("no images", func(t *testing.T) {
		var input *ec2.DescribeImagesInput

		ec2 := &EC2Client{
			describeImages: func(*ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
				return mockNoImages, nil
			},
		}

		returned, _ := ec2.GetImages(input)
		if returned != mockNoImages {
			t.Errorf("expect %v but returned %v", mockNoImages, returned)
		}
	})
}

func TestGetUnattachedVolumes(t *testing.T) {

	t.Run("empty Volumes array", func(t *testing.T) {
		returned, _ := GetUnattachedVolumes([]*ec2.Volume{})
		expected := []*ec2.Volume{}
		reflectDeepEqual(t, expected, returned)
	})

	t.Run("single attached volumes", func(t *testing.T) {
		returned, _ := GetUnattachedVolumes(mockSingleAttachedVolume)
		expected := []*ec2.Volume{}
		reflectDeepEqual(t, expected, returned)
	})

	t.Run("single unattached volume", func(t *testing.T) {
		returned, _ := GetUnattachedVolumes(mockSingleUnattachedVolume)
		reflectDeepEqual(t, mockSingleUnattachedVolume, returned)

	})

	t.Run("multiple volumes attached and unattached", func(t *testing.T) {
		returned, _ := GetUnattachedVolumes(mockMultipleVolumesAttachedAndUnattached)
		expected := []*ec2.Volume{
			{
				VolumeId: &volumeId,
				Tags:     []*ec2.Tag{&tag},
			},
		}
		reflectDeepEqual(t, expected, returned)
	})
}

func TestGenerateVolumeSnapshotDetails(t *testing.T) {
	t.Run("No snapshots", func(t *testing.T) {
		snapshots := []*ec2.Snapshot{}
		expected := make(map[string]VolumeSnapshotData)
		returned := GenerateVolumeSnapshotDetails(snapshots)
		reflectDeepEqual(t, expected, returned)
	})

	t.Run("single snapshots", func(t *testing.T) {
		snapshots := []*ec2.Snapshot{
			{
				SnapshotId: aws.String("snap-123"),
				VolumeSize: aws.Int64(200),
				VolumeId:   aws.String("vol-123"),
			},
		}
		expected := map[string]VolumeSnapshotData{
			"vol-123": {
				VolumeSize:        200,
				NumberOfSnapshots: 1,
				SnapshotIds:       []string{"snap-123"},
			},
		}
		returned := GenerateVolumeSnapshotDetails(snapshots)
		reflectDeepEqual(t, expected, returned)
	})

	t.Run("multiple snapshots same volume", func(t *testing.T) {
		snapshots := []*ec2.Snapshot{
			{
				SnapshotId: aws.String("snap-123"),
				VolumeSize: aws.Int64(200),
				VolumeId:   aws.String("vol-123"),
			},
			{
				SnapshotId: aws.String("snap-456"),
				VolumeSize: aws.Int64(200),
				VolumeId:   aws.String("vol-123"),
			},
		}
		expected := map[string]VolumeSnapshotData{
			"vol-123": {
				VolumeSize:        200,
				NumberOfSnapshots: 2,
				SnapshotIds:       []string{"snap-123", "snap-456"},
			},
		}
		returned := GenerateVolumeSnapshotDetails(snapshots)
		reflectDeepEqual(t, expected, returned)
	})

	t.Run("multiple snapshots different volumes", func(t *testing.T) {
		snapshots := []*ec2.Snapshot{
			{
				SnapshotId: aws.String("snap-123"),
				VolumeSize: aws.Int64(200),
				VolumeId:   aws.String("vol-123"),
			},
			{
				SnapshotId: aws.String("snap-456"),
				VolumeSize: aws.Int64(400),
				VolumeId:   aws.String("vol-456"),
			},
		}
		expected := map[string]VolumeSnapshotData{
			"vol-123": {
				VolumeSize:        200,
				NumberOfSnapshots: 1,
				SnapshotIds:       []string{"snap-123"},
			},
			"vol-456": {
				VolumeSize:        400,
				NumberOfSnapshots: 1,
				SnapshotIds:       []string{"snap-456"},
			},
		}
		returned := GenerateVolumeSnapshotDetails(snapshots)
		reflectDeepEqual(t, expected, returned)
	})
}

func TestGetVolumesAttachedToStoppedInstances(t *testing.T) {

	t.Run("no stopped instances therefore no volumes", func(t *testing.T) {
		e := &EC2Client{
			describeVolumes: func(*ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
				return mockDescribeVolumes, nil
			},
			describeInstances: func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
				return mockNoInstances, nil
			},
		}

		expected := []*ec2.Volume{}
		returned, _ := GetVolumesAttachedToStoppedInstances(e)
		reflectDeepEqual(t, expected, returned)
	})

	// t.Run("single stopped instance", func(t *testing.T) {
	// 	e := &EC2Client{
	// 		describeVolumes: func(*ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	// 			return mockDescribeVolumes, nil
	// 		},
	// 		describeInstances: func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	// 			return mockNoInstances, nil
	// 		},
	// 	}

	// 	expected := []*ec2.Volume{}
	// 	returned, _ := GetVolumesAttachedToStoppedInstances(e)
	// })
}

func TestGetInstanceIds(t *testing.T) {
	t.Run("No instances", func(t *testing.T) {
		expected := []*string{}
		returned := GetInstanceIds(mockNoInstances)

		reflectDeepEqual(t, expected, returned)
	})

	t.Run("single instance", func(t *testing.T) {
		expected := []*string{
			aws.String("testInstanceId"),
		}
		returned := GetInstanceIds(mockDescribeInstances)

		reflectDeepEqual(t, expected, returned)
	})

	t.Run("multiple instances", func(t *testing.T) {
		expected := []*string{
			aws.String("testInstanceId"),
			aws.String("testInstanceId"),
		}
		returned := GetInstanceIds(mockDescribeInstancesMultiple)

		reflectDeepEqual(t, expected, returned)
	})
}

func reflectDeepEqual(t *testing.T, expected, returned interface{}) {
	t.Helper()
	if !reflect.DeepEqual(returned, expected) {
		t.Errorf("expected %q but returned %q", expected, returned)
	}
}
