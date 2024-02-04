package rpc

import (
	"context"
	"errors"
	"io"

	pb "github.com/tempcke/rpm/api/rpc/proto"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/specifications"
	"github.com/tempcke/rpm/usecase"
)

var _ specifications.Driver = Driver{}

// Driver is used to run specifications tests against a stood up server
// however in doing so the code also becomes a client wrapper with an identical interface
// to the httpDriver
type Driver struct {
	client pb.RPMClient
}

func NewDriver(client pb.RPMClient) Driver {
	return Driver{
		client: client,
	}
}
func (d Driver) StoreProperty(ctx context.Context, p entity.Property) (entity.ID, error) {
	client, err := d.getClient()
	if err != nil {
		return "", err
	}
	req := pb.StorePropertyReq{
		Property: pb.ToProperty(p),
	}
	res, err := client.StoreProperty(ctx, &req)
	if err != nil {
		return "", err
	}
	return res.PropertyID, nil
}
func (d Driver) GetProperty(ctx context.Context, id entity.ID) (*entity.Property, error) {
	client, err := d.getClient()
	if err != nil {
		return nil, err
	}
	req := pb.GetPropertyReq{PropertyID: id}
	res, err := client.GetProperty(ctx, &req)
	if err != nil {
		return nil, err
	}
	p := entity.Property{
		ID:        res.Property.GetPropertyID(),
		Street:    res.Property.GetStreet(),
		City:      res.Property.GetCity(),
		StateCode: res.Property.GetState(),
		Zip:       res.Property.GetZip(),
	}
	return &p, nil
}
func (d Driver) ListProperties(ctx context.Context, f usecase.PropertyFilter) ([]entity.Property, error) {
	client, err := d.getClient()
	if err != nil {
		return nil, err
	}
	req := pb.FromPropertyFilter(f)
	stream, err := client.ListProperties(ctx, req)
	if err != nil {
		return nil, err
	}

	var properties = make([]entity.Property, 0)
	for {
		p, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		properties = append(properties, entity.Property{
			ID:        p.GetPropertyID(),
			Street:    p.GetStreet(),
			City:      p.GetCity(),
			StateCode: p.GetState(),
			Zip:       p.GetZip(),
		})
	}
	return properties, nil
}
func (d Driver) RemoveProperty(ctx context.Context, id entity.ID) error {
	client, err := d.getClient()
	if err != nil {
		return err
	}
	_, err = client.RemoveProperty(ctx, &pb.RemovePropertyReq{PropertyID: id})
	if err != nil {
		return err
	}
	return nil
}

func (d Driver) StoreTenant(ctx context.Context, tenant entity.Tenant) (*entity.Tenant, error) {
	client, err := d.getClient()
	if err != nil {
		return nil, err
	}
	req := pb.StoreTenantReq{
		Tenant: pb.ToTenant(tenant),
	}
	res, err := client.StoreTenant(ctx, &req)
	if err != nil {
		return nil, err
	}
	tenant.ID = res.TenantID
	return &tenant, nil
}
func (d Driver) GetTenant(ctx context.Context, id entity.ID) (*entity.Tenant, error) {
	client, err := d.getClient()
	if err != nil {
		return nil, err
	}
	req := pb.GetTenantReq{TenantID: id}
	res, err := client.GetTenant(ctx, &req)
	if err != nil {
		return nil, err
	}
	return res.Tenant.ToTenant().Ptr(), nil
}
func (d Driver) ListTenants(ctx context.Context) ([]entity.Tenant, error) {
	client, err := d.getClient()
	if err != nil {
		return nil, err
	}
	stream, err := client.ListTenants(ctx, &pb.ListTenantsReq{})
	if err != nil {
		return nil, err
	}
	var tenants []entity.Tenant
	for {
		pbTenant, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		tenants = append(tenants, pbTenant.ToTenant())
	}
	return tenants, nil
}

func (d Driver) getClient() (pb.RPMClient, error) {
	if d.client == nil {
		return nil, errors.New("client not initialized")
	}
	return d.client, nil
}
