package accountDetails

import (
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

type STSClient struct {
	Client            stsiface.STSAPI
	getCallerIdentity func(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error)
}

func (s *STSClient) GetCallerIdentity(input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	if s.getCallerIdentity != nil {
		result, err := s.getCallerIdentity(input)
		return result, err
	}

	result, err := s.Client.GetCallerIdentity(input)
	return result, err
}
