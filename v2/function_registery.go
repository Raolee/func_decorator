package v2

import (
	"errors"
	"fmt"
	"reflect"
)

type FunctionNode struct {
	ID       string
	Function AnyFunction
	Next     set[string]
}

func NewFunctionNode(id string, f AnyFunction) *FunctionNode {
	return &FunctionNode{
		ID:       id,
		Function: f,
		Next:     newSet[string](),
	}
}

type FunctionRegistry interface {
	RegisterFunction(id string, f AnyFunction)
	GetFunctionNode(id string) (*FunctionNode, bool)
	DeregisterFunction(id string)
	ConnectFunctionNode(fromId, toId string, adapters ...AnyFunction) error
}

type functionRegistry struct {
	nodes map[string]*FunctionNode
}

func NewFunctionRegistry() FunctionRegistry {
	return &functionRegistry{
		nodes: make(map[string]*FunctionNode),
	}
}

func (r *functionRegistry) RegisterFunction(id string, f AnyFunction) {
	r.nodes[id] = NewFunctionNode(id, f)
}

func (r *functionRegistry) GetFunctionNode(id string) (*FunctionNode, bool) {
	node, ok := r.nodes[id]
	return node, ok
}

func (r *functionRegistry) DeregisterFunction(id string) {
	delete(r.nodes, id)
}

func (r *functionRegistry) ConnectFunctionNode(fromId, toId string, adapters ...AnyFunction) error {
	fromNode, ok := r.nodes[fromId]
	if !ok {
		return errors.New(fmt.Sprintf("'fromNode.id=%s' not exists", fromId))
	}
	toNode, ok := r.nodes[toId]
	if !ok {
		return errors.New(fmt.Sprintf("'toNode.id=%s' not exists", toId))
	}

	// Node 연결 규칙 검사
	if err := r.validateConnectionFunction(fromNode, toNode, adapters...); err != nil {
		return err
	}

	// TODO : adapters 를 넣어줘야 함
	fromNode.Next.Add(toId) // Node 에 직접 다음 것을 넣어줌
	return nil
}

// FunctionNode 가 연결될 수 있는 규칙을 검사함
func (r *functionRegistry) validateConnectionFunction(fromNode, toNode *FunctionNode, adapters ...AnyFunction) error {

	// 두 노드가 같은 노드 일 때
	if fromNode.ID == toNode.ID {
		return fmt.Errorf("fromNode and toNode have the same ID (%s)", fromNode.ID)
	}
	// 이미 연결된 관계
	if fromNode.Next.Exists(toNode.ID) {
		return fmt.Errorf("A connection to toNode is already defined in fromNode")
	}
	// toNode 가 fromNode 를 가리키는 상황
	if toNode.Next.Exists(fromNode.ID) {
		return fmt.Errorf("the connection between fromNode and toNode is recursive. ")
	}
	// funcNode 응답 타입과 toNode 요청 타입이 다를 때 => adapters 가 타입을 맞춰줘야 함
	if !fromNode.Function.GetResponseType().AssignableTo(toNode.Function.GetRequestType()) {

		// adapters 가 비었을 때
		if len(adapters) == 0 {
			return fmt.Errorf("when the response type of fromNode and the request type of toNode are different, adapters can't nil")
		}

		// fromNode 응답 타입과 adapter[0] 의 요청 타입이 일치하는가?
		checkList := make([]reflect.Type, 0)

		// fromNode 의 응답타입 추가
		checkList = append(checkList, fromNode.Function.GetResponseType())
		// adapter 의 요청타입 및 응답타입 추가
		for _, adapter := range adapters {
			checkList = append(checkList, adapter.GetRequestType(), adapter.GetResponseType())
		}
		// toNode 의 요청 타입 추가
		checkList = append(checkList, toNode.Function.GetRequestType())

		// 응답, 요청 타입을 짝으로 전부 체크
		for i := 0; i < len(checkList); i += 2 {
			if !checkList[i].AssignableTo(checkList[i+1]) {
				// 하나라도 타입이 맞지 않으면 에러!
				return fmt.Errorf("checkList[%d] (%s) not equal checkList[%d] (%s)",
					i,
					checkList[i],
					i+1,
					checkList[i+1],
				)
			}
		}
	}

	return nil
}
