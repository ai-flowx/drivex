package mycommon

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/ai-flowx/drivex/pkg/mylog"
)

func CheckStatusCode(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		errMsg, err := io.ReadAll(resp.Body)
		if err != nil {
			mylog.Logger.Error("Failed to read response body",
				zap.Int("status", resp.StatusCode),
				zap.Error(err))
			return errors.New("failed to read error response body")
		}

		mylog.Logger.Error("Unexpected status code",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(errMsg)))

		return fmt.Errorf("status %d: %s", resp.StatusCode, string(errMsg))
	}

	return nil
}
