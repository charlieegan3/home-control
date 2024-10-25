package plugs

import "testing"

func TestParseStatus_ON(t *testing.T) {
	rawStatus := `
{t}</table><hr/>{t}{s}</th><th></th><th style='text-align:center'><th></th><td>{e}{s}Voltage{m}</td><td style='text-align:left'>248</td><td>&nbsp;</td><td> V{e}{s}Current{m}</td><td style='text-align:left'>0.439</td><td>&nbsp;</td><td> A{e}{s}Active Power{m}</td><td style='text-align:left'>77</td><td>&nbsp;</td><td> W{e}{s}Apparent Power{m}</td><td style='text-align:left'>109</td><td>&nbsp;</td><td> VA{e}{s}Reactive Power{m}</td><td style='text-align:left'>78</td><td>&nbsp;</td><td> VAr{e}{s}Power Factor{m}</td><td style='text-align:left'>0.70</td><td>&nbsp;</td><td>{e}{s}Energy Today{m}</td><td style='text-align:left'>1.689</td><td>&nbsp;</td><td> kWh{e}{s}Energy Yesterday{m}</td><td style='text-align:left'>1.834</td><td>&nbsp;</td><td> kWh{e}{s}Energy Total{m}</td><td style='text-align:left'>59.540</td><td>&nbsp;</td><td> kWh{e}</table><hr/>{t}</table>{t}<tr><td style='width:100%;text-align:center;font-weight:bold;font-size:62px'>ON</td></tr><tr></tr></table>
`

	parsedStatus, err := ParseStatus(rawStatus)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if parsedStatus.PowerOn != true {
		t.Fatalf("unexpected active power: %v", parsedStatus.PowerOn)
	}

	if parsedStatus.ActivePower != 77 {
		t.Fatalf("unexpected active power: %v", parsedStatus.ActivePower)
	}
}

func TestParseStatus_OFF(t *testing.T) {
	rawStatus := `
{t}</table><hr/>{t}{s}</th><th></th><th style='text-align:center'><th></th><td>{e}{s}Voltage{m}</td><td style='text-align:left'>248</td><td>&nbsp;</td><td> V{e}{s}Current{m}</td><td style='text-align:left'>0.439</td><td>&nbsp;</td><td> A{e}{s}Active Power{m}</td><td style='text-align:left'>76</td><td>&nbsp;</td><td> W{e}{s}Apparent Power{m}</td><td style='text-align:left'>109</td><td>&nbsp;</td><td> VA{e}{s}Reactive Power{m}</td><td style='text-align:left'>78</td><td>&nbsp;</td><td> VAr{e}{s}Power Factor{m}</td><td style='text-align:left'>0.70</td><td>&nbsp;</td><td>{e}{s}Energy Today{m}</td><td style='text-align:left'>1.689</td><td>&nbsp;</td><td> kWh{e}{s}Energy Yesterday{m}</td><td style='text-align:left'>1.834</td><td>&nbsp;</td><td> kWh{e}{s}Energy Total{m}</td><td style='text-align:left'>59.540</td><td>&nbsp;</td><td> kWh{e}</table><hr/>{t}</table>{t}<tr><td style='width:100%;text-align:center;font-weight:bold;font-size:62px'>OFF</td></tr><tr></tr></table>
`

	parsedStatus, err := ParseStatus(rawStatus)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if parsedStatus.PowerOn != false {
		t.Fatalf("unexpected active power: %v", parsedStatus.PowerOn)
	}

	if parsedStatus.ActivePower != 76 {
		t.Fatalf("unexpected active power: %v", parsedStatus.ActivePower)
	}
}

func TestParseStatus_InvalidData(t *testing.T) {
	rawStatus := `
Invalid data that should cause parsing to fail
`

	_, err := ParseStatus(rawStatus)
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}
