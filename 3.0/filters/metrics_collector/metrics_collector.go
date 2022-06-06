package metrics_collector

import (
	"context"
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
	LabelProtocol                 = "protocol"
	LabelMethod                   = "method"
)

type MetricsCollector struct{}

var ErrConsumerRequestTimeoutStr = "maybe the client read timeout or fail to decode tcp stream in Writer.Write"

func init() {
	extension.SetFilter("metricsCollector", NewMetricsCollector)
}

func NewMetricsCollector() filter.Filter {
	return &MetricsCollector{}
}

func (f *MetricsCollector) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	labels := LabelMap(invoker.GetURL().Protocol, invocation.MethodName())
	prometheus.IncCounterWithLabel(RequestCounter, labels)

	startTime := time.Now()
	result := invoker.Invoke(ctx, invocation)
	result.AddAttachment(StartTimeAttachment, startTime)

	return result
}

func (f *MetricsCollector) OnResponse(_ context.Context, result protocol.Result, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	labels := LabelMap(invoker.GetURL().Protocol, invocation.MethodName())

	startTimeIFace := result.Attachment(StartTimeAttachment, "")
	if startTime, ok := startTimeIFace.(time.Time); ok {
		prometheus.IncSummaryWithLabel(RequestDurationSummary, float64(time.Now().Sub(startTime).Nanoseconds()), labels)
	} else {
		logger.Warnf("StartTime is not a time.Time: %v", startTimeIFace)
	}

	if result.Error() == nil {
		prometheus.IncCounterWithLabel(RequestSuccessCounter, labels)
	} else if clusterutils.DoesAdaptiveServiceReachLimitation(result.Error()) {
		prometheus.IncCounterWithLabel(RequestReachLimitationCounter, labels)
	} else if offline_simulator.IsServerOfflineErr(result.Error()) {
		prometheus.IncCounterWithLabel(RequestOfflineDroppedCounter, labels)
	} else if context.DeadlineExceeded.Error() == result.Error().Error() || ErrConsumerRequestTimeoutStr == result.Error().Error() {
		prometheus.IncCounterWithLabel(RequestTimeoutCounter, labels)
	} else {
		prometheus.IncCounterWithLabel(RequestUnknownErrorCounter, labels)
	}

	return result
}

func LabelMap(protocol, method string) map[string]string {
	return map[string]string{
		LabelProtocol: protocol,
		LabelMethod:   method,
	}
}
