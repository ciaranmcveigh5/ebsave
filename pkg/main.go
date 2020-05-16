package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ciaranmcveigh5/ebsave/pkg/volumes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/pricing"
)

func main() {

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})
	svc := ec2.New(sess)

	input := &ec2.DescribeVolumesInput{}

	customEC2 := volumes.EC2Client{
		Client: svc,
	}

	result, _ := customEC2.GetVolumes(input)

	input2 := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String("i-0a38dc703205cdeb2"),
		},
	}

	result2, _ := svc.DescribeInstances(input2)

	volume := result.Volumes[0]
	fmt.Println(result2)

	awsRegion := map[string]string{
		"ca-central-1":   "Canada (Central)",
		"ap-northeast-3": "Asia Pacific (Osaka-Local)",
		"us-east-1":      "US East (N. Virginia)",
		"ap-northeast-2": "Asia Pacific (Seoul)",
		"us-gov-west-1":  "AWS GovCloud (US)",
		"us-east-2":      "US East (Ohio)",
		"ap-northeast-1": "Asia Pacific (Tokyo)",
		"ap-south-1":     "Asia Pacific (Mumbai)",
		"ap-southeast-2": "Asia Pacific (Sydney)",
		"ap-southeast-1": "Asia Pacific (Singapore)",
		"sa-east-1":      "South America (Sao Paulo)",
		"us-west-2":      "US West (Oregon)",
		"eu-west-1":      "EU (Ireland)",
		"eu-west-3":      "EU (Paris)",
		"eu-west-2":      "EU (London)",
		"us-west-1":      "US West (N. California)",
		"eu-central-1":   "EU (Frankfurt)",
	}

	ebsName := map[string]string{
		"standard": "Magnetic",
		"gp2":      "General Purpose",
		"io1":      "Provisioned IOPS",
		"st1":      "Throughput Optimized HDD",
		"sc1":      "Cold HDD",
	}

	// getVolumeTypeDescription
	// getVolumeRegion

	pricingSvc := pricing.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

	// az to region function
	// stringLength := len(*volume.AvailabilityZone)
	// str := *volume.AvailabilityZone
	// last := str[stringLength-1:]
	// ebsRegion := strings.TrimSuffix(str, last)

	input3 := &pricing.GetProductsInput{
		MaxResults:  aws.Int64(1),
		ServiceCode: aws.String("AmazonEC2"),
		Filters: []*pricing.Filter{
			&pricing.Filter{
				Type:  aws.String("TERM_MATCH"),
				Field: aws.String("volumeType"),
				Value: aws.String(ebsName["io1"]),
			},
			&pricing.Filter{
				Type:  aws.String("TERM_MATCH"),
				Field: aws.String("location"),
				Value: aws.String(awsRegion["eu-west-1"]),
			},
		},
	}

	// input4 := &pricing.GetProductsInput{
	// 	MaxResults:  aws.Int64(1),
	// 	ServiceCode: aws.String("AmazonEC2"),
	// }

	result3, err := pricingSvc.GetProducts(input3)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case pricing.ErrCodeInternalErrorException:
				fmt.Println(pricing.ErrCodeInternalErrorException, aerr.Error())
			case pricing.ErrCodeInvalidParameterException:
				fmt.Println(pricing.ErrCodeInvalidParameterException, aerr.Error())
			case pricing.ErrCodeNotFoundException:
				fmt.Println(pricing.ErrCodeNotFoundException, aerr.Error())
			case pricing.ErrCodeInvalidNextTokenException:
				fmt.Println(pricing.ErrCodeInvalidNextTokenException, aerr.Error())
			case pricing.ErrCodeExpiredNextTokenException:
				fmt.Println(pricing.ErrCodeExpiredNextTokenException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}

	type Data struct {
		Terms struct {
			OnDemand struct {
				GDX844Q9TQAG2YZ2JRTCKXETXF struct {
					PriceDimensions struct {
						GDX844Q9TQAG2YZ2JRTCKXETXF6YS6EN2CT7 struct {
							PricePerUnit struct {
								USD string `json:"USD"`
							} `json:"pricePerUnit"`
						} `json:"GDX844Q9TQAG2YZ2.JRTCKXETXF.6YS6EN2CT7"`
					} `json:"priceDimensions"`
				} `json:"GDX844Q9TQAG2YZ2.JRTCKXETXF"`
			} `json:"OnDemand"`
		} `json:"terms"`
	}

	d := Data{}
	jsonString, _ := json.Marshal(result3.PriceList[0])
	json.Unmarshal(jsonString, &d)
	price, _ := strconv.ParseFloat(d.Terms.OnDemand.GDX844Q9TQAG2YZ2JRTCKXETXF.PriceDimensions.GDX844Q9TQAG2YZ2JRTCKXETXF6YS6EN2CT7.PricePerUnit.USD, 64)

	fmt.Println(price)
	fmt.Println(input3)
	fmt.Println(volume)
	fmt.Println(result3)

	// result3.PriceList

	// priceResult := ""

	// for _, m := range result3.PriceList {
	// 	fmt.Println("----------------")
	// 	getUSDFromMap(m, &priceResult)
	// 	fmt.Println(priceResult)
	// 	fmt.Println("xxxxxxxxxxxxxxxxx")
	// }

	// tUSD(result3.PriceList[0])

}

