package mycommon

import (
	"errors"
	"go.uber.org/zap"

	"github.com/ai-flowx/drivex/pkg/mylog"
)

const (
	adjustmentFloatValue = 0.01
)

var (
	modelParamsMap = map[string]ModelParams{
		// TBD: FIXME
	}
)

type ModelParams struct {
	TemperatureRange Range
	TopPRange        Range
	MaxTokens        int
}

type Range struct {
	Min float32
	Max float32
}

func GetModelParams(modelName string) (ModelParams, error) {
	params, ok := modelParamsMap[modelName]
	if !ok {
		return ModelParams{}, errors.New("unsupported model")
	}

	return params, nil
}

func adjustFloatValue(value, _min, _max float32) float32 {
	if value < 0 {
		value = 0
	}

	if value < _min {
		value = _min + adjustmentFloatValue
	} else if value >= _max {
		value = _max - adjustmentFloatValue
	}

	return value
}

func AdjustParamsToRange(modelName string, temperature, topP float32, maxTokens int) (temp, p float32, tokens int, err error) {
	params, err := GetModelParams(modelName)
	if err != nil {
		return temperature, topP, maxTokens, err
	}

	temp = adjustFloatValue(temperature, params.TemperatureRange.Min, params.TemperatureRange.Max)

	p = adjustFloatValue(topP, params.TopPRange.Min, params.TopPRange.Max)

	if maxTokens < 0 {
		tokens = 0
	}

	if maxTokens > params.MaxTokens {
		tokens = params.MaxTokens
	}

	mylog.Logger.Debug("", zap.Float32("adjusted_temperature", temp),
		zap.Float32("adjusted_topP", p),
		zap.Int("adjusted_maxTokens", tokens))

	return temp, p, tokens, nil
}
