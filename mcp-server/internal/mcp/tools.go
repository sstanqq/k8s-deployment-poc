package mcp

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerTools(srv *mcp.Server, store RequestStore, host HostProvider) {
	// TODO: move spec to .yaml file
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "nodeinfo",
		Description: "return host and internal info",
	}, wrapTool(store, func(ctx context.Context, req *mcp.CallToolRequest, input NodeInfoInput) (
		*mcp.CallToolResult,
		NodeInfoOutput,
		error,
	) {
		return getNodeInfo(ctx, req, input, host)
	}))
}

func wrapTool[I any, O any](
	store RequestStore,
	fn func(ctx context.Context, req *mcp.CallToolRequest, input I) (*mcp.CallToolResult, O, error),
) func(ctx context.Context, req *mcp.CallToolRequest, input I) (*mcp.CallToolResult, O, error) {

	return func(ctx context.Context, req *mcp.CallToolRequest, input I) (*mcp.CallToolResult, O, error) {
		result, output, err := fn(ctx, req, input)

		timestamp := time.Now()

		toolName := req.Params.Name
		sessionID := req.Session.ID()

		if err != nil {
			_ = store.SaveError(toolName, sessionID, fmt.Errorf("[%s] %w", timestamp.Format(time.RFC3339), err))
		} else {
			_ = store.SaveResult(toolName, sessionID, output)
		}

		log.Printf("\"%s\" tool called at %s, session: %s, error: %v", toolName, timestamp.Format(time.RFC3339), sessionID, err)
		return result, output, err
	}
}

type NodeInfo struct {
	Hostname   string `json:"hostname"`
	InternalIP string `json:"internal_ip"`
}

type NodeInfoInput struct {
}

type NodeInfoOutput struct {
	NodeInfo *NodeInfo `json:"node_info"`
}

func getNodeInfo(ctx context.Context, req *mcp.CallToolRequest, input NodeInfoInput, host HostProvider) (*mcp.CallToolResult, NodeInfoOutput, error) {
	hostname, ip, err := host.GetNodeInfo()
	if err != nil {
		return nil, NodeInfoOutput{}, fmt.Errorf("failed to get node info: %w", err)
	}

	output := NodeInfoOutput{&NodeInfo{
		Hostname:   hostname,
		InternalIP: ip,
	}}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf(
					"Hostname: %s, Internal IP: %s", hostname, ip,
				),
			},
		},
	}

	return result, output, nil
}
