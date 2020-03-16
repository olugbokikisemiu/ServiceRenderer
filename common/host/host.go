// Package host provides tools for identifying multiple hosts in a network run.
package host

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// noCommit is the default host postfix in case of no commit information.
const noCommit = "no_commit"

var (
	// generated random id used for differenting host on same machine.
	id primitive.ObjectID
	// host is the generated unique identifier for the host.
	host string
)

func init() {
	id = primitive.NewObjectID()
}

// SetHostCommit sets hostname and commit version to host identifier.
// NOTE: Short version of 9 chars will be used for hostID and commithash.
func SetHostCommit(hostname, commit string) {
	version := commit
	if len(version) < 1 {
		version = noCommit
	}
	if len(version) > 9 {
		version = version[:9]
	}

	host = id.Hex()[len(id.Hex())-9:] + ":" + hostname + ":" + version
}

// ID returns unique identifier for the host.
func ID() string {
	return id.Hex()
}

// Host returns a unique identifier for the host with commit and IP information.
// Use SetHostCommit(hostname, commit) to set host and commit information else
// default host ID will be returned.
//
// Following is the format for a unique host name.
//
// 			ShortHostID:HostnameIP:ShortCommitHash
//
func Host() string {
	return host
}
