package backend

import "testing"

func TestHexCodeToInt(t *testing.T) {
	got, err := HexCode("#FFFFFF").HexCodeToInt()
	if err != nil {
		t.Fatalf("HexCodeToInt returned error: %v", err)
	}
	if got != 0xFFFFFF {
		t.Errorf("HexCodeToInt(#FFFFFF) = %d, want %d", got, 0xFFFFFF)
	}

	got, err = HexCode("#000000").HexCodeToInt()
	if err != nil {
		t.Fatalf("HexCodeToInt returned error: %v", err)
	}
	if got != 0 {
		t.Errorf("HexCodeToInt(#000000) = %d, want 0", got)
	}
}

func TestHexCodeToIntInvalid(t *testing.T) {
	if _, err := HexCode("#ZZZZZZ").HexCodeToInt(); err == nil {
		t.Error("HexCodeToInt did not return an error for invalid hex")
	}
}

func TestParseRGB(t *testing.T) {
	r, g, b, err := HexCode("#FF0000").ParseRGB()
	if err != nil {
		t.Fatalf("ParseRGB returned error: %v", err)
	}
	if r != 255 || g != 0 || b != 0 {
		t.Errorf("ParseRGB(#FF0000) = (%d,%d,%d), want (255,0,0)", r, g, b)
	}

	r, g, b, err = HexCode("#00FF00").ParseRGB()
	if err != nil {
		t.Fatalf("ParseRGB returned error: %v", err)
	}
	if r != 0 || g != 255 || b != 0 {
		t.Errorf("ParseRGB(#00FF00) = (%d,%d,%d), want (0,255,0)", r, g, b)
	}
}

func TestParseRGBInvalid(t *testing.T) {
	if _, _, _, err := HexCode("#ZZZZZZ").ParseRGB(); err == nil {
		t.Error("ParseRGB did not return an error for invalid hex")
	}
}

func TestCombineRGB(t *testing.T) {
	c, err := CombineRGB(255, 0, 0)
	if err != nil {
		t.Fatalf("CombineRGB returned error: %v", err)
	}
	if c != HexCode("#FF0000") {
		t.Errorf("CombineRGB(255,0,0) = %q, want #FF0000", c)
	}

	c, err = CombineRGB(0, 0, 0)
	if err != nil {
		t.Fatalf("CombineRGB returned error: %v", err)
	}
	if c != HexCode("#000000") {
		t.Errorf("CombineRGB(0,0,0) = %q, want #000000", c)
	}

	c, err = CombineRGB(18, 52, 86)
	if err != nil {
		t.Fatalf("CombineRGB returned error: %v", err)
	}
	if c != HexCode("#123456") {
		t.Errorf("CombineRGB(18,52,86) = %q, want #123456", c)
	}
}

func TestCombineRGBInvalid(t *testing.T) {
	for _, v := range []int{-1, 256} {
		if _, err := CombineRGB(v, 0, 0); err == nil {
			t.Errorf("CombineRGB(%d,0,0) did not return an error", v)
		}
	}
}

func TestMuteColorTarget(t *testing.T) {
	got, err := HexCode("#FFFFFF").MuteColorTarget(HexCode("#000000"), 0.5)
	if err != nil {
		t.Fatalf("MuteColorTarget returned error: %v", err)
	}
	if got != HexCode("#808080") {
		t.Errorf("MuteColorTarget(white->black @0.5) = %q, want #808080", got)
	}
}

func TestMuteColorTargetInvalid(t *testing.T) {
	if _, err := HexCode("#ZZZZZZ").MuteColorTarget(HexCode("#000000"), 0.5); err == nil {
		t.Error("MuteColorTarget did not return an error for invalid source")
	}
	if _, err := HexCode("#FFFFFF").MuteColorTarget(HexCode("#ZZZZZZ"), 0.5); err == nil {
		t.Error("MuteColorTarget did not return an error for invalid target")
	}
}
