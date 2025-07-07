package aws

import (
	"github.com/songquanpeng/one-api/relay/adaptor/aws/converse"
	"github.com/songquanpeng/one-api/relay/adaptor/aws/utils"
)

func GetAdaptor(model string) utils.AwsAdapter {
	// Check if model is supported by the unified Converse adapter
	if _, exists := converse.GetAwsModelID(model); exists {
		return &converse.Adaptor{}
	}
	return nil
}
