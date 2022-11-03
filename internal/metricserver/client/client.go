package metricserver

import (
	"context"
	"errors"

	"github.com/kmx0/devops/internal/metricserver/proto"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// SendMetricsBatch ...
func SendMetricsBatch(ctx context.Context, conn *grpc.ClientConn, metrics []types.Metrics) error {
	if len(metrics) == 0 {
		return errors.New("массив метрик пуст")
	}
	resp := make([]*proto.Metric, 0)
	for _, m := range metrics {
		msg := convertToProto(m)
		if msg == nil {
			return errors.New("неверный тип метрики")
		}
		resp = append(resp, msg)
	}
	client := proto.NewMetricsServiceClient(conn)
	// client.Update()
	logrus.Info("Updating: ", metrics)
	_, err := client.UpdateMetricBatch(ctx, &proto.UpdateMetricBatchRequest{Metrics: resp})
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// SendMetrics ...
func SendMetricSingle(ctx context.Context, conn *grpc.ClientConn, metric types.Metrics) error {
	// resp := make([]*proto.Metric, 0)
	msg := convertToProto(metric)
	if msg == nil {
		return errors.New("неверный тип метрики")
	}
	// resp = append(resp, msg)
	client := proto.NewMetricsServiceClient(conn)
	// client.Update()
	logrus.Info("Updating single metric: ", metric)
	_, err := client.UpdateMetric(ctx, &proto.UpdateMetricRequest{Metric: msg})
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// SendMetrics ...
func GetMetric(ctx context.Context, conn *grpc.ClientConn, metric types.Metrics) (*proto.Metric, error) {
	// resp := make([]*proto.Metric, 0)

	msg := convertToProto(metric)

	if msg == nil {
		return nil, errors.New("неверный тип метрики")
	}

	client := proto.NewMetricsServiceClient(conn)
	logrus.Info("Requesting single metric: ", metric)
	resp, err := client.GetMetric(ctx, &proto.GetMetricRequest{Metric: msg})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return resp, nil
}

func convertToProto(m types.Metrics) *proto.Metric {
	metric := &proto.Metric{Id: m.ID, Hash: m.Hash, Type: GetMetricProtoType(&m)}
	switch metric.Type {
	case proto.Type_COUNTER:
		metric.Value = &proto.Metric_Counter{Counter: *m.Delta}
	case proto.Type_GAUGE:
		metric.Value = &proto.Metric_Gauge{Gauge: *m.Value}
	default:
		return nil
	}
	return metric
}

func GetMetricProtoType(m *types.Metrics) proto.Type {
	switch m.MType {
	case "counter":
		return proto.Type_COUNTER
	case "gauge":
		return proto.Type_GAUGE
	default:
		return proto.Type_UNKNOWN
	}
}
