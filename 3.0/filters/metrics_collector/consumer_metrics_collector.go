package metrics_collector

import (
	"context"
	"os"
	"time"
)

import (
	clusterutils "dubbo.apache.org/dubbo-go/v3/cluster/utils"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/metrics/prometheus"
	"dubbo.apache.org/dubbo-go/v3/protocol"
)

import (
	"github.com/dubbogo/dubbo-go-benchmark/3.0/filters/offline_simulator"
)

const (
	StartTimeAttachment           = "StartTime"
	RequestCounter                = "request_count"
	RequestDurationSummary        = "request_duration_ns"
	RequestSuccessCounter         = "request_success"
	RequestUnknownErrorCounter    = "request_unknown_error"
	RequestTimeoutCounter         = "request_timeout"
	RequestOfflineDroppedCounter  = "request_offline_dropped"
	RequestReachLimitationCounter = "request_reach_limitation"
)

type ConsumerMetricsCollector struct{}

var ErrConsumerRequestTimeoutStr = "maybe the client read timeout or fail to decode tcp stream in Writer.Write"

func init() {
	extension.SetFilter("consumerMetricsCollector", NewConsumerMetricsCollector)
}

func NewConsumerMetricsCollector() filter.Filter {
	return &ConsumerMetricsCollector{}
}

func (f *ConsumerMetricsCollector) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	prometheus.IncCounterWithLabel(RequestCounter, ConsumerLabelMap())
	startTime := time.Now()
	result := invoker.Invoke(ctx, invocation)
	result.AddAttachment(StartTimeAttachment, startTime)
	return result
}

func (f *ConsumerMetricsCollector) OnResponse(_ context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	startTimeIFace := result.Attachment(StartTimeAttachment, "")
	if startTime, ok := startTimeIFace.(time.Time); ok {
		prometheus.IncSummaryWithLabel(RequestDurationSummary, float64(time.Now().Sub(startTime).Nanoseconds()), ConsumerLabelMap())
	} else {
		logger.Warnf("StartTime is not a time.Time: %v", startTimeIFace)
	}

	if result.Error() == nil {
		prometheus.IncCounterWithLabel(RequestSuccessCounter, ConsumerLabelMap())
	} else if clusterutils.DoesAdaptiveServiceReachLimitation(result.Error()) {
		prometheus.IncCounterWithLabel(RequestReachLimitationCounter, ConsumerLabelMap())
	} else if offline_simulator.IsServerOfflineErr(result.Error()) {
		prometheus.IncCounterWithLabel(RequestOfflineDroppedCounter, ConsumerLabelMap())
	} else if context.DeadlineExceeded.Error() == result.Error().Error() || ErrConsumerRequestTimeoutStr == result.Error().Error() {
		prometheus.IncCounterWithLabel(RequestTimeoutCounter, ConsumerLabelMap())
	} else {
		prometheus.IncCounterWithLabel(RequestUnknownErrorCounter, ConsumerLabelMap())
	}
	return result
}

func ConsumerLabelMap() map[string]string {
	return map[string]string{
		"tps":          os.Getenv("TPS"),
		"duration":     os.Getenv("DURATION"),
		"func_name":    os.Getenv("FUNC_NAME"),
		"service_type": "consumer",
	}
}
