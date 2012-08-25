package sim

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rwcarlsen/goclus/msg"
	"io/ioutil"
	"reflect"
	"strings"
)

type ProtoInfo struct {
	ImportPath string
	Config     map[string]interface{}
}

type AgentInfo struct {
	Id        string
	ProtoId   string
	ParentId  string
	IsService bool
}

type Loader struct {
	Prototypes map[string]*ProtoInfo
	Agents     []*AgentInfo
	Engine     *Engine
	agentLib   map[string]reflect.Type
	protos     map[string]interface{}
	imports    map[string]string
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
	if parent != nil {
		a.(msg.Communicator).SetParent(parent)
	}
	return a
}

func (l *Loader) newPrototype(importPath string) interface{} {
	if tp, ok := l.agentLib[importPath]; ok {
		return reflect.New(tp).Interface()
	}
	panic("loader: No registered agent for import path '" + importPath + "'")
}

func (l *Loader) NewAgentFromProto(protoId string, parent msg.Communicator) interface{} {
	importPath := l.imports[protoId]
	a := l.NewAgent(importPath, parent)
	data, _ := json.Marshal(l.protos[protoId])
	json.Unmarshal(data, a)
	return a
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
		l.protos[protoId] = p
		l.imports[protoId] = info.ImportPath
	}

	// configure prototypes
	for id, p := range l.protos {
		if l.Prototypes[id].Config == nil {
			continue
		}
		data, _ := json.Marshal(l.Prototypes[id].Config)
		err := json.Unmarshal(data, p)
		if err != nil {
			return prettyMarshalErr(id, data, err)
		}
	}

	// create agents from prototypes
	agents := []interface{}{}
	agentMap := map[string]interface{}{}
	for _, info := range l.Agents {
		a := l.NewAgentFromProto(info.ProtoId, nil)
		agentMap[info.Id] = a
		agents = append(agents, a)

		ag, ok := a.(Agent)
		if !ok {
			panic("loader: Agent '" + info.Id + "' does not implement required sim.Agent methods")
		}

		ag.SetId(info.Id)
		l.Engine.RegisterAll(ag)

		// register as service
		if info.IsService {
			if err := l.Engine.RegisterService(ag); err != nil {
				panic("loader: " + err.Error())
			}
		}
	}

	for i, info := range l.Agents {
		// set parents
		if c, ok := agents[i].(msg.Communicator); ok {
			if par, ok := agentMap[info.ParentId]; ok {
				c.SetParent(par.(msg.Communicator))
			}
		}

	}

	l.Engine.Load = l
	return nil
}

func prettyParseError(js string, err error) error {
	syntax, ok := err.(*json.SyntaxError)
	if !ok {
		return err
	}
	start, _ := strings.LastIndex(js[:syntax.Offset], "\n")+1, len(js)
	line, pos := strings.Count(js[:start], "\n"), int(syntax.Offset)-start-1
	msg := fmt.Sprint("loader: ", err, " at line ", line+1, " pos ", pos)
	return errors.New(msg)
}

func prettyMarshalErr(id string, data []byte, err error) error {
	msg := "loader: improper schema on prototype '" + id + "'" + ": " + err.Error()
	return errors.New(msg)
}
