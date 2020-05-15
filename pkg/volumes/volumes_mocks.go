package volumes

import "github.com/aws/aws-sdk-go/service/ec2"

var volumeId = "testVolumeId"
var tagKey = "testTagKey"
var tagValue = "testTagValue"

var tag = ec2.Tag{
	Key:   &tagKey,
	Value: &tagValue,
}

var instanceId = "testInstanceId"
var instanceState = "stopped"

var attachment = ec2.VolumeAttachment{
	InstanceId: &instanceId,
}

var singleAttachedVolume = ec2.Volume{
	VolumeId:    &volumeId,
	Tags:        []*ec2.Tag{&tag},
	Attachments: []*ec2.VolumeAttachment{&attachment},
}

var singleUnattachedVolume = ec2.Volume{
	VolumeId: &volumeId,
	Tags:     []*ec2.Tag{&tag},
}

var mockSingleAttachedVolume = []*ec2.Volume{
	&singleAttachedVolume,
}

var mockSingleUnattachedVolume = []*ec2.Volume{
	&singleUnattachedVolume,
}

var mockMultipleVolumesAttachedAndUnattached = []*ec2.Volume{
	&singleAttachedVolume,
	&singleUnattachedVolume,
}

var mockDescribeVolumes = &ec2.DescribeVolumesOutput{
	Volumes: mockSingleAttachedVolume,
}

var mockDescribeInstances = &ec2.DescribeInstancesOutput{
	Reservations: []*ec2.Reservation{
		{
			Instances: []*ec2.Instance{
				{
					InstanceId: &instanceId,
					State: &ec2.InstanceState{
						Name: &instanceState,
					},
				},
			},
		},
	},
}

var mockNoInstances = &ec2.DescribeInstancesOutput{
	Reservations: []*ec2.Reservation{
		{
			Instances: []*ec2.Instance{},
		},
	},
}

var mockDescribeInstancesMultiple = &ec2.DescribeInstancesOutput{
	Reservations: []*ec2.Reservation{
		{
			Instances: []*ec2.Instance{
				{
					InstanceId: &instanceId,
					State: &ec2.InstanceState{
						Name: &instanceState,
					},
				},
				{
					InstanceId: &instanceId,
					State: &ec2.InstanceState{
						Name: &instanceState,
					},
				},
			},
		},
	},
}

var mockNoSnapshots = &ec2.DescribeSnapshotsOutput{
	Snapshots: []*ec2.Snapshot{},
}

var mockNoImages = &ec2.DescribeImagesOutput{
	Images: []*ec2.Image{},
}
