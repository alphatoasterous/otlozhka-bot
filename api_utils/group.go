package api_utils

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"log"
)

func GetGroupInfo(vkCommunity *api.VK) api.GroupsGetByIDResponse {
	// GetGroupInfo gets information about the community/group page utilizing community token
	group, err := vkCommunity.GroupsGetByID(nil)
	if err != nil {
		log.Fatal(err)
	}
	return group
}
