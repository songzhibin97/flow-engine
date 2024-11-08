package approval

import (
	"context"
	"fmt"
	"sort"
)

// NodeInterface 定义节点的基本行为
type NodeInterface interface {
	GetID() string     // 获取节点ID
	GetType() NodeType // 获取节点类型
	GetExtras() []byte // 获取节点扩展信息
}

// EdgeInterface 定义连接线的基本行为
type EdgeInterface interface {
	GetID() string     // 获取边的ID
	GetSrcID() string  // 获取源节点ID
	GetDstID() string  // 获取目标节点ID
	GetType() EdgeType // 获取连接线类型
	GetPriority() int  // 获取边的优先级
	GetExtras() []byte // 获取连接线的扩展信息
}

// ProcessInterface 定义流程操作的接口
type ProcessInterface interface {
	GetNodeByID(nodeID string) (NodeInterface, bool)             // 根据节点ID获取节点
	GetOutgoingEdges(nodeID string) []EdgeInterface              // 获取指定节点的出边
	GetIncomingEdges(nodeID string) []EdgeInterface              // 获取指定节点的入边
	GetNextNodes(nodeID string) ([]NodeInterface, error)         // 获取给定节点的下一个节点列表
	GetNextNodeByEdge(edge EdgeInterface) (NodeInterface, error) // 根据边获取目标节点
}

// ProcessRegistryInterface 用于管理节点和边处理器
type ProcessRegistryInterface interface {
	RegisterNodeHandler(nodeType NodeType, handler NodeHandlerFunc)
	RegisterNodeAfterHandler(nodeType NodeType, handler NodeHandlerAfterFunc)
	RegisterNodeBeforeHandler(nodeType NodeType, handler NodeHandlerBeforeFunc)
	RegisterEdgeHandler(edgeType EdgeType, handler EdgeHandlerFunc)
	GetNodeHandler(nodeType NodeType) (NodeHandlerFunc, bool)
	GetNodeAfterHandler(nodeType NodeType) (NodeHandlerAfterFunc, bool)
	GetNodeBeforeHandler(nodeType NodeType) (NodeHandlerBeforeFunc, bool)
	GetEdgeHandler(edgeType EdgeType) (EdgeHandlerFunc, bool)
}

// CirculateHandler 用于处理点与点直接流转的逻辑抽象
type CirculateHandler interface {
	CirculateHandler(ctx context.Context, edge EdgeInterface) (bool, error)
}

// NodeHandlerFunc 定义节点处理函数
type NodeHandlerFunc func(ctx context.Context, node NodeInterface) error

// NodeHandlerBeforeFunc 定义节点处理前的处理函数
type NodeHandlerBeforeFunc func(ctx context.Context, node NodeInterface) error

// NodeHandlerAfterFunc 定义节点处理后的处理函数
type NodeHandlerAfterFunc func(ctx context.Context, node NodeInterface) error

// EdgeHandlerFunc 定义边处理函数
type EdgeHandlerFunc func(ctx context.Context, edge EdgeInterface) (bool, error)

// Graph 实现 ProcessInterface，用于存储和管理流程中的节点和边
type Graph struct {
	nodes map[string]NodeInterface
	edges map[string]EdgeInterface
}

func NewProcess(process *Process) *Graph {
	p := &Graph{
		nodes: make(map[string]NodeInterface),
		edges: make(map[string]EdgeInterface),
	}
	for _, node := range process.Nodes {
		p.nodes[node.ID] = node
	}
	for _, edge := range process.Edges {
		p.edges[edge.ID] = edge
	}

	// 边按照优先级排一下序
	sort.Slice(process.Edges, func(i, j int) bool {
		return process.Edges[i].Priority < process.Edges[j].Priority
	})

	return p
}

func (p *Graph) GetNodeByID(nodeID string) (NodeInterface, bool) {
	node, exists := p.nodes[nodeID]
	return node, exists
}

