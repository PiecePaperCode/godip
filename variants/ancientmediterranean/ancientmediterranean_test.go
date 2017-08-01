package ancientmediterranean

import (
	"reflect"
	"testing"

	"github.com/zond/godip/state"
	"github.com/zond/godip/variants/classical/orders"

	cla "github.com/zond/godip/variants/classical/common"
	dip "github.com/zond/godip/common"
)

func init() {
	dip.Debug = true
}

func assertOrderValidity(t *testing.T, validator dip.Validator, order dip.Order, nat dip.Nation, err error) {
	if gotNat, e := order.Validate(validator); e != err {
		t.Errorf("%v should validate to %v, but got %v", order, err, e)
	} else if gotNat != nat {
		t.Errorf("%v should validate with %q as issuer, but got %q", order, nat, gotNat)
	}
}

func assertMove(t *testing.T, j *state.State, src, dst dip.Province, success bool) {
	if success {
		unit, _, ok := j.Unit(src)
		if !ok {
			t.Errorf("Should be a unit at %v", src)
		}
		j.SetOrder(src, orders.Move(src, dst))
		j.Next()
		if err, ok := j.Resolutions()[src]; ok && err != nil {
			t.Errorf("Move from %v to %v should have worked, got %v", src, dst, err)
		}
		if now, _, ok := j.Unit(src); ok && reflect.DeepEqual(now, unit) {
			t.Errorf("%v should have moved from %v", now, src)
		}
		if now, _, _ := j.Unit(dst); !reflect.DeepEqual(now, unit) {
			t.Errorf("%v should be at %v now", unit, dst)
		}
	} else {
		unit, _, _ := j.Unit(src)
		j.SetOrder(src, orders.Move(src, dst))
		j.Next()
		if err, ok := j.Resolutions()[src]; !ok || err == nil {
			t.Errorf("Move from %v to %v should not have worked", src, dst)
		}
		if now, _, _ := j.Unit(src); !reflect.DeepEqual(now, unit) {
			t.Errorf("%v should not have moved from %v", now, src)
		}
	}
}

func startState(t *testing.T) *state.State {
	judge, err := AncientMediterraneanStart()
	if err != nil {
		t.Fatalf("%v", err)
	}
	return judge
}

func TestByzantium(t *testing.T) {
	judge := startState(t)

	// Test can't bypass Byzantium.
	judge.SetUnit("aeg", dip.Unit{cla.Fleet, Greece})
	assertOrderValidity(t, judge, orders.Move("aeg", "bla"), "", cla.ErrIllegalMove)
	
	// Test can sail through it.
	assertOrderValidity(t, judge, orders.Move("aeg", "byz"), Greece, nil)
	judge.SetUnit("byz", dip.Unit{cla.Fleet, Greece})
	assertOrderValidity(t, judge, orders.Move("byz", "bla"), Greece, nil)
	
	// Check can't convoy via Byzantium.
	assertOrderValidity(t, judge, orders.Move("ath", "bit"), "", cla.ErrMissingConvoyPath)
}

func TestHighSeas(t *testing.T) {
	judge := startState(t)
	
	judge.SetUnit("aus", dip.Unit{cla.Fleet, Greece})
	assertOrderValidity(t, judge, orders.Move("aus", "mes"), Greece, nil)
	assertOrderValidity(t, judge, orders.Move("aus", "lib"), Greece, nil)
	assertOrderValidity(t, judge, orders.Move("aus", "got"), Greece, nil)
	judge.SetUnit("mes", dip.Unit{cla.Fleet, Greece})
	assertOrderValidity(t, judge, orders.Move("mes", "lib"), Greece, nil)
	assertOrderValidity(t, judge, orders.Move("mes", "got"), Greece, nil)
	judge.SetUnit("lib", dip.Unit{cla.Fleet, Greece})
	assertOrderValidity(t, judge, orders.Move("lib", "got"), Greece, nil)
}

