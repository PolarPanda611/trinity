package main

import "github.com/PolarPanda611/trinity"

func main() {
	t := trinity.New("local")
	// rg := t.NewAPIGroup("/api/v1")
	// trinity.NewAPIInGroup(rg, "exchange_rate_types", servicev1.ExchangeRateTypeViewSet, []string{"Retrieve", "List", "Create", "Update", "Delete"})
	t.Serve()
}
