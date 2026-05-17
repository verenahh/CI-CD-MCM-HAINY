package model

import "testing"

func TestValidateEmptyName(t *testing.T) {
	p := Product{Name: "", Price: 10.0}
	if p.Validate() {
		t.Error("expected validation to fail for empty name")
	}
}

func TestValidateNegativePrice(t *testing.T) {
	p := Product{Name: "Aktuelle Strompreise an da Börse", Price: -5.0}
	if p.Validate() {
		t.Error("expected validation to fail for negative price")
	}
}

func TestValidateValidProduct(t *testing.T) {
	p := Product{Name: "sehr vailde", Price: 4.00}
	if !p.Validate() {
		t.Error("expected validation to pass for valid product")
	}
}

func TestValidateZeroPrice(t *testing.T) {
	p := Product{Name: "Nix wert", Price: 0}
	if !p.Validate() {
		t.Error("expected validation to pass for zero price")
	}
}
