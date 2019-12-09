package main

import (
	"flag"

	"github.com/PolarPanda611/trinity"
)

var runmode = flag.String("runmode", "release", "running mode options: debug , test , release ")

func main() {
	flag.Parse()
	t := trinity.New(*runmode)
	// rg := t.NewAPIGroup("/api/v1")
	// trinity.NewAPIInGroup(rg, "exchange_rate_types", servicev1.ExchangeRateTypeViewSet, []string{"Retrieve", "List", "Create", "Update", "Delete"})
	t.Serve()
}
