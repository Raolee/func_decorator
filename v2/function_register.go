package v2

import (
	"errors"
	"fmt"
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
	ConnectFunctionNode(fromId, toId string) error
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

func (r *functionRegistry) ConnectFunctionNode(fromId, toId string) error {
	fromNode, ok := r.nodes[fromId]
	if !ok {
		return errors.New(fmt.Sprintf("'fromNode.id=%s' don't exists", fromId))
	}
	_, ok = r.nodes[toId]
	if !ok {
		return errors.New(fmt.Sprintf("'toNode.id=%s' don't exists", toId))
	}

	fromNode.Next.Add(toId)
	return nil
}
