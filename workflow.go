package whooktown

import (
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
)

// WorkflowClient provides workflow management functionality
type WorkflowClient struct {
	http *httpClient
}

// List returns all workflows for the account
func (c *WorkflowClient) List(ctx context.Context) ([]Workflow, error) {
	var workflows []Workflow
	if err := c.http.Get(ctx, "/workflow", &workflows); err != nil {
		return nil, err
	}
	return workflows, nil
}

// CreateWorkflowRequest represents a request to create a workflow
type CreateWorkflowRequest struct {
	ID      uuid.UUID            `json:"id,omitempty"`
	Name    string               `json:"name"`
	Worker  string               `json:"worker,omitempty"`
	Version string               `json:"version,omitempty"`
	Graph   map[string]*FlowNode `json:"graph"`
	Enabled bool                 `json:"enabled,omitempty"`
}

// Create creates a new workflow
func (c *WorkflowClient) Create(ctx context.Context, req *CreateWorkflowRequest) (*Workflow, error) {
	var workflow Workflow
	if err := c.http.Post(ctx, "/workflow", req, &workflow); err != nil {
		return nil, err
	}
	return &workflow, nil
}

// CreateFromJSON creates a new workflow from a JSON graph
func (c *WorkflowClient) CreateFromJSON(ctx context.Context, name string, graphJSON json.RawMessage) (*Workflow, error) {
	body := map[string]interface{}{
		"name":  name,
		"graph": graphJSON,
	}
	var workflow Workflow
	if err := c.http.Post(ctx, "/workflow", body, &workflow); err != nil {
		return nil, err
	}
	return &workflow, nil
}

// Delete deletes a workflow
func (c *WorkflowClient) Delete(ctx context.Context, workflowID uuid.UUID) error {
	return c.http.Delete(ctx, "/workflow/"+workflowID.String())
}

// SetEnabled enables or disables a workflow
func (c *WorkflowClient) SetEnabled(ctx context.Context, workflowID uuid.UUID, enabled bool) error {
	body := map[string]bool{
		"enabled": enabled,
	}
	return c.http.Patch(ctx, "/workflow/"+workflowID.String()+"/enabled", body, nil)
}

// Enable enables a workflow
func (c *WorkflowClient) Enable(ctx context.Context, workflowID uuid.UUID) error {
	return c.SetEnabled(ctx, workflowID, true)
}

// Disable disables a workflow
func (c *WorkflowClient) Disable(ctx context.Context, workflowID uuid.UUID) error {
	return c.SetEnabled(ctx, workflowID, false)
}

// GetOperations returns available workflow operations
func (c *WorkflowClient) GetOperations(ctx context.Context) (map[string]Operation, error) {
	var operations map[string]Operation
	if err := c.http.Get(ctx, "/workflow/operation", &operations); err != nil {
		return nil, err
	}
	return operations, nil
}

// GetRunning returns currently running workflows
func (c *WorkflowClient) GetRunning(ctx context.Context) (map[string]interface{}, error) {
	var running map[string]interface{}
	if err := c.http.Get(ctx, "/workflow/running", &running); err != nil {
		return nil, err
	}
	return running, nil
}

// Health checks the workflow engine health
func (c *WorkflowClient) Health(ctx context.Context) (map[string]interface{}, error) {
	var health map[string]interface{}
	if err := c.http.Get(ctx, "/workflow/health", &health); err != nil {
		return nil, err
	}
	return health, nil
}

// === Workflow Builder Helpers ===

// NewInputNode creates an input node for the workflow
func NewInputNode(id, sensorID string) *FlowNode {
	return &FlowNode{
		ID:       id,
		Operator: "input",
		Name:     sensorID,
	}
}

// NewOutputNode creates an output node for the workflow
func NewOutputNode(id, sensorID string, inputs []string) *FlowNode {
	return &FlowNode{
		ID:       id,
		Operator: "output",
		Name:     sensorID,
		Inputs:   inputs,
	}
}

// NewConstNode creates a constant value node
func NewConstNode(id, value string) *FlowNode {
	return &FlowNode{
		ID:       id,
		Operator: "const",
		Name:     value,
	}
}

// NewSelectNode creates a select node with conditions
func NewSelectNode(id string, inputs, values, conditions []string) *FlowNode {
	return &FlowNode{
		ID:        id,
		Operator:  "select",
		Inputs:    inputs,
		Values:    values,
		Condition: conditions,
	}
}

// NewAndNode creates an AND logic node
func NewAndNode(id string, inputs []string) *FlowNode {
	return &FlowNode{
		ID:       id,
		Operator: "and",
		Inputs:   inputs,
	}
}

// NewOrNode creates an OR logic node
func NewOrNode(id string, inputs []string) *FlowNode {
	return &FlowNode{
		ID:       id,
		Operator: "or",
		Inputs:   inputs,
	}
}

// NewNotNode creates a NOT logic node
func NewNotNode(id string, input string) *FlowNode {
	return &FlowNode{
		ID:       id,
		Operator: "not",
		Inputs:   []string{input},
	}
}

// NewCompareNode creates a comparison node (lt, le, gt, ge, eq, ne)
func NewCompareNode(id, operator string, inputs []string) *FlowNode {
	return &FlowNode{
		ID:       id,
		Operator: operator,
		Inputs:   inputs,
	}
}

// NewTrafficControlNode creates a traffic control node
func NewTrafficControlNode(id, layoutID string, density int, speed string, enabled bool, inputs []string) *FlowNode {
	return &FlowNode{
		ID:       id,
		Operator: "traffic_control",
		Inputs:   inputs,
		LayoutID: layoutID,
		Density:  density,
		Speed:    speed,
		Enabled:  &enabled,
	}
}

// NewCameraControlNode creates a camera control node
func NewCameraControlNode(id, layoutID, pathID, action string, inputs []string) *FlowNode {
	return &FlowNode{
		ID:       id,
		Operator: "camera_control",
		Inputs:   inputs,
		LayoutID: layoutID,
		PathID:   pathID,
		Action:   action,
	}
}

// NewGroupControlNode creates a group control node
func NewGroupControlNode(id, groupID, outputField, outputValue string, inputs []string) *FlowNode {
	return &FlowNode{
		ID:          id,
		Operator:    "group_control",
		Inputs:      inputs,
		GroupID:     groupID,
		OutputField: outputField,
		OutputValue: outputValue,
	}
}

// WithLatch adds latch configuration to a node
func (n *FlowNode) WithLatch(latchValue string) *FlowNode {
	n.Latch = true
	n.LatchValue = latchValue
	return n
}
