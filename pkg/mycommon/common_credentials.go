package mycommon

import (
	"strconv"

	"github.com/ai-flowx/drivex/pkg/config"
	"github.com/ai-flowx/drivex/pkg/mycomdef"
)

func GetACredentials(s *config.ModelDetails) (cred map[string]interface{}, id string) {
	if len(s.CredentialList) > 0 {
		key := s.ServiceID + "credentials"
		index := config.GetLBIndex(config.LoadBalancingStrategy, key, len(s.CredentialList))
		id = s.ServiceID + "_credentials_" + strconv.Itoa(index)
		return s.CredentialList[index], id
	}

	return s.Credentials, id
}

func GetCredentialLimit(credentials map[string]interface{}) (limitType string, limitn float64, timeout int) {
	limitData, ok := credentials["limit"].(map[string]interface{})
	if !ok {
		return "", 0, 0 // 没有找到或类型不匹配
	}

	if to, ok := limitData["timeout"].(int); ok {
		timeout = to
	}

	if qps, ok := limitData[mycomdef.KeynameQps].(float64); ok {
		return mycomdef.KeynameQps, qps, timeout
	}

	if qpm, ok := limitData[mycomdef.KeynameQpm].(float64); ok {
		return mycomdef.KeynameQpm, qpm, timeout
	}

	if rpm, ok := limitData[mycomdef.KeynameRpm].(float64); ok {
		return mycomdef.KeynameQpm, rpm, timeout
	}

	if concurrency, ok := limitData[mycomdef.KeynameConcurrency].(float64); ok {
		return mycomdef.KeynameConcurrency, concurrency, timeout
	}

	return "", 0, 0 // 默认返回
}
