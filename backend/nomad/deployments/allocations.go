package deployments

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/jippi/hashi-ui/backend/nomad/helper"
	"github.com/jippi/hashi-ui/backend/structs"
)

const (
	WatchAllocations   = "NOMAD_WATCH_DEPLOYMENT_ALLOCATIONS"
	UnwatchAllocations = "NOMAD_UNWATCH_DEPLOYMENT_ALLOCATIONS"
	fetchedAllocations = "NOMAD_FETCHED_DEPLOYMENT_ALLOCATIONS"
)

type allocations struct {
	action structs.Action
}

func NewAllocations(action structs.Action) *allocations {
	return &allocations{
		action: action,
	}
}

func (w *allocations) Do(client *api.Client, q *api.QueryOptions) (*structs.Action, error) {
	allocations, meta, err := client.Deployments().Allocations(w.action.Payload.(string), q)
	if err != nil {
		return nil, fmt.Errorf("watch: unable to fetch %s: %s", w.Key(), err)
	}

	if !helper.QueryChanged(q, meta) {
		return nil, nil
	}

	return &structs.Action{
		Index:   meta.LastIndex,
		Payload: allocations,
		Type:    fetchedAllocations,
	}, nil
}

func (w *allocations) Key() string {
	return fmt.Sprintf("/deployment/%s/allocations", w.action.Payload.(string))
}

func (w *allocations) IsMutable() bool {
	return false
}