func TestDiolkos(t *testing.T) {
	judge := startState(t)

	// Test can't bypass Athens or Sparta.
	judge.SetUnit("ion", dip.Unit{cla.Fleet, Greece})
	assertOrderValidity(t, judge, orders.Move("ion", "aeg"), "", cla.ErrIllegalMove)
	
	// Test can walk from Athens to Sparta.
	assertOrderValidity(t, judge, orders.Move("ath", "spa"), Greece, nil)
}

func TestSicily(t *testing.T) {
	judge := startState(t)

	// Test can walk or sail between Sicily and Neapolis.
	judge.SetUnit("sic", dip.Unit{cla.Army, Rome})
	assertOrderValidity(t, judge, orders.Move("sic", "nea"), Rome, nil)
	assertOrderValidity(t, judge, orders.Move("nea", "sic"), Rome, nil)
	
	// Test can sail through the 'Strait of Messina'.
	judge.SetUnit("tys", dip.Unit{cla.Fleet, Rome})
	assertOrderValidity(t, judge, orders.Move("tys", "aus"), Rome, nil)
}

func TestCorsica(t *testing.T) {
	judge := startState(t)

	// Test can walk or sail between Corsica and Sardinia.
	judge.SetUnit("cor", dip.Unit{cla.Army, Rome})
	judge.SetUnit("sad", dip.Unit{cla.Fleet, Rome})
	assertOrderValidity(t, judge, orders.Move("cor", "sad"), Rome, nil)
	assertOrderValidity(t, judge, orders.Move("sad", "cor"), Rome, nil)
}

func TestNileDelta(t *testing.T) {
	judge := startState(t)
	
	// Happy paths near Nile Delta
	assertOrderValidity(t, judge, orders.Move("the", "sii"), Egypt, nil)
	assertOrderValidity(t, judge, orders.SupportHold("the", "mem"), Egypt, nil)
	assertOrderValidity(t, judge, orders.SupportHold("the", "ale"), Egypt, nil)
	assertOrderValidity(t, judge, orders.SupportMove("mem", "the", "ale"), Egypt, nil)
	assertOrderValidity(t, judge, orders.Move("ale", "sii"), Egypt, nil)
	assertOrderValidity(t, judge, orders.SupportMove("the", "ale", "sii"), Egypt, nil)
	judge.SetUnit("gop", dip.Unit{cla.Fleet, Rome})
	assertOrderValidity(t, judge, orders.Move("the", "jer"), Egypt, nil)
	
	// Illegal moves near Nile Delta
	judge.SetUnit("ree", dip.Unit{cla.Fleet, Rome})
	assertOrderValidity(t, judge, orders.Move("ree", "ale"), "", cla.ErrIllegalMove)
	assertOrderValidity(t, judge, orders.Move("ree", "gop"), "", cla.ErrIllegalMove)
	judge.RemoveUnit("mem")
	judge.SetUnit("mem", dip.Unit{cla.Fleet, Rome})
	assertOrderValidity(t, judge, orders.Move("mem", "sii"), "", cla.ErrIllegalMove)
	assertOrderValidity(t, judge, orders.Move("mem", "gop"), "", cla.ErrIllegalMove)
}

func TestConvoyBaleares(t *testing.T) {
	judge := startState(t)

	// Test convoys through Baleares.
	judge.SetUnit("sag", dip.Unit{cla.Army, Rome})
	judge.SetUnit("bal", dip.Unit{cla.Fleet, Rome})
	judge.SetUnit("lig", dip.Unit{cla.Fleet, Rome})
	assertOrderValidity(t, judge, orders.Move("sag", "cor"), Rome, nil)
	assertOrderValidity(t, judge, orders.Convoy("bal", "sag", "cor"), Rome, nil)

	// Test an army in Baleares can't be part of a convoy chain.
	judge.RemoveUnit("bal")
	judge.SetUnit("bal", dip.Unit{cla.Army, Rome})
	assertOrderValidity(t, judge, orders.Move("sag", "cor"), "", cla.ErrMissingConvoyPath)
	assertOrderValidity(t, judge, orders.Convoy("bal", "sag", "cor"), "", cla.ErrIllegalConvoyer)
}

