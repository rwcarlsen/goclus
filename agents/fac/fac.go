
package fac

type Fac struct {
  queuedOrders []*msg.Message

  inCommod string
  inUnits string
  inBuff *rsrc.Buffer

  outCommod string
  outUnits string
  outBuff *rsrc.Buffer

  createRate float64
  convertAmt float64
  convertPeriod int
  convertOffset int
}

func (f *Fac) Tick(tm time.Duration) {
  
}

func (f *Fac) Tock(tm time.Duration) {
  
}

func (f *Fac) Receive(m *msg.Message) {
  
}
