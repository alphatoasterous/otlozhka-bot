package handlers

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"log"
	"time"
)

type WallpostStorage struct {
	timestamp int64
	keepAlive int64

	wallPosts []object.WallWallpost
}

func NewWallpostStorage(keepAlive int64) *WallpostStorage {
	return &WallpostStorage{
		timestamp: 0,
		keepAlive: keepAlive,
	}
}

func (wpStorage *WallpostStorage) GetWallposts() []object.WallWallpost {
	log.Println("WPStorage: Getting Wallposts from storage... Stored posts amount:", len(wpStorage.wallPosts))
	return wpStorage.wallPosts
}

func (wpStorage *WallpostStorage) UpdateWallpostStorage(vkUser *api.VK, domain string) {
	currentTimestamp := time.Now().Unix()
	if currentTimestamp-wpStorage.timestamp >= wpStorage.keepAlive {
		log.Println("WPStorage: Wallposts in wallpost storage are stale, updating...")
		postponedPosts, err := GetAllPostponedWallposts(vkUser, domain)
		if err != nil {
			log.Fatal(err)
		}
		wpStorage.wallPosts = postponedPosts
		wpStorage.timestamp = currentTimestamp
	} else {
		log.Println("WPStorage: Wallposts in wallpost storage are not stale, skipping update...")
	}
}

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

func GetWallpostsByPeerID(peerID int, posts []object.WallWallpost) []object.WallWallpost {
	var foundPosts []object.WallWallpost
	for _, post := range posts {
		if post.SignerID == peerID {
			foundPosts = append(foundPosts, post)
		}
	}
	return foundPosts
}
