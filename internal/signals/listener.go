package signals

import (
	"os"
	"os/signal"
)

type SignalListener struct {
	C chan os.Signal
}

func NewSignalListener(sig ...os.Signal) SignalListener {
	var sl = SignalListener{
		C: make(chan os.Signal, 1),
	}
	signal.Notify(sl.C, sig...)

	return sl
}
