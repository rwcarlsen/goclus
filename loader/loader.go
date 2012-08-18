
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
  Config map[string]interface{}
}

type AgentInfo struct {
  Id string
  ProtoId string
  ParentId string
  IndexId string
}

type SimInput struct {
  Prototypes map[string]*ProtoInfo
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

  // load input file
  input := &SimInput{}
  err = json.Unmarshal(data, &input)
  if err != nil {
    return nil, prettyParseError(string(data), err)
  }
  eng := input.Engine
  
  // create prototypes
  protos := map[string]interface{}{}
  imports := map[string]string{}
  for protoId, info := range input.Prototypes {
    pv := reflect.New(agentLib[info.ImportPath])
    protos[protoId]= pv.Interface()
    imports[protoId] = info.ImportPath
  }

  // configure prototypes
  for id, p := range protos {
    data, _ := json.Marshal(input.Prototypes[id].Config)
    json.Unmarshal(data, p)
    fmt.Println(p)
  }

  // create agents from prototypes
  agents := []interface{}{}
  agentMap := map[string]interface{}{}
  for _, info := range input.Agents {
    p := protos[info.ProtoId]
    pv := reflect.Indirect(reflect.ValueOf(p))
    pt := pv.Type()
    av := reflect.New(agentLib[imports[info.ProtoId]])
    for i := 0; i < pv.NumField(); i++ {
      name := pt.Field(i).Name
      pfield := pv.FieldByName(name)
      afield := reflect.Indirect(av).FieldByName(name)
      if afield.CanSet() {
        afield.Set(pfield)
      }
    }

    a := av.Interface()
    agents = append(agents, a)
    if info.IndexId != "" {
      err := eng.RegisterComm(info.IndexId, a.(msg.Communicator))
      if err != nil {
        panic(err.Error())
      }
    }
    fmt.Println(agents[len(agents)-1])
  }

  // set parents
  for i, info := range input.Agents {
    tp := reflect.TypeOf(agents[i])
    if setParent, ok := tp.MethodByName("SetParent"); ok {
      if par, ok := agentMap[info.ParentId]; ok {
        setParent.Func.Call([]reflect.Value{reflect.ValueOf(par)})
      }
    } else {
      return nil, errors.New("loader: non-communicator cannot have parent")
    }
  }

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
