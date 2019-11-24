package main

import (
	"encoding/json"
	"flag"
	"time"

	"github.com/mongodb/grip"
	"github.com/mongodb/grip/message"
	"github.com/mongodb/grip/send"
	"github.com/tychoish/signalprocessing"
)

func main() {
	var (
		seed         int64
		pvalue       float64
		permutations int
		inputJSON    string
	)

	flag.Int64Var(&seed, "seed", time.Now().Unix(), "specify the psudo randomness seed")
	flag.Float64Var(&pvalue, "pvalue", 0.5, "specify the pvalue")
	flag.IntVar(&permutations, "permutations", 1000, "specify the permutations")
	flag.StringVar(&inputJSON, "inputJSON", "", "pass a json string that is an array of floats")
	flag.Parse()

	sender := send.MakePlainLogger()
	grip.EmergencyFatal(sender.SetFormatter(send.MakePlainFormatter()))
	grip.EmergencyFatal(grip.SetSender(sender))

	grip.EmergencyFatal(message.When(inputJSON == "", "no input"))

	data := []float64{}
	err := json.Unmarshal([]byte(inputJSON), &data)
	grip.EmergencyFatal(message.WrapError(err, "problem parsing input data"))

	cd := signalprocessing.NewQHatDetector(pvalue, permutations, seed)
	cp, err := cd.DetectChanges(data)
	grip.EmergencyFatal(message.WrapError(err, "problem analyzing data"))

	out, err := json.Marshal(cp)
	grip.EmergencyFatal(message.WrapError(err, "writing output data"))
	grip.Info(out)
}
