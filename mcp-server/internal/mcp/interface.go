package mcp

type RequestStore interface {
	SaveResult(toolName, sessionID string, result any) error
	SaveError(toolName, sessionID string, err error) error
}

type HostProvider interface {
	GetNodeInfo() (hostname, internalIP string, err error)
}
