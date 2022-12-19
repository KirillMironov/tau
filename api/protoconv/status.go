package protoconv

import (
	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/api"
)

func StatusEntriesToProto(entries []tau.StatusEntry) []*api.StatusEntry {
	protoEntries := make([]*api.StatusEntry, 0, len(entries))

	for _, entry := range entries {
		protoEntries = append(protoEntries, &api.StatusEntry{
			Title: entry.Title,
			Value: entry.Value,
		})
	}

	return protoEntries
}
