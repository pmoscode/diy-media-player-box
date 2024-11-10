package storage

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
)

type Manager struct {
	mutex     sync.Mutex
	items     map[string]*item
	threshold int
}

func (m *Manager) Init(processNames string) {
	processes := strings.Split(processNames, ",")

	for _, name := range processes {
		m.addItem(name)
	}
}

func (m *Manager) IncrementHeartbeat(name string) {
	m.mutex.Lock()
	_, exists := m.items[name]

	if !exists {
		m.addItem(name)
	}

	itemInstance := m.items[name]
	itemInstance.IncrementHeartbeatCount()
	m.mutex.Unlock()
}

func (m *Manager) IncrementCheck() {
	m.mutex.Lock()
	for _, itemInstance := range m.items {
		itemInstance.IncrementCheckCount()
	}
	m.mutex.Unlock()
}

func (m *Manager) GetExceededThresholdNames() []string {
	exceededThreshold := make([]string, 0)

	m.mutex.Lock()
	for name, itemInstance := range m.items {
		if itemInstance.checkCount-itemInstance.heartbeatCount >= m.threshold {
			exceededThreshold = append(exceededThreshold, name)
		}
	}
	m.mutex.Unlock()

	return exceededThreshold
}

func (m *Manager) ResetAllCounterFor(processName string) {
	m.mutex.Lock()
	itemInstance := m.items[processName]
	itemInstance.ResetHeartbeatCounter()
	itemInstance.ResetCheckCounter()
	m.mutex.Unlock()
}

func (m *Manager) Debug() {
	builder := strings.Builder{}

	for name, itemInstance := range m.items {
		builder.WriteString("######### Debug #########\n")
		builder.WriteString(fmt.Sprintf("### Name: %s\n", name))
		builder.WriteString(fmt.Sprintf("### CntHeart: %d\n", itemInstance.heartbeatCount))
		builder.WriteString(fmt.Sprintf("### CntCheck: %d\n", itemInstance.checkCount))
		builder.WriteString("############################\n\n")
	}

	log.Println(builder.String())
}

func (m *Manager) addItem(name string) error {
	_, exists := m.items[name]

	if exists {
		return errors.New("item with name '" + name + "' already exists!")
	}

	m.items[name] = newItem()

	return nil
}
