package accountDetails

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
)

func TestGetCallerIdentity(t *testing.T) {
	t.Run("return account id", func(t *testing.T) {
		var input *sts.GetCallerIdentityInput

		var mockGetCallerIdentity = &sts.GetCallerIdentityOutput{
			Account: aws.String("abc"),
		}

		s := &STSClient{
			getCallerIdentity: func(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
				return mockGetCallerIdentity, nil
			},
		}

		returned, _ := s.GetCallerIdentity(input)

		if returned != mockGetCallerIdentity {
			t.Errorf("expect %v but returned %v", mockGetCallerIdentity, returned)
		}
	})
}
