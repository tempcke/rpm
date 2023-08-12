//go:build withDocker
// +build withDocker

package main_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/joho/godotenv"
	pb "github.com/tempcke/rpm/api/rpc/proto"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/specifications"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var _ specifications.Driver = (*grpcDriver)(nil)

// grpcDriver is used to run specifications tests against a stood up server
// however in doing so the code also becomes a client wrapper with an identical interface
// to the httpDriver
type grpcDriver struct {
	Addr     string
	connOnce sync.Once
	conn     *grpc.ClientConn
	connErr  error
	client   pb.RPMClient
}

func (d *grpcDriver) StoreProperty(ctx context.Context, p entity.Property) (entity.ID, error) {
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
func (d *grpcDriver) GetProperty(ctx context.Context, id entity.ID) (*entity.Property, error) {
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
func (d *grpcDriver) ListProperties(ctx context.Context) ([]entity.Property, error) {
	client, err := d.getClient()
	if err != nil {
		return nil, err
	}
	stream, err := client.ListProperties(ctx, &pb.ListPropertiesReq{})
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
func (d *grpcDriver) RemoveProperty(ctx context.Context, id entity.ID) error {
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

func (d *grpcDriver) StoreTenant(ctx context.Context, tenant entity.Tenant) (*entity.Tenant, error) {
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
func (d *grpcDriver) GetTenant(ctx context.Context, id entity.ID) (*entity.Tenant, error) {
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
func (d *grpcDriver) ListTenants(ctx context.Context) ([]entity.Tenant, error) {
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

func (d *grpcDriver) getClient() (pb.RPMClient, error) {
	var (
		certFile      = conf.GetString(internal.EnvServiceCertFile)
		credentialOpt = grpc.WithTransportCredentials(insecure.NewCredentials())
	)

	if file := findCertFile(certFile); file != "" {
		creds, err := credentials.NewClientTLSFromFile(file, "")
		if err != nil {
			return nil, err
		}
		credentialOpt = grpc.WithTransportCredentials(creds)
	}
	d.connOnce.Do(func() {
		dialOpts := []grpc.DialOption{credentialOpt}
		d.conn, d.connErr = grpc.Dial(d.Addr, dialOpts...)
		d.client = pb.NewRPMClient(d.conn)
	})
	if err := d.connErr; err != nil {
		return nil, err
	}
	return d.client, nil
}

func findCertFile(relPath string) string {
	if relPath == "" {
		return ""
	}
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(relPath); err == nil {
			_ = godotenv.Load(relPath)
			return relPath
		}
		relPath = fmt.Sprintf("../%s", relPath)
	}
	return ""
}
