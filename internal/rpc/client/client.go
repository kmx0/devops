package rpc

import (
	"context"
	"errors"

	"github.com/kmx0/devops/internal/rpc/proto"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// SendMetrics ...
func SendMetrics(ctx context.Context, conn *grpc.ClientConn, metrics []types.Metrics) error {
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
	client := proto.NewAlertingClient(conn)
	// client.Update()
	logrus.Info("Updating: ", metrics)
	_, err := client.Update(ctx, &proto.UpdateRequest{Metrics: resp})
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
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
