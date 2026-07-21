package legalscan

import "testing"

func TestScanFindsPriceChange(t *testing.T) {
	res := ScanText("t.pdf", "Арендодатель вправе в одностороннем порядке изменить стоимость аренды кабинета.")
	if len(res.Flags) == 0 {
		t.Fatal("expected flags")
	}
	if res.Flags[0].Severity != "high" {
		t.Fatalf("severity=%s", res.Flags[0].Severity)
	}
}
