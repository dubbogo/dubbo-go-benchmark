package metrics_collector

import (
	"context"
	"strconv"
	"time"
)

import (
	clusterutils "dubbo.apache.org/dubbo-go/v3/cluster/utils"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
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
	AdaptiveServiceRemainingGauge = "adaptive_service_remaining"
	AdaptiveServiceInflightGauge  = "adaptive_service_inflight"
	LabelProtocol                 = "protocol"
	LabelMethod                   = "method"
	LabelTargetIP                 = "target_ip"
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
	labels := LabelMap(invoker.GetURL().Protocol, invocation.MethodName(), invoker.GetURL().Ip)
	prometheus.IncCounterWithLabel(RequestCounter, labels)

	startTime := time.Now()
	result := invoker.Invoke(ctx, invocation)
	result.AddAttachment(StartTimeAttachment, startTime)

	return result
}

func (f *MetricsCollector) OnResponse(_ context.Context, result protocol.Result, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	labels := LabelMap(invoker.GetURL().Protocol, invocation.MethodName(), invoker.GetURL().Ip)
	startTimeIFace := result.Attachment(StartTimeAttachment, "")
	if startTime, ok := startTimeIFace.(time.Time); ok {
		prometheus.IncSummaryWithLabel(RequestDurationSummary, float64(time.Now().Sub(startTime).Nanoseconds()), labels)
	} else {
		logger.Warnf("StartTime is not a time.Time: %v", startTimeIFace)
	}

	remainingIFace := result.Attachment(constant.AdaptiveServiceRemainingKey, "")
	if remainingStr, ok := remainingIFace.(string); ok {
		if remaining, err := strconv.ParseInt(remainingStr, 10, 64); err == nil {
			prometheus.SetGaugeWithLabel(AdaptiveServiceRemainingGauge, float64(remaining), labels)
		} else {
			logger.Warnf("parse remaining error: %v", err)
		}
	} else {
		logger.Warnf("RemainingAttachment is not a string: %v", remainingIFace)
	}

	inflightIFace := result.Attachment(constant.AdaptiveServiceInflightKey, "")
	if inflightStr, ok := inflightIFace.(string); ok {
		if inflight, err := strconv.ParseInt(inflightStr, 10, 64); err == nil {
			prometheus.SetGaugeWithLabel(AdaptiveServiceInflightGauge, float64(inflight), labels)
		} else {
			logger.Warnf("parse inflight error: %v", err)
		}
	} else {
		logger.Warnf("InflightAttachment is not a string: %v", inflightIFace)
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

func LabelMap(protocol, method, ip string) map[string]string {
	return map[string]string{
		LabelProtocol: protocol,
		LabelMethod:   method,
		LabelTargetIP: ip,
	}
}
