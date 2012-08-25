
type Request struct {
  
}

func OfferBalance

type Econ struct {
   id string
   req map[msg.Communicator]map[time.Time]float64
   off map[msg.Communicator]map[time.Time]float64
   Tracked []string
   eng *sim.Engine
}

func (e *Econ) SetId(id string) {
  m.id = id
}

func (e *Econ) Id() string {
  return m.id
}

func (e *Econ) Start(eng *sim.Engine) {
  e.eng = eng
}

func (e *Econ) Receive(m *msg.Message) {
  if kind, ok := m.Payload.(string); ok {
    commod := strings.Split("::").(string)
    mkt := 
  }
}

func (e *Econ) MsgNotify(m *msg.Message) {
  if !e.isTracked(m.Owner) || m.Trans == nil {
    return
  }

  qty = mg.Trans.Resource().Qty()
  if mg.Trans.Type() == trans.Offer {
    e.off[m.Owner][e.tm] += qty
  } else {
    e.req[m.Owner][e.tm] += qty
  }
}

func (e *Econ) isTracked(c msg.Communicator) bool {
  if a, ok := c.(sim.Agent); ok {
    for _, id := range e.Tracked {
      if a.Id() == id {
        return true
      }
    }
  }
  return false
}

func (e *Econ) Parent() msg.Communicator {
  return nil
}

func (e *Econ) SetParent(par msg.Communicator) {
}

