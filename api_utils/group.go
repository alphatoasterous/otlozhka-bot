package api_utils

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/rs/zerolog/log"
)

func GetGroupInfo(vkCommunity *api.VK) api.GroupsGetByIDResponse {
	// GetGroupInfo gets information about the community/group page utilizing community token
	group, err := vkCommunity.GroupsGetByID(nil)
	if err != nil {
		log.Fatal().Err(err)
	}
	return group
}

func IsManagerWithRights(role string) bool {
	switch role {
	case
		"editor",
		"administrator",
		"creator":
		return true
	}
	return false
}

func GetGroupManagerIDs(vkUser *api.VK, domain string) []int {
	groupManagers, err := vkUser.GroupsGetMembersFilterManagers(api.Params{"group_id": domain})
	if err != nil {
		log.Fatal().Err(err)
	}
	groupManagerIDs := make([]int, 0, len(groupManagers.Items))
	for _, groupManager := range groupManagers.Items {
		if IsManagerWithRights(groupManager.Role) {
			groupManagerIDs = append(groupManagerIDs, groupManager.ID)
		}
	}
	return groupManagerIDs
}
