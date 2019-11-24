# `signalprocessing` -- golang change point detection

## Overview

`singalprocessing` contains implementations of a `ChangeDetector`
interface in addition to CLI wrappers to facilitate testing. 

## API Use

In a go program, you can instantiate an object:

	var cd ChangeDetector
	
	cd = signalprocessing.NewQHatDetector(0.5, 1000, time.Now().Unix())

	// or

	cd = signalprocessing.NewEDMDetector(10)

	cp, err := cd.DetectChanges([]float64{...})
	if err != nil {
		return 
	}

	fmt.Println(cp)

See the [godoc](https://godoc.org/github.com/tychoish/signalprocessing) for
more information.

## CLI Use 

Fetch software and dependencies: 

    go get github.com/mongodb/grip
    go get github.com/tychoish/singalprocessing

Build binaries: 

	cd $GOPATH/src/github.com/tychoish/singalprocessing
	go build cmd/edivisive/edivisive.go
	go build cmd/edm/edm.go
	
Place binaries in the path: 

	mv edm /usr/local/bin/
	mv edivisive /usr/local/bin/

Example calls: 

	edivisive.go --inputJSON '[1,1,1,1,2,2,2,2,3,3,3]'                                                                                           
	edm --inputJSON '[1,1,1,1,1,1,1,1,1,1,2,2,2,2,2,2,2,2,2,20,20,20,20,55,55,55,55,55,3,3,3,3,3,3,3,3,3,3,3,3,3,3,3]'

Pass `--help` to either binary for full documentation of the input
options. 

## Development

The following work items: 

- Implement an outlier detection algorithm.

- Implement methods to convert []ChangePoint{} to bson and json
  directly.

- Improve test harness to be able to easily add validation workloads
  to the regression suite.
