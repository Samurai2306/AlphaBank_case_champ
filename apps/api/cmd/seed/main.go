package main

import (
	"fmt"

	"github.com/alphabank-case-champ/copilot-api/internal/adapters/memory"
)

func main() {
	s := memory.NewSeeded()
	p := s.Profile()
	fmt.Printf("seeded persona %s (%s) revenue=%.0f txns=%d piggy=%.0f\n",
		p.PersonaKey, p.DisplayName, p.MonthlyRevenueEstimate, len(s.Transactions()), s.Piggy().Balance)
}
