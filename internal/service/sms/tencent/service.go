package tencent

import (
	"context"
	"fmt"

	isms "github.com/tsukaychan/webook/internal/service/sms"
	"github.com/tsukaychan/webook/pkg/ratelimit"

	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
	limiter  ratelimit.Limiter
}

const LimitKey = "sms:tencent"

func (s *Service) Send(ctx context.Context, tplId string, args []isms.ArgVal, phones ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = &tplId
	req.PhoneNumberSet = s.strToStringPtrSlice(phones)
	req.TemplateParamSet = s.argToStringPtrSlice(args)

	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	// TODO multi error
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("send sms failed, code: %s, reason: %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) strToStringPtrSlice(values []string) []*string {
	var res []*string
	for i := range values {
		res = append(res, &values[i])
	}
	return res
}

func (s *Service) argToStringPtrSlice(values []isms.ArgVal) []*string {
	var res []*string
	for i := range values {
		res = append(res, &values[i].Val)
	}
	return res
}

func NewService(client *sms.Client, appId string, signName string, limiter ratelimit.Limiter) *Service {
	return &Service{
		appId:    &appId,
		signName: &signName,
		client:   client,
		limiter:  limiter,
	}
}
