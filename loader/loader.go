
package loader

import (
  "io/ioutil"
  "encoding/json"
  "reflect"
  "github.com/rwcarlsen/goclus/sim"
  "github.com/rwcarlsen/goclus/msg"
)

var agentLib map[string]reflect.Type

func Register(a interface{}) {
  v := reflect.Indirect(reflect.ValueOf(a))
  t := v.Type()
  name := t.PkgPath() + "." + t.Name()
  agentLib[name] = t
}

type AgentInfo struct {
  ImportPath string
  Id string
  ParentId string
}

type SimInput struct {
  Agents []AgentInfo
  Engine *sim.Engine
}

func NewAgent(importPath string, parent msg.Communicator) interface{} {
  a := reflect.New(agentLib[importPath]).Interface()
  a.(msg.Communicator).SetParent(parent)
  return a
}

func LoadSim(fname string) (*sim.Engine, error) {
  data, err := ioutil.ReadFile(fname)
  if err != nil {
    return nil, err
  }

  input := &SimInput{}
  err = json.Unmarshal(data, &input)
  if err != nil {
    return nil, err
  }

  // create all agents
  agents := []interface{}{}
  agentmap := map[string]interface{}{}
  for _, info := range input.Agents {
    a := reflect.New(agentLib[info.ImportPath]).Interface()
    agentmap[info.Id] = a
    agents = append(agents, a)
  }

  // set parents
  for i, info := range input.Agents {
    ac := agents[i].(msg.Communicator)
    if par, ok := agentmap[info.ParentId]; ok {
      ac.SetParent(par.(msg.Communicator))
    }
  }

  // get the engine
  eng := input.Engine

  // register for ticks, tocks, and resolves
  for _, a := range agents {
    switch t := a.(type) {
      case sim.Ticker:
        eng.RegisterTick(t)
    }
    switch t := a.(type) {
      case sim.Tocker:
        eng.RegisterTock(t)
    }
    switch t := a.(type) {
      case sim.Resolver:
        eng.RegisterResolve(t)
    }
  }

  return eng, nil
}
