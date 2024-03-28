package handlers

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"log"
	"sync"
	"time"
)

type WallpostStorage struct {
	mutex sync.Mutex

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

func (wpStorage *WallpostStorage) GetAndUpdateWallpostStorage(vkUser *api.VK, domain string) []object.WallWallpost {
	wpStorage.mutex.Lock()
	defer wpStorage.mutex.Unlock()

	currentTimestamp := time.Now().Unix()
	if currentTimestamp-wpStorage.timestamp > wpStorage.keepAlive {
		postponedPosts, err := GetAllPostponedWallposts(vkUser, domain)
		if err != nil {
			log.Fatal(err)
		}
		wpStorage.wallPosts = postponedPosts
		wpStorage.timestamp = currentTimestamp
	}
	return wpStorage.wallPosts
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

func getWallpostsByPeerID(peerID int, posts []object.WallWallpost) []object.WallWallpost {
	var foundPosts []object.WallWallpost
	for _, post := range posts {
		if post.SignerID == peerID {
			foundPosts = append(foundPosts, post)
		}
	}
	return foundPosts
}
