package handlers

import (
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/rs/zerolog/log"
)

// WallpostStorage manages the storage and retrieval of wall posts.
// It holds a collection of wall posts and timestamps to manage data freshness.
// Storing wallposts in memory reduces VK API calls and does not affect user experience, which is optimal.
type WallpostStorage struct {
	timestamp int64
	keepAlive int64

	wallPosts []object.WallWallpost
}

// NewWallpostStorage initializes a new WallpostStorage with a specified keepAlive duration.
// The keepAlive parameter determines how long (in seconds) the posts are considered fresh.
// Returns a pointer to the newly created WallpostStorage.
func NewWallpostStorage(keepAlive int64) *WallpostStorage {
	return &WallpostStorage{
		timestamp: 0,
		keepAlive: keepAlive,
	}
}

// GetWallposts retrieves all the wall posts currently stored in WallpostStorage.
// It logs the retrieval process and the number of posts fetched.
// Returns a slice of WallWallpost objects.
func (wpStorage *WallpostStorage) GetWallposts() []object.WallWallpost {
	log.Print("WPStorage: Getting Wallposts from storage... Stored posts amount:", wpStorage.GetWallpostCount())
	return wpStorage.wallPosts
}

// GetWallpostCount returns the number of wall posts currently stored.
func (wpStorage *WallpostStorage) GetWallpostCount() int {
	return len(wpStorage.wallPosts)
}

// CheckWallpostStorageNeedsUpdate checks if the wall posts in the storage are stale based on the keepAlive setting.
// Logs a message indicating whether the posts are stale or not.
// Returns true if the posts are stale and need an update; false otherwise.
func (wpStorage *WallpostStorage) CheckWallpostStorageNeedsUpdate() bool {
	currentTimestamp := time.Now().Unix()
	if currentTimestamp-wpStorage.timestamp >= wpStorage.keepAlive {
		log.Info().Msg("WPStorage: Wallposts in wallpost storage are stale")
		return true
	} else {
		log.Info().Msg("WPStorage: Wallposts in wallpost storage are not stale")
		return false
	}
}

// UpdateWallpostStorage fetches and updates the wall posts from a specified VK domain.
// It calls GetAllPostponedWallposts to retrieve new data, logs critical errors, and updates the internal timestamp.
func (wpStorage *WallpostStorage) UpdateWallpostStorage(vkUser *api.VK, domain string) {

	postponedPosts, err := GetAllPostponedWallposts(vkUser, domain)
	if err != nil {
		log.Fatal().Err(err)
	}
	wpStorage.wallPosts = postponedPosts
	wpStorage.timestamp = time.Now().Unix()

}

// flattenWallpostArray takes a two-dimensional slice of WallWallpost objects and flattens it into a single slice.
// It first calculates the total number of WallWallpost objects across all inner slices to pre-allocate the necessary
// space for the resulting slice. This helps in optimizing the memory allocation during the flattening process.
// After pre-allocation, it appends each inner slice's elements to the flattened slice sequentially.
// The function returns the flattened slice of WallWallpost objects.
// An empty slice is returned if the input is empty or only contains empty inner slices.
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

// GetAllPostponedWallposts retrieves all postponed wall posts for a given community(via `domain` string) using the VK API.
// It uses the `vkUser` client(with user access rights for wall.get access) to fetch posts in batches of up to 100 posts per request until all posts are retrieved.
// The function accepts a `*api.VK` instance representing the VK API User client and a `domain` string to specify the target community.
// It repeatedly calls the `wall.get` method with an increasing offset until all posts are fetched.
// For more information about the method, check VK API documentation page: https://dev.vk.com/wall.get
// This method filters for "postponed" posts using the 'filter' field in the API request parameters.
// On successful retrieval of all posts, the function returns a flat slice of WallWallpost objects.
// If an error occurs during the API calls, it tries to retry five times, while logging the failure.
// If retries fails, it crashes miserably.
// The return includes a slice of all postponed WallWallpost objects and an error, if any occurred.
func GetAllPostponedWallposts(vkUser *api.VK, domain string) ([]object.WallWallpost, error) {
	const maxWallPostCount = 100
	const maxRetries = 5
	const retrySleepTime = time.Second * 2

	var (
		allPosts [][]object.WallWallpost
		offset   int
	)

	// Retry function with panic recovery and a retry limit
	// This was done due to possibility of post count being non-constant value.
	// (e.g. postponed post got deleted/published while executing this function)
	tryFetchingWallposts := func() ([]object.WallWallpost, error) {
		defer func() {
			if r := recover(); r != nil {
				log.Warn().Msg("Recovered from panic, retrying...")
			}
		}()

		for {
			response, err := vkUser.WallGet(api.Params{
				"domain": domain,
				"offset": offset,
				"filter": "postponed",
				"count":  maxWallPostCount,
			})
			if err != nil {
				log.Panic().Err(err).Msg("Failed to fetch wall posts")
				return nil, err
			}

			allPosts = append(allPosts, response.Items)

			// Check if we've fetched all posts
			if offset+maxWallPostCount >= response.Count {
				break
			}

			offset += maxWallPostCount
		}

		// Return flattened wall posts
		return flattenWallpostArray(allPosts), nil
	}

	// Call the retry function, limiting to maxRetries
	var result []object.WallWallpost
	var err error

	for retries := 0; retries < maxRetries; retries++ {
		result, err = tryFetchingWallposts()
		if err == nil {
			break
		}

		log.Warn().Int("attempt", retries+1).Msg("Retrying wallpost fetch due to error")
		if retries == maxRetries-1 {
			log.Fatal().Err(err).Msg("Maximum retry attempts reached. Exiting...")
			return nil, err
		}

		time.Sleep(retrySleepTime)
	}

	return result, nil

}

// GetWallpostsByPeerID filters a slice of WallWallpost objects based on the SignerID of every WallWallpost object.
// It returns a new slice containing only the posts where the SignerID matches the provided peerID.
// The function takes an integer `peerID` and a slice of WallWallpost objects `posts` as parameters.
// The returned slice may be empty if no posts match the given peerID.
func GetWallpostsByPeerID(peerID int, posts []object.WallWallpost) []object.WallWallpost {
	var foundPosts []object.WallWallpost
	for _, post := range posts {
		if post.SignerID == peerID {
			foundPosts = append(foundPosts, post)
		}
	}
	return foundPosts
}
