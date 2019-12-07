package tcpserver

import "testing"

func Test_splitName(t *testing.T) {
	tests := []struct {
		name       string
		pointName  string
		grouping   bool
		wantName   string
		wantKey    string
		wantMetric string
	}{
		{"single", "disk", false, "disk", "value", ""},
		{"single grouped", "disk", true, "disk", "value", ""},
		{"double", "disk.free", false, "disk", "free", ""},
		{"double grouped", "disk.free", true, "disk", "value", "free"},
		{"triple", "disk.var.free", false, "disk", "var.free", ""},
		{"triple grouped", "disk.var.free", true, "disk", "_free", "var"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotKey, gotMetric := splitName(tt.pointName, tt.grouping)
			if gotName != tt.wantName {
				t.Errorf("splitName() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotKey != tt.wantKey {
				t.Errorf("splitName() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if gotMetric != tt.wantMetric {
				t.Errorf("splitName() gotMetric = %v, want %v", gotMetric, tt.wantMetric)
			}
		})
	}
}
