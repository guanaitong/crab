package ms

import "math/rand"

type LoadBalance interface {
	DoSelect(list []*ServiceInstance) *ServiceInstance
}

type RandomLoadBalance struct {
}

func (lb *RandomLoadBalance) DoSelect(list []*ServiceInstance) *ServiceInstance {
	length := len(list)
	if length == 0 {
		return nil
	}
	if length == 1 {
		return list[0]
	}
	return list[rand.Intn(length)]
}
