package ebsavePricing

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/pricing"
)

var priceEu = 0.1100000000
var priceUs = 0.1000000000

var pricingInfo = map[string]interface{}{
	"hoo": "testing",
	"terms": map[string]interface{}{
		"OnDemand": map[string]interface{}{
			"USD": "0.1100000000",
		},
	},
	"xyz": "foo",
}

var mockGetProducts = pricing.GetProductsOutput{
	PriceList: []aws.JSONValue{
		pricingInfo,
	},
}

var PricingStub = &PricingClient{
	getProducts: func(*pricing.GetProductsInput) (*pricing.GetProductsOutput, error) {
		return &mockGetProducts, nil
	},
}

func TestGetVolumeCost(t *testing.T) {
	t.Run("single volume 200GB gp2 eu-west-1", func(t *testing.T) {

		volume := ec2.Volume{
			Size:             aws.Int64(200),
			VolumeType:       aws.String("gp2"),
			AvailabilityZone: aws.String("eu-west-1a"),
		}

		expected := AssetCost{
			Timeframe: "monthly",
			Currency:  "USD",
			Value:     22.00,
		}
		returned, err := PricingStub.GetVolumeCost(&volume)

		if !reflect.DeepEqual(expected, returned) {
			t.Errorf("expected %v but returned %v", expected, returned)
		}

		assertError(t, nil, err)

	})
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