func getUSDFromMap(m map[string]interface{}, result *string) {
	for k, v := range m {
		_, ok := v.(map[string]interface{})
		if ok {
			getUSDFromMap(v.(map[string]interface{}), result)
		}

		if k == "USD" {
			*result = v.(string)
		}
	}
}

// func isCompleted() func() bool {
// 	completed := false
// 	return isCompleted
// }

// func getUSDFromMap(m map[string]interface{}) float64 {

// 	onDemand := m["terms"].(map[string]interface{})["OnDemand"].(map[string]interface{})
// 	for _, v := range onDemand {
// 		for _, v2 := range v.(map[string]interface{}) {
// 			for k3, v3 := range v2.(map[string]interface{}) {
// 				fmt.Println(k3)
// 				fmt.Println(v3)
// 			}
// 		}
// 	}
// 	return 0.00
// }

// func getUSDFromMap2(m map[string]interface{}) float64 {
// 	v := reflect.ValueOf(m)
// 	if v.Kind() == reflect.Map {
// 		for _, key := range v.MapKeys() {
// 			x := v.MapIndex(key)
// 			if x.Kind() == reflect.Map {
// 				fmt.Println("testing.....")
// 				fmt.Println(key)
// 				fmt.Println(x)
// 				getUSDFromMap(x)
// 				return 0.00
// 			}
// 		}
// 	}
// 	return 0.00
// }

func isMap(t interface{}) bool {
	switch t.(type) {
	case map[string]interface{}:
		return true
	default:
		return false
	}
}

// [map
// 	[product:
// 		map[
// 			attributes:
// 				map[
// 					location: EU (Ireland)
// 					locationType: AWS Region
// 					maxIopsBurstPerformance: "3000 for volumes <= 1 TiB"
// 					maxIopsvolume: 16000
// 					maxThroughputvolume: 250 MiB/s
// 					maxVolumeSize:16 TiB
// 					operation:
// 					servicecode: AmazonEC2
// 					servicename: Amazon Elastic Compute Cloud
// 					storageMedia: SSD-backed
// 					usagetype:EU-EBS:VolumeUsage.gp2
// 					volumeApiName:gp2
// 					volumeType:General Purpose
// 				]
// 				productFamily: Storage
// 				sku:GDX844Q9TQAG2YZ2
// 			]
// 			publicationDate: 2020-05-08T21:47:16Z
// 			serviceCode: AmazonEC2
// 			terms:
// 				map[
// 					OnDemand:
// 						map[
// 							GDX844Q9TQAG2YZ2.JRTCKXETXF:
// 								map[
// 									effectiveDate: 2020-05-01T00:00:00Z
// 									offerTermCode: JRTCKXETXF
// 									priceDimensions:
// 										map[
// 											GDX844Q9TQAG2YZ2.JRTCKXETXF.6YS6EN2CT7:
// 												map[
// 													appliesTo: []
// 													beginRange:0
// 													description: "$0.11 per GB-month of General Purpose SSD (gp2) provisioned storage - EU (Ireland)"
// 													endRange: Inf
// 													pricePerUnit:
// 														map[
// 															USD:0.1100000000
// 														]
// 														rateCode: GDX844Q9TQAG2YZ2.JRTCKXETXF.6YS6EN2CT7
// 														unit:GB-Mo
// 												]
// 										]
// 									sku: GDX844Q9TQAG2YZ2
// 									termAttributes:
// 										map[

// 										]
// 								]
// 						]
// 				]
// 			version: 20200508214716
// 	]
// ]
