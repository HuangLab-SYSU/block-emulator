package measure

import (
	"blockEmulator/message"
)

type MeasureModule interface {
	UpdateMeasureRecord(*message.BlockInfoMsg)
	HandleExtraMessage([]byte)
	OutputMetricName() string
	OutputMetricTitle() string
	OutputRecord() ([]float64, float64)
}
