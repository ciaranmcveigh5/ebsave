package ebsavePricing

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/pricing/pricingiface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/pricing"
)

type AssetCost struct {
	Timeframe string  `json:"timeframe,omitempty"`
	Currency  string  `json:"currency,omitempty"`
	Value     float64 `json:"value,omitempty"`
}

var awsRegion = map[string]string{
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

var ebsName = map[string]string{
	"standard": "Magnetic",
	"gp2":      "General Purpose",
	"io1":      "Provisioned IOPS",
	"st1":      "Throughput Optimized HDD",
	"sc1":      "Cold HDD",
}

type PricingClient struct {
	Client      pricingiface.PricingAPI
	getProducts func(*pricing.GetProductsInput) (*pricing.GetProductsOutput, error)
}

func (p *PricingClient) GetProducts(input *pricing.GetProductsInput) (*pricing.GetProductsOutput, error) {
	if p.getProducts != nil {
		result, err := p.getProducts(input)
		return result, err
	}

	result, err := p.Client.GetProducts(input)
	return result, err
}

func (p *PricingClient) GetVolumeCost(volume *ec2.Volume) AssetCost {
	volumeCost := AssetCost{}
	volumeRegion := getVolumeRegion(volume)

	input := &pricing.GetProductsInput{
		MaxResults:  aws.Int64(1),
		ServiceCode: aws.String("AmazonEC2"),
		Filters: []*pricing.Filter{
			&pricing.Filter{
				Type:  aws.String("TERM_MATCH"),
				Field: aws.String("volumeType"),
				Value: aws.String(ebsName[*volume.VolumeType]),
			},
			&pricing.Filter{
				Type:  aws.String("TERM_MATCH"),
				Field: aws.String("location"),
				Value: aws.String(awsRegion[volumeRegion]),
			},
		},
	}
	result, err := p.GetProducts(input)
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
		return volumeCost
	}

	volumePrice := ""
	pricingMap := result.PriceList[0]["terms"].(map[string]interface{})["OnDemand"].(map[string]interface{})
	getUSDFromMap(pricingMap, &volumePrice)
	volumePriceFloat, _ := strconv.ParseFloat(volumePrice, 64)

	volumeCost.Timeframe = "monthly"
	volumeCost.Currency = "USD"
	volumeCost.Value = (float64(*volume.Size) * volumePriceFloat)

	return volumeCost
}

func getVolumeRegion(volume *ec2.Volume) string {
	stringLength := len(*volume.AvailabilityZone)
	str := *volume.AvailabilityZone
	last := str[stringLength-1:]
	volumeRegion := strings.TrimSuffix(str, last)
	return volumeRegion
}

func getUSDFromMap(m map[string]interface{}, result *string) {
	if *result == "" {
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
}
