package api_utils

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"log"
)

func flattenWallpostArray(posts [][]object.WallWallpost) []object.WallWallpost {
	totalLength := 0
	for _, innerArray := range posts {
		totalLength += len(innerArray)
	}
	flattened := make([]object.WallWallpost, 0, totalLength)
	for _, innerArray := range posts {
		flattened = append(flattened, innerArray...)
	}
	return flattened
}

func GetAllPostponedWallposts(vkUser *api.VK, domain string) ([]object.WallWallpost, error) {
	// TODO: Create some storage for postponed posts.
	var allPosts [][]object.WallWallpost
	var offset int
	for {
		const maxWallPostCount = 100
		response, err := vkUser.WallGet(
			api.Params{"domain": domain, "offset": offset, "filter": "postponed", "count": maxWallPostCount})
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		allPosts = append(allPosts, response.Items)
		if len(allPosts)*maxWallPostCount >= response.Count {
			break
		}
		offset += maxWallPostCount
	}
	return flattenWallpostArray(allPosts), nil
}

func FindByFromID(posts []object.WallWallpost, fromID int) []object.WallWallpost {
	var foundPosts []object.WallWallpost
	for _, post := range posts {
		if post.SignerID == fromID {
			foundPosts = append(foundPosts, post)
		}
	}
	return foundPosts
}
