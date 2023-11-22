package zitimesh

const AgentServiceConfgSchema = `
{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "type": "object",
    "properties": {
        "port": {
            "type": "integer"
        },
        "network": {
            "type": "string",
            "enum": ["tcp", "udp"]
        }
    },
    "required": ["port", "network"]
}
`

type NetworkType string

const (
	NetworkTypeTCP NetworkType = "tcp"
	NetworkTypeUDP NetworkType = "udp"
)

type AgentService struct {
	// Port is the port number to expose
	Port int `json:"port"`

	// Network is the type of network to use ("tcp" or "udp")
	Network NetworkType `json:"network"`
}
