package collector

import (
	"context"

	authapi "github.com/fabric8-services/fabric8-notification/auth/api"
	"github.com/fabric8-services/fabric8-notification/wit"
	"github.com/fabric8-services/fabric8-notification/wit/api"
)

func NewVersionResolver(authClient *authapi.Client, witClient *api.Client) ReceiverResolver {
	return func(ctx context.Context, url string) ([]Receiver, map[string]interface{}, error) {
		codebases, err := wit.GetCodebases(ctx, witClient, url)
		if err != nil {
			return nil, nil, err
		}
		spaceIDs := collectCodebasesSpaces(codebases)
		if len(spaceIDs) == 0 {
			return []Receiver{}, map[string]interface{}{}, nil
		}
		spaces, err := GetSpaces(ctx, witClient, spaceIDs)
		if err != nil {
			return nil, nil, err
		}

		users := collectSpacesUsers(spaces)
		if len(users) == 0 {
			return []Receiver{}, map[string]interface{}{}, nil
		}
		resolved, err := resolveAllUsers(ctx, authClient, SliceUniq(users), []*authapi.UserData{}, false)
		if err != nil {
			return nil, nil, err
		}
		return resolved, map[string]interface{}{}, nil
	}
}
