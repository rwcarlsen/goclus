
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

type ProtoInfo struct {
  ImportPath string
  Config map[string]interface{}
}

type AgentInfo struct {
  Id string
  ProtoId string
  ParentId string
  IndexId string
}

type Loader struct {
  Prototypes map[string]*ProtoInfo
  Agents []*AgentInfo
  Engine *sim.Engine
  agentLib map[string]reflect.Type
  protos map[string]interface{}
  imports map[string]string
}

func (l *Loader) Register(a interface{}) {
  if l.agentLib == nil {
    l.agentLib = map[string]reflect.Type{}
  }
  v := reflect.Indirect(reflect.ValueOf(a))
  t := v.Type()
  name := t.PkgPath() + "." + t.Name()
  l.agentLib[name] = t
}

func (l *Loader) NewAgent(importPath string, parent msg.Communicator) interface{} {
  a := l.newPrototype(importPath)
  l.registerWithEngine(a)
  a.(msg.Communicator).SetParent(parent)
  return a
}

func (l *Loader) newPrototype(importPath string) interface{} {
  return reflect.New(l.agentLib[importPath]).Interface()
}

func (l *Loader) NewAgentFromProto(protoId string, parent msg.Communicator) interface{} {
  importPath := l.imports[protoId]
  a := l.NewAgent(importPath, parent)
  data, _ := json.Marshal(l.protos[protoId])
  json.Unmarshal(data, a)
  return a
}

func (l *Loader) registerWithEngine(a interface{}) {
  switch t := a.(type) {
    case sim.Ticker:
      l.Engine.RegisterTick(t)
  }
  switch t := a.(type) {
    case sim.Tocker:
      l.Engine.RegisterTock(t)
  }
  switch t := a.(type) {
    case sim.Resolver:
      l.Engine.RegisterResolve(t)
  }
}

func (l *Loader) LoadSim(fname string) error {
  // load input file
  data, err := ioutil.ReadFile(fname)
  if err != nil {
    return err
  }

  err = json.Unmarshal(data, l)
  if err != nil {
    return prettyParseError(string(data), err)
  }

  // create prototypes
  l.protos = map[string]interface{}{}
  l.imports = map[string]string{}
  for protoId, info := range l.Prototypes {
    p := l.newPrototype(info.ImportPath)
    l.protos[protoId]= p
    l.imports[protoId] = info.ImportPath
  }

  // configure prototypes
  for id, p := range l.protos {
    data, _ := json.Marshal(l.Prototypes[id].Config)
    json.Unmarshal(data, p)
    fmt.Println("prototype: ", p)
  }

  // create agents from prototypes
  agents := []interface{}{}
  agentMap := map[string]interface{}{}
  for _, info := range l.Agents {
    a := l.NewAgentFromProto(info.ProtoId, nil)
    agentMap[info.Id] = a
    agents = append(agents, a)
    if info.IndexId != "" {
      err := l.Engine.RegisterComm(info.IndexId, a.(msg.Communicator))
      if err != nil {
        panic("loader: " + err.Error())
      }
    }
    fmt.Println(agents[len(agents)-1])
  }

  for i, info := range l.Agents {
    // set parents
    if a, ok := agents[i].(msg.Communicator); ok {
      if par, ok := agentMap[info.ParentId]; ok {
        a.SetParent(par.(msg.Communicator))
      }
    } else {
      return errors.New("loader: non-communicator cannot have parent")
    }

    // set Id if can
    if a, ok := agents[i].(sim.Agent); ok {
      a.SetId(info.Id)
    }
  }
  return nil
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
