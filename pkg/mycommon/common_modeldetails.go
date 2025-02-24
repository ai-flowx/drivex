package mycommon

import (
	"github.com/ai-flowx/drivex/pkg/config"
	"github.com/ai-flowx/drivex/pkg/mycomdef"
)

func getLimitDetails(limit config.Limit) (name string, qps float64, timeout int) {
	switch {
	case limit.QPS > 0:
		return mycomdef.KeynameQps, limit.QPS, limit.Timeout
	case limit.QPM > 0:
		return mycomdef.KeynameQpm, limit.QPM, limit.Timeout
	case limit.RPM > 0:
		return mycomdef.KeynameQpm, limit.RPM, limit.Timeout
	case limit.Concurrency > 0:
		return mycomdef.KeynameConcurrency, limit.Concurrency, limit.Timeout
	default:
		return "", 0, 0
	}
}

func GetServiceModelDetailsLimit(s *config.ModelDetails) (name string, qps float64, timeout int) {
	return getLimitDetails(s.Limit)
}
