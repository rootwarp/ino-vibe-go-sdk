package device

import pb "bitbucket.org/ino-on/ino-vibe-api"

// Filter is constraints for query device.
type Filter struct {
	InstallStatus isFilterInstallStatus
	GroupID       isFilterGroupID
}

type isFilterInstallStatus interface {
	isFilterInstallStatus()
}

// FilterInstallStatus is install status value of filter.
type FilterInstallStatus struct {
	Value pb.InstallStatus
}

func (f FilterInstallStatus) isFilterInstallStatus() {}

type isFilterGroupID interface {
	isFilterGroupID()
}

// FilterGroupID is group ID value of filter
type FilterGroupID struct {
	Value string
}

func (f FilterGroupID) isFilterGroupID() {}
