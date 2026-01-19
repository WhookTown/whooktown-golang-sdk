package whooktown

import (
	"context"

	"github.com/gofrs/uuid"
)

// GroupsClient provides asset group management functionality
type GroupsClient struct {
	http *httpClient
}

// ListGroups returns asset groups for a layout
func (c *GroupsClient) ListGroups(ctx context.Context, layoutID uuid.UUID) ([]AssetGroup, error) {
	var groups []AssetGroup
	if err := c.http.Get(ctx, "/ui/groups/"+layoutID.String(), &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// CreateGroupRequest represents a request to create an asset group
type CreateGroupRequest struct {
	LayoutID uuid.UUID `json:"layout_id"`
	Name     string    `json:"name"`
}

// CreateGroup creates a new asset group
func (c *GroupsClient) CreateGroup(ctx context.Context, req *CreateGroupRequest) (*AssetGroup, error) {
	var group AssetGroup
	if err := c.http.Post(ctx, "/ui/groups", req, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

// UpdateGroupRequest represents a request to update an asset group
type UpdateGroupRequest struct {
	Name string `json:"name,omitempty"`
}

// UpdateGroup updates an asset group
func (c *GroupsClient) UpdateGroup(ctx context.Context, groupID uuid.UUID, req *UpdateGroupRequest) (*AssetGroup, error) {
	var group AssetGroup
	if err := c.http.Put(ctx, "/ui/groups/"+groupID.String(), req, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

// DeleteGroup deletes an asset group
func (c *GroupsClient) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	return c.http.Delete(ctx, "/ui/groups/"+groupID.String())
}

// AddMember adds a building to an asset group
func (c *GroupsClient) AddMember(ctx context.Context, groupID, buildingID uuid.UUID) (*AssetGroup, error) {
	body := map[string]uuid.UUID{
		"building_id": buildingID,
	}
	var group AssetGroup
	if err := c.http.Post(ctx, "/ui/groups/"+groupID.String()+"/members", body, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

// RemoveMember removes a building from an asset group
func (c *GroupsClient) RemoveMember(ctx context.Context, groupID, buildingID uuid.UUID) (*AssetGroup, error) {
	var group AssetGroup
	if err := c.http.Delete(ctx, "/ui/groups/"+groupID.String()+"/members/"+buildingID.String()); err != nil {
		return nil, err
	}
	return &group, nil
}
