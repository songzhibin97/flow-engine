package approval

import (
	"context"
	"testing"

	"github.com/jackc/pgtype"
)

// 节点类型
const (
	nodeTypeLeaveRequest  NodeType = "leave_request"
	nodeTypeLeaveApproval NodeType = "leave_approval"
	nodeTypeLeaveComplete NodeType = "leave_complete"
)

// 边类型
const (
	edgeTypeApproval EdgeType = "approval"
	edgeTypeComplete EdgeType = "complete"
)

// 请假申请节点处理
func handleLeaveRequest(ctx context.Context, node NodeInterface) error {
	// 获取请假信息,保存到数据库等
	return nil
}

// 请假审批节点处理
func handleLeaveApproval(ctx context.Context, node NodeInterface) error {
	// 获取审批信息,更新请假状态
	return nil
}

// 请假完成节点处理
func handleLeaveComplete(ctx context.Context, node NodeInterface) error {
	// 更新请假状态为已完成
	return nil
}

// 请假审批边处理
func handleApprovalEdge(ctx context.Context, edge EdgeInterface) (bool, error) {
	// 检查审批是否通过,返回true表示通过
	return true, nil
}

// 请假完成边处理
func handleCompleteEdge(ctx context.Context, edge EdgeInterface) (bool, error) {
	// 检查是否可以完成请假,返回true表示可以
	return true, nil
}

var emptyExtras = pgtype.JSONB{
	Status: pgtype.Null,
}

func createLeaveProcess() *Process {
	return &Process{
		Nodes: []Node{
			{
				ID:         "leave_request",
				Name:       "请假申请",
				NodeType:   nodeTypeLeaveRequest,
				NodeExtras: emptyExtras,
			},
			{
				ID:         "leave_approval",
				Name:       "请假审批",
				NodeType:   nodeTypeLeaveApproval,
				NodeExtras: emptyExtras,
			},
			{
				ID:         "leave_complete",
				Name:       "请假完成",
				NodeType:   nodeTypeLeaveComplete,
				NodeExtras: emptyExtras,
			},
		},
		Edges: []Edge{
			{
				ID:         "approval_edge",
				DstID:      "leave_approval",
				SrcID:      "leave_request",
				Priority:   1,
				EdgeType:   edgeTypeApproval,
				EdgeExtras: emptyExtras,
			},
			{
				ID:         "complete_edge",
				DstID:      "leave_complete",
				SrcID:      "leave_approval",
				Priority:   1,
				EdgeType:   edgeTypeComplete,
				EdgeExtras: emptyExtras,
			},
		},
	}
}

func TestNodeHandlers(t *testing.T) {
	// 测试请假申请节点处理
	node := &Node{
		ID:       "leave_request",
		NodeType: nodeTypeLeaveRequest,
	}
	err := handleLeaveRequest(context.Background(), node)
	if err != nil {
		t.Errorf("handleLeaveRequest 出错: %v", err)
	}

	// 测试请假审批节点处理
	node.ID = "leave_approval"
	node.NodeType = nodeTypeLeaveApproval
	err = handleLeaveApproval(context.Background(), node)
	if err != nil {
		t.Errorf("handleLeaveApproval 出错: %v", err)
	}

	// 测试请假完成节点处理
	node.ID = "leave_complete"
	node.NodeType = nodeTypeLeaveComplete
	err = handleLeaveComplete(context.Background(), node)
	if err != nil {
		t.Errorf("handleLeaveComplete 出错: %v", err)
	}
}

func TestEdgeHandlers(t *testing.T) {
	// 测试请假审批边处理
	edge := &Edge{
		ID:       "approval_edge",
		DstID:    "leave_approval",
		SrcID:    "leave_request",
		EdgeType: edgeTypeApproval,
	}
	ok, err := handleApprovalEdge(context.Background(), edge)
	if err != nil {
		t.Errorf("handleApprovalEdge 出错: %v", err)
	}
	if !ok {
		t.Errorf("handleApprovalEdge 返回 false")
	}

	// 测试请假完成边处理
	edge.ID = "complete_edge"
	edge.DstID = "leave_complete"
	edge.SrcID = "leave_approval"
	edge.EdgeType = edgeTypeComplete
	ok, err = handleCompleteEdge(context.Background(), edge)
	if err != nil {
		t.Errorf("handleCompleteEdge 出错: %v", err)
	}
	if !ok {
		t.Errorf("handleCompleteEdge 返回 false")
	}
}

