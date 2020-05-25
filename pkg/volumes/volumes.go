package volumes

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type EC2Client struct {
	Client            ec2iface.EC2API
	describeVolumes   func(*ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error)
	describeInstances func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
	describeSnapshots func(*ec2.DescribeSnapshotsInput) (*ec2.DescribeSnapshotsOutput, error)
	describeImages    func(*ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error)
}

type InstanceData struct {
	Reservations []struct {
		Instances []struct {
			InstanceId string `json:"InstanceId"`
		} `json:"Instances"`
	} `json:"Reservations"`
}

type SnapshotData struct {
	Snapshots []struct {
		Description string `json:"Description"`
		SnapshotId  string `json:"SnapshotId"`
		VolumeSize  int64  `json:"VolumeSize"`
	} `json:"Snapshots"`
}

type AmiData struct {
	Images []struct {
		ImageId string `json:"ImageId"`
	} `json:"Images"`
}

type VolumeSnapshotData struct {
	VolumeSize        int64
	NumberOfSnapshots int64
	SnapshotIds       []string
}

func (e *EC2Client) DescribeVolumes(input *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	if e.describeVolumes != nil {
		result, err := e.describeVolumes(input)
		return result, err
	}

	result, err := e.Client.DescribeVolumes(input)
	return result, err
}

func (e *EC2Client) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	if e.describeInstances != nil {
		result, err := e.describeInstances(input)
		return result, err
	}

	result, err := e.Client.DescribeInstances(input)
	return result, err
}

func (e *EC2Client) DescribeSnapshots(input *ec2.DescribeSnapshotsInput) (*ec2.DescribeSnapshotsOutput, error) {
	if e.describeSnapshots != nil {
		result, err := e.describeSnapshots(input)
		return result, err
	}

	result, err := e.Client.DescribeSnapshots(input)
	return result, err
}

func (e *EC2Client) DescribeImages(input *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	if e.describeImages != nil {
		result, err := e.describeImages(input)
		return result, err
	}

	result, err := e.Client.DescribeImages(input)
	return result, err
}

func (e *EC2Client) GetVolumes(input *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	result, err := e.DescribeVolumes(input)

	if err != nil {
		awsErrorLog(err)
		return nil, err
	}

	return result, nil
}

func (e *EC2Client) GetInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	result, err := e.DescribeInstances(input)

	if err != nil {
		awsErrorLog(err)
		return nil, err
	}

	return result, nil
}

func (e *EC2Client) GetSnapshots(input *ec2.DescribeSnapshotsInput) (*ec2.DescribeSnapshotsOutput, error) {
	result, err := e.DescribeSnapshots(input)

	if err != nil {
		awsErrorLog(err)
		return nil, err
	}

	return result, nil
}

func (e *EC2Client) GetImages(input *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	result, err := e.DescribeImages(input)

	if err != nil {
		awsErrorLog(err)
		return nil, err
	}

	return result, nil
}

func GetUnattachedVolumes(volumes []*ec2.Volume) []*ec2.Volume {
	var unattachedVolumes = []*ec2.Volume{}

	for _, volume := range volumes {
		if len(volume.Attachments) == 0 {
			unattachedVolumes = append(unattachedVolumes, volume)
		}
	}

	return unattachedVolumes
}

func GetInstanceIds(instances *ec2.DescribeInstancesOutput) ([]*string, error) {

	var instanceIds = []*string{}

	d := InstanceData{}
	jsonString, err := json.Marshal(instances)
	if err != nil {
		return instanceIds, err
	}
	err = json.Unmarshal(jsonString, &d)
	if err != nil {
		return instanceIds, err
	}

	for _, reservation := range d.Reservations {
		for _, instance := range reservation.Instances {
			instanceIds = append(instanceIds, aws.String(instance.InstanceId))
		}
	}

	return instanceIds, nil
}

func GenerateVolumeSnapshotDetails(snapshots []*ec2.Snapshot) map[string]VolumeSnapshotData {
	volumeSnapshotInfo := make(map[string]VolumeSnapshotData)

	if len(snapshots) == 0 {
		return volumeSnapshotInfo
	}

	for _, snapshot := range snapshots {
		volumeSnapshotInfo[*snapshot.VolumeId] = VolumeSnapshotData{
			VolumeSize:        *snapshot.VolumeSize,
			NumberOfSnapshots: volumeSnapshotInfo[*snapshot.VolumeId].NumberOfSnapshots + 1,
			SnapshotIds:       append(volumeSnapshotInfo[*snapshot.VolumeId].SnapshotIds, *snapshot.SnapshotId),
		}
	}

	return volumeSnapshotInfo
}

func awsErrorLog(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		default:
			fmt.Println(aerr.Error())
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
	}
}
