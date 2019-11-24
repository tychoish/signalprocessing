package main

import (
	"encoding/json"
	"flag"

	"github.com/mongodb/grip"
	"github.com/mongodb/grip/message"
	"github.com/mongodb/grip/send"
	"github.com/tychoish/signalprocessing"
)

func main() {
	var (
		minSize   int
		inputJSON string
	)
	flag.IntVar(&minSize, "minSize", 10, "specify the minimum size")
	flag.StringVar(&inputJSON, "inputJSON", "", "pass a json string that is an array of floats")
	flag.Parse()

	sender := send.MakePlainLogger()
	grip.EmergencyFatal(sender.SetFormatter(send.MakePlainFormatter()))
	grip.EmergencyFatal(grip.SetSender(sender))

	data := []float64{}
	err := json.Unmarshal([]byte(inputJSON), &data)
	grip.EmergencyFatal(message.WrapError(err, "problem parsing input data"))

	cd := signalprocessing.NewEDMDetector(minSize)
	cp, err := cd.DetectChanges(data)
	grip.EmergencyFatal(message.WrapError(err, "problem analyzing data"))

	out, err := json.Marshal(cp)
	grip.EmergencyFatal(message.WrapError(err, "writing output data"))
	grip.Info(out)
}