func (p *Graph) GetOutgoingEdges(nodeID string) []EdgeInterface {
	var edges []EdgeInterface
	for _, edge := range p.edges {
		if edge.GetSrcID() == nodeID {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (p *Graph) GetIncomingEdges(nodeID string) []EdgeInterface {
	var edges []EdgeInterface
	for _, edge := range p.edges {
		if edge.GetDstID() == nodeID {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (p *Graph) GetNextNodes(nodeID string) ([]NodeInterface, error) {
	var nextNodes []NodeInterface
	for _, edge := range p.GetOutgoingEdges(nodeID) {
		if nextNode, exists := p.GetNodeByID(edge.GetDstID()); exists {
			nextNodes = append(nextNodes, nextNode)
		}
	}
	return nextNodes, nil
}

func (p *Graph) GetNextNodeByEdge(edge EdgeInterface) (NodeInterface, error) {
	if nextNode, exists := p.GetNodeByID(edge.GetDstID()); exists {
		return nextNode, nil
	}
	return nil, fmt.Errorf("未找到边 %s 的目标节点", edge.GetID())
}

// ProcessRegistry 实现 ProcessRegistryInterface，用于注册和获取节点和边的处理函数
type ProcessRegistry struct {
	nodeHandlers       map[NodeType]NodeHandlerFunc
	nodeAfterHandlers  map[NodeType]NodeHandlerAfterFunc
	nodeBeforeHandlers map[NodeType]NodeHandlerBeforeFunc
	edgeHandlers       map[EdgeType]EdgeHandlerFunc
}

func NewProcessRegistry() *ProcessRegistry {
	return &ProcessRegistry{}
}

func (r *ProcessRegistry) RegisterNodeHandler(nodeType NodeType, handler NodeHandlerFunc) {
	if r.nodeHandlers == nil {
		r.nodeHandlers = make(map[NodeType]NodeHandlerFunc)
	}
	r.nodeHandlers[nodeType] = handler
}

func (r *ProcessRegistry) RegisterNodeAfterHandler(nodeType NodeType, handler NodeHandlerAfterFunc) {
	if r.nodeAfterHandlers == nil {
		r.nodeAfterHandlers = make(map[NodeType]NodeHandlerAfterFunc)
	}
	r.nodeAfterHandlers[nodeType] = handler
}

func (r *ProcessRegistry) RegisterNodeBeforeHandler(nodeType NodeType, handler NodeHandlerBeforeFunc) {
	if r.nodeBeforeHandlers == nil {
		r.nodeBeforeHandlers = make(map[NodeType]NodeHandlerBeforeFunc)
	}
	r.nodeBeforeHandlers[nodeType] = handler
}

func (r *ProcessRegistry) RegisterEdgeHandler(edgeType EdgeType, handler EdgeHandlerFunc) {
	if r.edgeHandlers == nil {
		r.edgeHandlers = make(map[EdgeType]EdgeHandlerFunc)
	}
	r.edgeHandlers[edgeType] = handler
}

func (r *ProcessRegistry) GetNodeHandler(nodeType NodeType) (NodeHandlerFunc, bool) {
	handler, exists := r.nodeHandlers[nodeType]
	return handler, exists
}

func (r *ProcessRegistry) GetNodeAfterHandler(nodeType NodeType) (NodeHandlerAfterFunc, bool) {
	handler, exists := r.nodeAfterHandlers[nodeType]
	return handler, exists
}

func (r *ProcessRegistry) GetNodeBeforeHandler(nodeType NodeType) (NodeHandlerBeforeFunc, bool) {
	handler, exists := r.nodeBeforeHandlers[nodeType]
	return handler, exists
}

func (r *ProcessRegistry) GetEdgeHandler(edgeType EdgeType) (EdgeHandlerFunc, bool) {
	handler, exists := r.edgeHandlers[edgeType]
	return handler, exists
}

// ProcessHandler 用于执行节点和边的处理，控制流程流转
type ProcessHandler struct {
	process   ProcessInterface
	registry  ProcessRegistryInterface
	circulate CirculateHandler
}

func (ph *ProcessHandler) ProcessNode(ctx context.Context, nodeID string) error {
	node, exists := ph.process.GetNodeByID(nodeID)
	if !exists {
		return fmt.Errorf("节点 %s 不存在", nodeID)
	}

	// 获取节点处理函数并执行
	nodeHandler, exists := ph.registry.GetNodeHandler(node.GetType())
	if !exists {
		return fmt.Errorf("未找到节点类型 %s 的处理器", node.GetType())
	}
	if err := nodeHandler(ctx, node); err != nil {
		return err
	}

	// 处理节点的出边
	for _, edge := range ph.process.GetOutgoingEdges(nodeID) {
		ok, err := ph.ProcessEdge(ctx, edge)
		if err != nil {
			return err
		}
		if ok {
			// 执行节点处理后的处理函数
			nodeAfterHandler, exists := ph.registry.GetNodeAfterHandler(node.GetType())
			if exists {
				if err := nodeAfterHandler(ctx, node); err != nil {
					return err
				}
			}

			// 状态变更
			if ph.circulate != nil {
				ok, err = ph.circulate.CirculateHandler(ctx, edge)
				if err != nil {
					return err
				}
			}

			// 根据边的目标节点ID继续处理
			nextNode, err := ph.process.GetNextNodeByEdge(edge)
			if err != nil {
				return err
			}
			nodeBeforeHandler, exists := ph.registry.GetNodeBeforeHandler(nextNode.GetType())
			if exists {
				if err := nodeBeforeHandler(ctx, nextNode); err != nil {
					return err
				}
			}

			break
		}
	}
	return nil
}

func (ph *ProcessHandler) ProcessEdge(ctx context.Context, edge EdgeInterface) (bool, error) {
	// 获取边处理函数并执行
	edgeHandler, exists := ph.registry.GetEdgeHandler(edge.GetType())
	if !exists {
		return false, fmt.Errorf("未找到边类型 %s 的处理器", edge.GetType())
	}
	return edgeHandler(ctx, edge)
}
