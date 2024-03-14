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

func GetAllPostponedWallposts(vk_user *api.VK, domain string) ([]object.WallWallpost, error) {
	// TODO: Create some storage for postponed posts.
	var allPosts [][]object.WallWallpost
	var offset int
	for {
		response, err := vk_user.WallGet(api.Params{"domain": domain, "offset": offset, "filter": "postponed"})
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		allPosts = append(allPosts, response.Items)
		if len(allPosts)*100 >= response.Count {
			break
		}
		offset += 100
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
