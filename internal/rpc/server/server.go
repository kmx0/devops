package rpc

import (
	"context"
	"net"

	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/repositories"
	"github.com/kmx0/devops/internal/rpc/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type RPCServer struct {
	trusted  *net.IPNet
	s        repositories.Repository
	g        *grpc.Server
	servOpts []grpc.ServerOption
	key      []byte
	proto.UnimplementedAlertingServer
}

var _ proto.AlertingServer = (*RPCServer)(nil)

func (s *RPCServer) Update(ctx context.Context, req *proto.UpdateRequest) (*proto.Empty, error) {
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
	return &proto.Empty{}, nil
}

func (s *RPCServer) saveMetric(req *proto.Metric, peer string) error {
	m := metrics.Metrics{
		ID:   req.Id,
		Hash: req.Hash,
	}
	switch req.GetType() {
	case proto.Type_COUNTER:
		m.MType = metrics.CounterType
		v := req.GetCounter()
		m.Delta = &v
	case proto.Type_GAUGE:
		m.MType = metrics.GaugeType
		v := req.GetGauge()
		m.Value = &v
	default:
		return status.Error(codes.InvalidArgument, repositories.ErrWrongMetricType.Error())
	}
	if len(s.key) != 0 {
		recived := m.Hash
		err := m.Sign(s.key)
		if err != nil || recived != m.Hash {
			return status.Error(codes.InvalidArgument, "подпись не соответствует ожиданиям")
		}
	}
	if err := s.s.UpdateMetric(context.TODO(), peer, m); err != nil {
		switch err {
		case repositories.ErrWrongMetricURL:
			return status.Error(codes.NotFound, err.Error())
		case repositories.ErrWrongMetricValue:
			return status.Error(codes.InvalidArgument, err.Error())
		case repositories.ErrWrongValueInStorage:
			return status.Error(codes.Unimplemented, err.Error())
		default:
			return status.Error(codes.Internal, err.Error())
		}
	}
	return nil
}

func NewRPCServer(cfg config.Config, store repositories.Repository, listen string) (*RPCServer, error) {
	serv := &RPCServer{
		s: store,
	}
	RPCServer := grpc.NewServer(serv.servOpts...)
	proto.RegisterAlertingServer(RPCServer, serv)
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

	if !trustedSubnet.Contains(net.ParseIP(realIP)) {

		return false
	}
	return true
}
