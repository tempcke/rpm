package rpc

import (
	"context"
	"errors"

	"github.com/tempcke/rpm/actions"
	pb "github.com/tempcke/rpm/api/rpc/proto"
	"github.com/tempcke/rpm/internal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.RPMServer = (*Server)(nil)

type Server struct {
	pb.UnimplementedRPMServer
	actions actions.Actions
}

func NewServer(actions actions.Actions) *Server {
	server := Server{
		actions: actions,
	}
	return &server
}

func (s *Server) StoreProperty(ctx context.Context, req *pb.StorePropertyReq) (*pb.StorePropertyRes, error) {
	pIn := req.GetProperty().ToProperty()
	id, err := s.actions.StoreProperty(ctx, pIn)
	if err != nil {
		// FIXME: determine and return correct error code
		return nil, status.Error(codes.Unknown, err.Error())
	}

	res := pb.StorePropertyRes{PropertyID: id}
	return &res, nil
}
func (s *Server) RemoveProperty(ctx context.Context, req *pb.RemovePropertyReq) (*pb.RemovePropertyRes, error) {
	if err := s.actions.RemoveProperty(ctx, req.GetPropertyID()); err != nil {
		return nil, err // FIXME: use proper status error
	}
	res := &pb.RemovePropertyRes{}
	return res, nil
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
func (s *Server) ListProperties(req *pb.ListPropertiesReq, stream pb.RPM_ListPropertiesServer) error {
	var ctx = context.Background() // TODO: is there a better context to use?
	filter := req.ToPropertyFilter()
	properties, err := s.actions.ListProperties(ctx, filter)
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

func (s *Server) StoreTenant(ctx context.Context, req *pb.StoreTenantReq) (*pb.StoreTenantRes, error) {
	in := req.GetTenant().ToTenant()
	out, err := s.actions.StoreTenant(ctx, in)
	if err != nil {
		// FIXME: determine and return correct error code
		return nil, status.Error(codes.Unknown, err.Error())
	}
	res := pb.StoreTenantRes{TenantID: out.ID}
	return &res, nil
}
func (s *Server) GetTenant(ctx context.Context, req *pb.GetTenantReq) (*pb.GetTenantRes, error) {
	out, err := s.actions.GetTenant(ctx, req.TenantID)
	if err != nil {
		if errors.Is(err, internal.ErrEntityNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		// FIXME: determine and return correct error code
		return nil, status.Error(codes.Unknown, err.Error())
	}
	res := pb.GetTenantRes{Tenant: pb.ToTenant(*out)}
	return &res, nil
}
func (s *Server) ListTenants(filter *pb.ListTenantsReq, stream pb.RPM_ListTenantsServer) error {
	var ctx = context.Background() // TODO: is there a better context to use?
	list, err := s.actions.ListTenants(ctx)
	if err != nil {
		// FIXME: determine and return correct error code
		return status.Error(codes.Unknown, err.Error())
	}
	for _, e := range list {
		if err := stream.Send(pb.ToTenant(e)); err != nil {
			return err
		}
	}

	return nil
}
