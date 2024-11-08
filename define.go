package approval

import "github.com/jackc/pgtype"

type NodeType string

type EdgeType string

type Process struct {
	Nodes []Node `json:"nodes"` // 节点列表
	Edges []Edge `json:"edges"` // 连接线列表
}

// Node 表示节点，包含节点的基本信息和审批配置等
type Node struct {
	Name       string       `json:"name"`        // 节点名称
	ID         string       `json:"id"`          // 节点ID
	NodeType   NodeType     `json:"nodeType"`    // 节点类型
	NodeExtras pgtype.JSONB `json:"node_extras"` // 节点扩展信息
}

func (n Node) GetID() string {
	return n.ID
}

func (n Node) GetType() NodeType {
	return n.NodeType
}

func (n Node) GetExtras() []byte {
	return n.NodeExtras.Bytes
}

// Edge 表示连接线，包含连接节点的信息
type Edge struct {
	ID       string `json:"id"`       // 连接线的唯一标识
	DstID    string `json:"dstId"`    // 目标节点的ID
	SrcID    string `json:"srcId"`    // 源节点的ID
	Priority int    `json:"priority"` // 优先级

	// 扩展信息
	EdgeType   EdgeType     `json:"line_type"`
	EdgeExtras pgtype.JSONB `json:"line_extras"`
}

func (e Edge) GetID() string {
	return e.ID
}

func (e Edge) GetSrcID() string {
	return e.SrcID
}

func (e Edge) GetDstID() string {
	return e.DstID
}

func (e Edge) GetType() EdgeType {
	return e.EdgeType
}

func (e Edge) GetPriority() int {
	return e.Priority
}

func (e Edge) GetExtras() []byte {
	return e.EdgeExtras.Bytes
}
