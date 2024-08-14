package v1alpha1

import "time"

type SyncPolicy struct {
	Enabled bool          `json:"enabled,required"`
	Refresh time.Duration `json:"refresh,required"`
}
