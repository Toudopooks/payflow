package service

import (
	"testing"
	"time"
)

func TestBuildInstallmentsDistributesRemainder(t *testing.T) {
	items, err := BuildInstallments("plan", 10000, 3, time.Date(2026, 1, 31, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	want := []int64{3334, 3333, 3333}
	var total int64
	for i, item := range items {
		if item.Amount != want[i] {
			t.Errorf("installment %d: got %d, want %d", i, item.Amount, want[i])
		}
		total += item.Amount
	}
	if total != 10000 {
		t.Fatalf("got total %d", total)
	}
	if got := items[1].DueAt.Day(); got != 28 {
		t.Fatalf("February due day = %d, want 28", got)
	}
}

func TestBuildInstallmentsRejectsZeroValueInstallments(t *testing.T) {
	if _, err := BuildInstallments("plan", 1, 4, time.Now()); err == nil {
		t.Fatal("expected error")
	}
}

func TestBuildInstallmentsRejectsUnsupportedCount(t *testing.T) {
	if _, err := BuildInstallments("plan", 1000, 5, time.Now()); err == nil {
		t.Fatal("expected error")
	}
}
