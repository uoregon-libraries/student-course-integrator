package translator

import (
	"hash/crc32"
	"time"
)

// DuckIDToUniversityID returns the university id (95 number) for the given
// duckid, or an error if the service can't be reached to do the lookup
func DuckIDToUniversityID(duckid string) (string, error) {
	// Simulate the cost of an API hit
	time.Sleep(time.Millisecond * 50)

	var ids = []string{"95x000001", "95x000002", "95x000003", "95x000004", "95x000005"}
	var i = crc32.ChecksumIEEE([]byte(duckid)) % uint32(len(ids))
	return ids[i], nil
}

// UniversityIDToDuckID returns the duckid for the given university id (95
// number), or an error if the service can't be reached to do the lookup
func UniversityIDToDuckID(uid string) (string, error) {
	// Simulate the cost of an API hit
	time.Sleep(time.Millisecond * 50)

	var ids = []string{"tester1", "tester2", "tester3", "tester4", "tester5"}
	var i = crc32.ChecksumIEEE([]byte(uid)) % uint32(len(ids))
	return ids[i], nil
}
