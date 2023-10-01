package registry

import (
	"context"
	"github.com/Untanky/modern-auth/internal/core"
	"github.com/google/uuid"
)

type registryServer struct {
	store core.KeyValueStore[string, *RegistrationInfo]
	index core.KeyValueStore[string, core.List[string]]

	registerChan   chan<- *RegistrationInfo
	unregisterChan chan<- *RegistrationInfo
}

func NewRegistryServer(
	store core.KeyValueStore[string, *RegistrationInfo],
	index core.KeyValueStore[string, core.List[string]],
	registerChan chan<- *RegistrationInfo,
	unregisterChan chan<- *RegistrationInfo,
) RegistryServer {
	return &registryServer{
		store:          store,
		index:          index,
		registerChan:   registerChan,
		unregisterChan: unregisterChan,
	}
}

func (r registryServer) Register(ctx context.Context, info *RegistrationInfo) (*RegistrationResponse, error) {
	id := uuid.New().String()
	info.Id = id
	err := r.store.WithContext(ctx).Set(id, info)
	if err != nil {
		return nil, err
	}
	list, err := r.index.WithContext(ctx).Get(info.GetName())
	if err != nil {
		return nil, err
	}
	_, err = list.Append(id)
	if err != nil {
		return nil, err
	}

	r.registerChan <- info

	response := &RegistrationResponse{
		Id:    id,
		Token: "abc",
	}
	return response, nil
}

func (r registryServer) Unregister(ctx context.Context, response *RegistrationResponse) (*Empty, error) {
	store := r.store.WithContext(ctx)
	info, err := store.Get(response.Id)
	if err != nil {
		return nil, err
	}

	err = r.store.WithContext(ctx).Delete(response.Id)
	if err != nil {
		return nil, err
	}

	list, err := r.index.WithContext(ctx).Get(info.GetName())
	if err != nil {
		return nil, err
	}

	index, err := list.Index(response.Id)
	if err != nil {
		return nil, err
	}

	err = list.Remove(index)
	if err != nil {
		return nil, err
	}

	r.unregisterChan <- info

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
