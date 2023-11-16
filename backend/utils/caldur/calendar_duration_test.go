package caldur_test

import (
	"errors"
	"testing"

	"github.com/stephenzsy/small-kms/backend/utils/caldur"
)

func TestParsePeriod(t *testing.T) {
	p, err := caldur.Parse("P1Y")
	if err != nil {
		t.Fatal(err)
	}
	if p.Year != 1 {
		t.Fatalf("expected 1 year, got %d year", p.Year)
	}

	p, err = caldur.Parse("2Y")
	if !errors.Is(err, caldur.ErrInvalidDuration) {
		t.Fatalf("expected ErrInvalidPeriod, got %v", err)
	}

	p, err = caldur.Parse("P2Y1M")
	if err != nil {
		t.Fatal(err)
	}
	if p.Year != 2 {
		t.Fatalf("expected 1 year, got %d year", p.Year)
	}
	if p.Month != 1 {
		t.Fatalf("expected 1 month, got %d month", p.Month)
	}

	p, err = caldur.Parse("P3M")
	if err != nil {
		t.Fatal(err)
	}
	if p.Year != 0 {
		t.Fatalf("expected 0 year, got %d year", p.Year)
	}
	if p.Month != 3 {
		t.Fatalf("expected 3 months, got %d months", p.Month)
	}

	p, err = caldur.Parse("P1Y2M3W4DT5H6M7S")
	if err != nil {
		t.Fatal(err)
	}
	if p.Year != 1 {
		t.Fatalf("expected 1 year, got %d year", p.Year)
	}
	if p.Month != 2 {
		t.Fatalf("expected 2 months, got %d months", p.Month)
	}
	if p.Week != 3 {
		t.Fatalf("expected 3 weeks, got %d weeks", p.Week)
	}
	if p.Day != 4 {
		t.Fatalf("expected 4 days, got %d days", p.Day)
	}
	if p.Hour != 5 {
		t.Fatalf("expected 5 hours, got %d hours", p.Hour)
	}
	if p.Minute != 6 {
		t.Fatalf("expected 6 minutes, got %d minutes", p.Minute)
	}
	if p.Second != 7 {
		t.Fatalf("expected 7 seconds, got %d seconds", p.Second)
	}

}
