package v2

import (
	"context"
	"fmt"
	"sync"
)

type ExecuteResultMap struct {
	results sync.Map
}

func NewExecuteResultMap() *ExecuteResultMap {
	return &ExecuteResultMap{results: sync.Map{}}
}

func (e *ExecuteResultMap) addResult(key any, result *ExecuteResult) {
	e.results.Store(key, result)
}

func (e *ExecuteResultMap) Slice() []*ExecuteResult {
	rts := make([]*ExecuteResult, 0)
	e.results.Range(func(key, value any) bool {
		rts = append(rts, value.(*ExecuteResult))
		return true // 끝까지 돌아라!
	})
	return rts
}

type ExecuteResult struct {
	NodeFlow string
	NodeID   string
	Req      any
	Res      any
	Err      error
}

func NewExecuteResult(nodeFlow, nodeID string, req, res any, err error) *ExecuteResult {
	return &ExecuteResult{
		NodeFlow: nodeFlow,
		NodeID:   nodeID,
		Req:      req,
		Res:      res,
		Err:      err,
	}
}

type FunctionChainExecutor struct {
	registry         FunctionRegistry
	adapter          FunctionAdapter
	executeResultMap *ExecuteResultMap
	current          int32
}

func NewFunctionChainExecutor(registry FunctionRegistry, adapter FunctionAdapter) *FunctionChainExecutor {
	return &FunctionChainExecutor{
		registry:         registry,
		adapter:          adapter,
		executeResultMap: NewExecuteResultMap(),
		current:          0,
	}
}

func (e *FunctionChainExecutor) Execute(startID string, ctx context.Context, req any) (*ExecuteResultMap, error) {
	startNode, ok := e.registry.GetFunctionNode(startID)
	if !ok {
		return e.executeResultMap, fmt.Errorf("function node not found: %s", startID)
	}

	// Adapter 를 사용하여 인자 변환
	// TODO : StartID => Next 를 따져서 변홚 해야함
	newReq, err := e.adapter.Adapt(startID, ctx, req, func(ctx context.Context, after any, err error) {
		e.current++
		e.executeResultMap.addResult(e.current, NewExecuteResult(GetNodeFlowInContext(ctx), startID, req, after, err))
	})
	if err != nil {
		return e.executeResultMap, err
	}

	res, err := startNode.Function(ctx, newReq)
	e.current++
	e.executeResultMap.addResult(e.current, NewExecuteResult(GetNodeFlowInContext(ctx), startID, newReq, res, err))

	if err != nil {
		return e.executeResultMap, err
	}

	for _, nextID := range startNode.Next.GetElems() {
		ctx = SetNodeFlowInContext(ctx, startID)
		_, err := e.Execute(nextID, ctx, res)
		if err != nil {
			return e.executeResultMap, err
		}
	}

	return e.executeResultMap, nil
}
