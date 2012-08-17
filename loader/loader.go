
package loader

import (
  "strings"
  "fmt"
  "errors"
  "io/ioutil"
  "encoding/json"
  "reflect"
  "github.com/rwcarlsen/goclus/sim"
  "github.com/rwcarlsen/goclus/msg"
)

var agentLib map[string]reflect.Type

func Register(a interface{}) {
  if agentLib == nil {
    agentLib = map[string]reflect.Type{}
  }
  v := reflect.Indirect(reflect.ValueOf(a))
  t := v.Type()
  name := t.PkgPath() + "." + t.Name()
  agentLib[name] = t
}

type ProtoInfo struct {
  ImportPath string
  Id string
  Config string
}

type AgentInfo struct {
  Id string
  ProtoId string
  ParentId string
}

type SimInput struct {
  Prototypes []*ProtoInfo
  Agents []*AgentInfo
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
    return nil, prettyParseError(string(data), err)
  }

  // create prototypes
  protos := map[string]interface{}{}
  for _, info := range input.Prototypes {
    p := reflect.New(agentLib[info.ImportPath]).Interface()

    // configure the agent according to input
    data := []byte(info.Config)
    err := json.Unmarshal(data, p)
    if err != nil {
      return nil, prettyParseError(info.Config, err)
    }

    fmt.Println("method count: ", reflect.TypeOf(p).NumMethod())
    protos[info.Id] = reflect.Indirect(reflect.ValueOf(p)).Interface()
  }

  // create agents from prototypes
  agents := []interface{}{}
  agentMap := map[string]interface{}{}
  for _, info := range input.Agents {
    a := protos[info.ProtoId]
    agents = append(agents, &a)
  }

  // set parents
  for i, info := range input.Agents {
    tp := reflect.TypeOf(agents[i])
    fmt.Println("all methods of:", tp.Method(0))
    for i := 0; i < tp.NumMethod(); i++ {
      fmt.Println(tp.Method(i))
    }

    if setParent, ok := tp.MethodByName("SetParent"); ok {
      if par, ok := agentMap[info.ParentId]; ok {
        setParent.Func.Call([]reflect.Value{reflect.ValueOf(par)})
      }
    } else {
      return nil, errors.New("loader: non-communicator cannot have parent")
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

func prettyParseError(js string, err error) error {
  syntax, ok := err.(*json.SyntaxError)
  if !ok {
    return err
  }

  //start, end := strings.LastIndex(js[:syntax.Offset], "\n")+1, len(js)
  start, _ := strings.LastIndex(js[:syntax.Offset], "\n")+1, len(js)
  //if idx := strings.Index(js[start:], "\n"); idx >= 0 {
  //  end = start + idx
  //}
  line, pos := strings.Count(js[:start], "\n"), int(syntax.Offset) - start - 1
  msg := fmt.Sprint("loader: ", err, " at line ", line + 1, " pos ", pos)
  return errors.New(msg)
  //fmt.Printf("%s\n%s^", js[start:end], strings.Repeat(" ", pos))
}
