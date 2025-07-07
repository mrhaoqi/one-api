package converse

import (
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/songquanpeng/one-api/common/ctxkey"
	"github.com/songquanpeng/one-api/relay/adaptor/aws/utils"
	"github.com/songquanpeng/one-api/relay/meta"
	relaymodel "github.com/songquanpeng/one-api/relay/model"
)

var _ utils.AwsAdapter = new(Adaptor)

type Adaptor struct{}

func (a *Adaptor) ConvertRequest(c *gin.Context, relayMode int, request *relaymodel.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}

	// Get AWS model ID
	awsModelID, exists := GetAwsModelID(request.Model)
	if !exists {
		return nil, errors.Errorf("unsupported model: %s", request.Model)
	}

	// Store the original request and AWS model ID for later use
	c.Set(ctxkey.RequestModel, request.Model)
	c.Set("aws_model_id", awsModelID)
	c.Set(ctxkey.ConvertedRequest, request)

	return request, nil
}

func (a *Adaptor) DoResponse(c *gin.Context, awsCli *bedrockruntime.Client, meta *meta.Meta) (usage *relaymodel.Usage, err *relaymodel.ErrorWithStatusCode) {
	if meta.IsStream {
		return StreamHandler(c, awsCli)
	} else {
		return Handler(c, awsCli)
	}
}


