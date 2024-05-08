package api_utils

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/rs/zerolog/log"
)

// GetGroupInfo retrieves information about the community/group page using a `*api.VK* instance with community access.
// It makes an API call to vkCommunity.GroupsGetByID and returns the group information.
// If an error occurs during the API call, the function logs the error and terminates execution.
func GetGroupInfo(vkCommunity *api.VK) api.GroupsGetByIDResponse {
	group, err := vkCommunity.GroupsGetByID(nil)
	if err != nil {
		log.Fatal().Err(err)
	}
	return group
}

// IsManagerWithRights checks if a given role is associated with managerial rights.
// The function returns true for the roles "editor", "administrator", and "creator".
// It returns false for all other roles.
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

// GetGroupManagerIDs retrieves the IDs of group managers with managerial permissions from VK.
// It uses the vkUser client to make an API call with the specified domain as the group_id.
// The function filters the members list to include only those with managerial rights.
// Returns a slice of IDs or logs a fatal error if the API call fails.
// It should be noted, that vkUser client should have sufficient permissions in given domain or else things go south.
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
