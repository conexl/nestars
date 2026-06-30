package telegram

import (
	"sync"

	"nestars-bot/internal/domain"
)

type FSM struct {
	mu    sync.RWMutex
	state map[int64]domain.FlowState
	stack map[int64][]domain.FlowState
	data  map[int64]map[string]any
}

func NewFSM() *FSM {
	return &FSM{
		state: make(map[int64]domain.FlowState),
		stack: make(map[int64][]domain.FlowState),
		data:  make(map[int64]map[string]any),
	}
}

func (f *FSM) Current(chatID int64) domain.FlowState {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if s, ok := f.state[chatID]; ok {
		return s
	}
	return domain.StateMain
}

func (f *FSM) Push(chatID int64, s domain.FlowState) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.stack[chatID] = append(f.stack[chatID], s)
	f.state[chatID] = s
}

func (f *FSM) Pop(chatID int64) domain.FlowState {
	f.mu.Lock()
	defer f.mu.Unlock()
	st := f.stack[chatID]
	if len(st) == 0 {
		f.state[chatID] = domain.StateMain
		return domain.StateMain
	}
	st = st[:len(st)-1]
	f.stack[chatID] = st
	if len(st) == 0 {
		f.state[chatID] = domain.StateMain
		return domain.StateMain
	}
	f.state[chatID] = st[len(st)-1]
	return f.state[chatID]
}

func (f *FSM) Replace(chatID int64, s domain.FlowState) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.stack[chatID]) > 0 {
		f.stack[chatID][len(f.stack[chatID])-1] = s
	} else {
		f.stack[chatID] = []domain.FlowState{s}
	}
	f.state[chatID] = s
}

func (f *FSM) Reset(chatID int64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.state, chatID)
	delete(f.stack, chatID)
	delete(f.data, chatID)
}

func (f *FSM) SetData(chatID int64, key string, value any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.data[chatID] == nil {
		f.data[chatID] = make(map[string]any)
	}
	f.data[chatID][key] = value
}

func (f *FSM) GetData(chatID int64, key string) (any, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if d, ok := f.data[chatID]; ok {
		v, ok := d[key]
		return v, ok
	}
	return nil, false
}

func (f *FSM) DeleteData(chatID int64, key string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if d, ok := f.data[chatID]; ok {
		delete(d, key)
	}
}
