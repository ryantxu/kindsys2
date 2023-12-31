package playlist

//*******************************************************************************************
// NOTE!!
// This file is exploring generating the JSONSchema from golang... but that is paused for now
// The tests are just about the schema for now
//*******************************************************************************************

// Type of the item.
type ItemType string

// Defines values for ItemType.
const (
	// Deprecated -- will be removed in the future
	//	ItemTypeDashboardById  ItemType = "dashboard_by_id"
	ItemTypeDashboardByTag ItemType = "dashboard_by_tag"
	ItemTypeDashboardByUid ItemType = "dashboard_by_uid"
)

// Item defines model for Item.
type Item struct {
	// Title is an unused property -- it will be removed in the future
	Title string `json:"title,omitempty"`

	// Type of the item.
	Type ItemType `json:"type" jsonschema:"enum=dashboard_by_id,enum=dashboard_by_tag,enum=dashboard_by_uid"`

	// Value depends on type and describes the playlist item.
	//
	//  - dashboard_by_id: The value is an internal numerical identifier set by Grafana. This
	//  is not portable as the numerical identifier is non-deterministic between different instances.
	//  Will be replaced by dashboard_by_uid in the future. (deprecated)
	//  - dashboard_by_tag: The value is a tag which is set on any number of dashboards. All
	//  dashboards behind the tag will be added to the playlist.
	//  - dashboard_by_uid: The value is the dashboard UID
	Value string `json:"value" jsonschema:"minLength=1"`
}

// Spec defines model for Spec.
type Spec struct {
	// Interval sets the time between switching views in a playlist.
	// FIXME: Is this based on a standardized format or what options are available? Can datemath be used?
	Interval string `json:"interval"`

	// The ordered list of items that the playlist will iterate over.
	// FIXME! This should not be optional, but changing it makes the godegen awkward
	Items []Item `json:"items" jsonschema:"minLength=1"`

	// Name of the playlist.
	Name string `json:"name" jsonschema:"minLength=1"`

	// Unique playlist identifier. Generated on creation, either by the
	// creator of the playlist of by the application.
	Uid string `json:"uid"`
}

//----

type PlaylistResource struct {
	// ????
}
