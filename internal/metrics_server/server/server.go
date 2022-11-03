package metrics_server

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/crypto"
	"github.com/kmx0/devops/internal/metrics_server/proto"
	"github.com/kmx0/devops/internal/repositories"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RPCServer struct {
	trusted  *net.IPNet
	store    repositories.Repository
	g        *grpc.Server
	servOpts []grpc.ServerOption
	cfg      config.Config
	proto.UnimplementedMetricsServiceServer
}

var _ proto.MetricsServiceServer = (*RPCServer)(nil)

func (s *RPCServer) UpdateMetricBatch(ctx context.Context, req *proto.UpdateMetricBatchRequest) (*empty.Empty, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "адрес не определён")
	}
	realIP, _, err := net.SplitHostPort(p.Addr.String())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "access denied, bad ip")
	}
	if !CheckTrusted(s.trusted, realIP) {
		return nil, status.Error(codes.InvalidArgument, "access denied, ip not Trusted")
	}
	for _, v := range req.Metrics {
		err = s.saveMetric(v)
		if err != nil {
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
}

func (s *RPCServer) UpdateMetric(ctx context.Context, req *proto.UpdateMetricRequest) (*empty.Empty, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "адрес не определён")
	}
	realIP, _, err := net.SplitHostPort(p.Addr.String())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "access denied, bad ip")
	}
	if !CheckTrusted(s.trusted, realIP) {
		return nil, status.Error(codes.InvalidArgument, "access denied, ip not Trusted")
	}
	err = s.saveMetric(req.Metric)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *RPCServer) GetMetric(ctx context.Context, req *proto.GetMetricRequest) (*proto.Metric, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "адрес не определён")
	}
	realIP, _, err := net.SplitHostPort(p.Addr.String())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "access denied, bad ip")
	}
	if !CheckTrusted(s.trusted, realIP) {
		return nil, status.Error(codes.InvalidArgument, "access denied, ip not Trusted")
	}
	switch req.Metric.Type.String() {
	case "counter":
		delta, err := s.store.GetCounterJSON(req.Metric.Id)
		if err != nil {
			return nil, err
		}
		req.Metric.Value = &proto.Metric_Counter{
			Counter: delta,
		}
		if s.cfg.Key != "" {
			req.Metric.Hash = crypto.Hash(fmt.Sprintf("%s:counter:%d", req.Metric.Id, req.Metric.GetCounter()), s.cfg.Key)
		}
		return req.Metric, nil
	case "gauge":
		value, err := s.store.GetGaugeJSON(req.Metric.Id)
		if err != nil {
			return nil, err
		}
		req.Metric.Value = &proto.Metric_Gauge{
			Gauge: value,
		}
		if s.cfg.Key != "" {

			req.Metric.Hash = crypto.Hash(fmt.Sprintf("%s:gauge:%f", req.Metric.Id, req.Metric.GetGauge()), s.cfg.Key)
		}
		return req.Metric, nil
	default:
		return &proto.Metric{}, nil
	}
}

func (s *RPCServer) saveMetric(req *proto.Metric) error {
	metrics := types.Metrics{
		ID:   req.Id,
		Hash: req.Hash,
	}
	switch req.GetType() {
	case proto.Type_COUNTER:
		metrics.MType = "counter"
		v := req.GetCounter()
		metrics.Delta = &v
	case proto.Type_GAUGE:
		metrics.MType = "gauge"
		v := req.GetGauge()
		metrics.Value = &v
	default:
		return status.Error(codes.InvalidArgument, "Неизвестный тип метрики")
	}
	if metrics.MType == "counter" || metrics.MType == "gauge" {
		err := s.store.UpdateJSON(s.cfg.Key, metrics)

		if err != nil {
			logrus.Error(err)

			switch {
			case strings.Contains(err.Error(), `received nil pointer on Delta`) || strings.Contains(err.Error(), `received nil pointer on Value`):
				return status.Error(codes.InvalidArgument, "Отправлено пустое значение")
			case strings.Contains(err.Error(), `hash sum not matched`):
				return status.Error(codes.InvalidArgument, "Хэш-сумма не совпала")
			default:
				return status.Error(codes.Internal, "")
			}
		} else if s.cfg.StoreInterval == 0 || s.cfg.DBDSN != "" {
			s.store.SaveToDisk(s.cfg)

		}
		logrus.Info("Wrtiting data: ", metrics)

	} else {
		return status.Error(codes.Unimplemented, "")
	}
	return nil
}

func NewRPCServer(cfg config.Config, store repositories.Repository, listen string) (*RPCServer, error) {
	_, subnet, _ := net.ParseCIDR(cfg.TrustedSubnet)
	serv := &RPCServer{
		store:   store,
		cfg:     cfg,
		trusted: subnet,
	}
	RPCServer := grpc.NewServer(serv.servOpts...)
	proto.RegisterMetricsServiceServer(RPCServer, serv)
	serv.g = RPCServer
	if len(listen) != 0 {
		go func() {
			list, err := net.Listen("tcp", listen)
			if err != nil {
				logrus.Error(err.Error())
				return
			}
			defer list.Close()
			if err := RPCServer.Serve(list); err != nil {
				logrus.Error(err.Error())
			}
		}()
	}
	return serv, nil
}

func CheckTrusted(trustedSubnet *net.IPNet, realIP string) bool {

	if trustedSubnet == nil {
		return true
	}
	if !trustedSubnet.Contains(net.ParseIP(realIP)) {

		return false
	}
	return true
}
