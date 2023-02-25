package rpc

import (
	"context"
	"errors"

	"github.com/tempcke/rpm/actions"
	pb "github.com/tempcke/rpm/api/rpc/proto"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.RPMServer = (*Server)(nil)

type Server struct {
	Conf    config.Config
	actions actions.Actions

	pb.UnimplementedRPMServer
}

func NewServer(actions actions.Actions) *Server {
	server := Server{
		Conf:    config.GetConfig(),
		actions: actions,
	}
	return &server
}
func (s Server) WithConfig(conf config.Config) *Server {
	s.Conf = conf
	return &s
}

func (s *Server) StoreProperty(ctx context.Context, req *pb.StorePropertyReq) (*pb.StorePropertyRes, error) {
	pIn := entity.Property{
		ID:        req.PropertyID,
		Street:    req.Street,
		City:      req.City,
		StateCode: req.State,
		Zip:       req.Zip,
	}
	id, err := s.actions.AddRental(ctx, pIn)
	if err != nil {
		// FIXME: return a proper grpc error status
		return nil, err
	}

	res := pb.StorePropertyRes{PropertyID: id}
	return &res, nil
}

func (s *Server) GetProperty(ctx context.Context, req *pb.GetPropertyReq) (*pb.GetPropertyRes, error) {
	var propertyID = req.GetPropertyID()
	p, err := s.actions.GetProperty(ctx, propertyID)
	if err != nil {
		if errors.Is(err, internal.ErrEntityNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	res := pb.GetPropertyRes{
		Property: &pb.Property{
			PropertyID: p.ID,
			Street:     p.Street,
			City:       p.City,
			State:      p.StateCode,
			Zip:        p.Zip,
		},
	}
	return &res, nil
}

func (s *Server) RemoveProperty(ctx context.Context, req *pb.RemovePropertyReq) (*pb.RemovePropertyRes, error) {
	if err := s.actions.RemoveProperty(ctx, req.GetPropertyID()); err != nil {
		return nil, err // FIXME: use proper status error
	}
	res := &pb.RemovePropertyRes{}
	return res, nil
}

func (s *Server) ListProperties(filter *pb.PropertyFilter, stream pb.RPM_ListPropertiesServer) error {
	var ctx = context.Background() // TODO: is there a better context to use?
	properties, err := s.actions.ListProperties(ctx)
	if err != nil {
		return err // FIXME: use proper status error
	}
	for _, p := range properties {
		property := pb.Property{
			PropertyID: p.ID,
			Street:     p.Street,
			City:       p.City,
			State:      p.StateCode,
			Zip:        p.Zip,
		}
		if err := stream.Send(&property); err != nil {
			return err
		}
	}
	return nil
}