func TestProcessHandler(t *testing.T) {
	// 创建ProcessRegistry并注册处理函数
	registry := &ProcessRegistry{}
	registry.RegisterNodeHandler(nodeTypeLeaveRequest, handleLeaveRequest)
	registry.RegisterNodeHandler(nodeTypeLeaveApproval, handleLeaveApproval)
	registry.RegisterNodeHandler(nodeTypeLeaveComplete, handleLeaveComplete)
	registry.RegisterEdgeHandler(edgeTypeApproval, handleApprovalEdge)
	registry.RegisterEdgeHandler(edgeTypeComplete, handleCompleteEdge)

	// 创建请假流程
	leaveProcess := createLeaveProcess()

	// 创建ProcessHandler
	handler := &ProcessHandler{
		process:  NewProcess(leaveProcess),
		registry: registry,
	}

	// 测试ProcessNode方法
	err := handler.ProcessNode(context.Background(), "leave_request")
	if err != nil {
		t.Errorf("ProcessNode 出错: %v", err)
	}

	// 测试ProcessEdge方法
	edge := leaveProcess.Edges[0]
	ok, err := handler.ProcessEdge(context.Background(), edge)
	if err != nil {
		t.Errorf("ProcessEdge 出错: %v", err)
	}
	if !ok {
		t.Errorf("ProcessEdge 返回 false")
	}
}

func TestIntegrationLeaveProcess(t *testing.T) {
	// 创建ProcessRegistry并注册处理函数
	registry := &ProcessRegistry{}
	registry.RegisterNodeHandler(nodeTypeLeaveRequest, handleLeaveRequest)
	registry.RegisterNodeHandler(nodeTypeLeaveApproval, handleLeaveApproval)
	registry.RegisterNodeHandler(nodeTypeLeaveComplete, handleLeaveComplete)
	registry.RegisterEdgeHandler(edgeTypeApproval, handleApprovalEdge)
	registry.RegisterEdgeHandler(edgeTypeComplete, handleCompleteEdge)

	// 创建请假流程
	leaveProcess := createLeaveProcess()

	// 创建ProcessHandler
	handler := &ProcessHandler{
		process:  NewProcess(leaveProcess),
		registry: registry,
	}

	// 执行请假申请节点
	err := handler.ProcessNode(context.Background(), "leave_request")
	if err != nil {
		t.Errorf("执行请假申请节点时出错: %v", err)
	}

	// 执行请假审批边
	edge := leaveProcess.Edges[0]
	ok, err := handler.ProcessEdge(context.Background(), edge)
	if err != nil {
		t.Errorf("执行请假审批边时出错: %v", err)
	}
	if !ok {
		t.Errorf("请假审批边处理返回 false")
	}

	// 执行请假审批节点
	nextNode, err := handler.process.GetNextNodeByEdge(edge)
	if err != nil {
		t.Errorf("获取下一个节点时出错: %v", err)
	}
	err = handler.ProcessNode(context.Background(), nextNode.GetID())
	if err != nil {
		t.Errorf("执行请假审批节点时出错: %v", err)
	}

	// 执行请假完成边
	edge = leaveProcess.Edges[1]
	ok, err = handler.ProcessEdge(context.Background(), edge)
	if err != nil {
		t.Errorf("执行请假完成边时出错: %v", err)
	}
	if !ok {
		t.Errorf("请假完成边处理返回 false")
	}

	// 执行请假完成节点
	nextNode, err = handler.process.GetNextNodeByEdge(edge)
	if err != nil {
		t.Errorf("获取下一个节点时出错: %v", err)
	}
	err = handler.ProcessNode(context.Background(), nextNode.GetID())
	if err != nil {
		t.Errorf("执行请假完成节点时出错: %v", err)
	}
}
