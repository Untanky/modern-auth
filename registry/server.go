package registry

import (
	"context"
	"github.com/Untanky/modern-auth/internal/core"
	"github.com/google/uuid"
)

type registryServer struct {
	store core.KeyValueStore[string, *RegistrationInfo]
	index core.List[string]
}

func NewRegistryServer() RegistryServer {
	return &registryServer{}
}

func (r registryServer) Register(ctx context.Context, info *RegistrationInfo) (*RegistrationResponse, error) {
	id := uuid.New().String()
	err := r.store.WithContext(ctx).Set(id, info)
	if err != nil {
		return nil, err
	}
	_, err = r.index.WithContext(ctx).Append(id)
	if err != nil {
		return nil, err
	}

	response := &RegistrationResponse{
		Id:    id,
		Token: "abc",
	}
	return response, nil
}

func (r registryServer) Unregister(ctx context.Context, response *RegistrationResponse) (*Empty, error) {
	err := r.store.WithContext(ctx).Delete(response.Id)
	if err != nil {
		return nil, err
	}

	index, err := r.index.Index(response.Id)
	if err != nil {
		return nil, err
	}

	err = r.index.Remove(index)
	if err != nil {
		return nil, err
	}

	return &Empty{}, nil
}

func (r registryServer) Subscribe(request *EndpointRequest, server Registry_SubscribeServer) error {
	//TODO implement me
	panic("implement me")
}

func (r registryServer) mustEmbedUnimplementedRegistryServer() {
	//TODO implement me
	panic("implement me")
}
