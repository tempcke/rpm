package rpc_test

import (
	"context"
	"io"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/actions"
	"github.com/tempcke/rpm/api/rpc"
	pb "github.com/tempcke/rpm/api/rpc/proto"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/entity/fake"
	"github.com/tempcke/rpm/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

var ctx = context.Background()

func TestRPC_Property(t *testing.T) {
	var (
		repo      = repository.NewInMemoryRepo()
		server    = rpc.NewServer(actions.NewActionsWithRepo(repo))
		rpmClient = newClient(t, server)
		p1        = fake.Property()
	)

	// StoreProperty
	storeReq := pb.StorePropertyReq{
		Property: pb.ToProperty(p1),
	}
	storeRes, err := rpmClient.StoreProperty(ctx, &storeReq)
	require.NoError(t, err)
	require.NotNil(t, storeRes)
	require.Equal(t, p1.ID, storeRes.PropertyID)

	// GetProperty
	getRes, err := rpmClient.GetProperty(ctx, &pb.GetPropertyReq{PropertyID: p1.GetID()})
	require.NoError(t, err)
	assertPropertyMatch(t, p1, getRes.Property)

	// ListProperties
	stream, err := rpmClient.ListProperties(ctx, &pb.ListPropertiesReq{})
	require.NoError(t, err)
	var properties = make(map[string]*pb.Property, 0)
	for {
		p, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		properties[p.PropertyID] = p
	}
	assert.Contains(t, properties, p1.ID)
	assertPropertyMatch(t, p1, properties[p1.ID])

	// RemoveProperty
	remRes, err := rpmClient.RemoveProperty(ctx, &pb.RemovePropertyReq{PropertyID: p1.ID})
	require.NoError(t, err)
	require.NotNil(t, remRes)
	_, err = rpmClient.GetProperty(ctx, &pb.GetPropertyReq{PropertyID: p1.ID})
	require.Error(t, err)
	s := status.Convert(err)
	require.NotNil(t, s)
	assert.Equal(t, codes.NotFound, s.Code())
	// t.Log(s.String())      // rpc error: code = NotFound desc = entity not found: property[some-id]
	// t.Log(s.Message())     // entity not found: property[some-id]
	// t.Log(s.Err().Error()) // rpc error: code = NotFound desc = entity not found: property[some-id]
}

func TestRPC_Tenant(t *testing.T) {
	var (
		repo      = repository.NewInMemoryRepo()
		server    = rpc.NewServer(actions.NewActionsWithRepo(repo))
		rpmClient = newClient(t, server)
	)

	t.Run("success", func(t *testing.T) {
		// store in1
		in1 := fake.Tenant()
		storeReq := pb.StoreTenantReq{
			Tenant: pb.ToTenant(in1),
		}
		storeRes, err := rpmClient.StoreTenant(ctx, &storeReq)
		require.NoError(t, err)
		require.NotNil(t, storeRes)
		require.Equal(t, in1.ID, storeRes.TenantID)

		// get
		getReq := pb.GetTenantReq{TenantID: in1.GetID()}
		getRes, err := rpmClient.GetTenant(ctx, &getReq)
		require.NoError(t, err)
		require.NotNil(t, getRes)
		out := getRes.GetTenant().ToTenant()
		require.True(t, out.Equal(in1))

		// store in2
		in2 := fake.Tenant()
		_, err = rpmClient.StoreTenant(ctx,
			&pb.StoreTenantReq{Tenant: pb.ToTenant(in2)})
		require.NoError(t, err)

		// list, expect in1, in2
		listReq := pb.ListTenantsReq{}
		stream, err := rpmClient.ListTenants(ctx, &listReq)
		require.NoError(t, err)
		var list = make(map[entity.ID]entity.Tenant)
		for {
			pbTenant, err := stream.Recv()
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
			list[pbTenant.GetTenantID()] = pbTenant.ToTenant()
		}
		require.Len(t, list, 2)
		require.True(t, in1.Equal(list[in1.GetID()]))
		require.True(t, in2.Equal(list[in2.GetID()]))
	})

	t.Run("get tenant not found", func(t *testing.T) {
		getReq := pb.GetTenantReq{TenantID: entity.NewID()}
		getRes, err := rpmClient.GetTenant(ctx, &getReq)
		require.Nil(t, getRes)
		require.Error(t, err)
		s, ok := status.FromError(err)
		require.True(t, ok, "err was not a grpc status")
		assert.Equal(t, codes.NotFound.String(), s.Code().String(), err)
	})
}

func assertPropertyMatch(t *testing.T, expect entity.Property, actual *pb.Property) {
	t.Helper()
	require.NotNil(t, actual)
	assert.Equal(t, expect.ID, actual.PropertyID)
	assert.Equal(t, expect.Street, actual.Street)
	assert.Equal(t, expect.City, actual.City)
	assert.Equal(t, expect.StateCode, actual.State)
	assert.Equal(t, expect.Zip, actual.Zip)

}

func newClient(t testing.TB, server *rpc.Server) pb.RPMClient {
	// start server
	const bufSize = 1024 * 1024
	var (
		lis = bufconn.Listen(bufSize)
		s   = grpc.NewServer()
	)

	pb.RegisterRPMServer(s, server)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	dialOpts := []grpc.DialOption{
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.DialContext(ctx, "bufnet", dialOpts...)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	return pb.NewRPMClient(conn)
}
