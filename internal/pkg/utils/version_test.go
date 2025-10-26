package utils

import (
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		// 基本版本比较
		{"相等版本", "v1.0.0", "v1.0.0", 0},
		{"v1 < v2", "v1.0.0", "v1.1.0", -1},
		{"v1 > v2", "v1.1.0", "v1.0.0", 1},

		// 不带 v 前缀
		{"无前缀相等", "1.0.0", "1.0.0", 0},
		{"无前缀 v1 < v2", "1.0.0", "1.1.0", -1},
		{"无前缀 v1 > v2", "1.1.0", "1.0.0", 1},

		// 跨主版本
		{"主版本 v1 < v2", "v1.9.9", "v2.0.0", -1},
		{"主版本 v1 > v2", "v2.0.0", "v1.9.9", 1},

		// 次版本
		{"次版本 v1 < v2", "v1.1.0", "v1.2.0", -1},
		{"次版本 v1 > v2", "v1.2.0", "v1.1.0", 1},

		// 修订号
		{"修订号 v1 < v2", "v1.0.1", "v1.0.2", -1},
		{"修订号 v1 > v2", "v1.0.2", "v1.0.1", 1},

		// 预发布版本
		{"预发布 < 正式", "v1.0.0-alpha", "v1.0.0", -1},
		{"正式 > 预发布", "v1.0.0", "v1.0.0-alpha", 1},
		{"预发布比较", "v1.0.0-alpha", "v1.0.0-beta", -1},

		// 缺少修订号
		{"缺少修订号1", "v1.0", "v1.0.0", 0},
		{"缺少修订号2", "v1", "v1.0.0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompareVersions(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("CompareVersions(%q, %q) = %d, want %d", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

func TestIsVersionGreaterOrEqual(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want bool
	}{
		{"相等", "v1.0.0", "v1.0.0", true},
		{"更大", "v1.1.0", "v1.0.0", true},
		{"更小", "v1.0.0", "v1.1.0", false},
		{"主版本更大", "v2.0.0", "v1.9.9", true},
		{"次版本更大", "v1.2.0", "v1.1.9", true},
		{"修订号更大", "v1.0.2", "v1.0.1", true},
		{"预发布版本", "v1.0.0", "v1.0.0-alpha", true},
		{"预发布版本2", "v1.0.0-alpha", "v1.0.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsVersionGreaterOrEqual(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("IsVersionGreaterOrEqual(%q, %q) = %v, want %v", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

// 测试实际使用场景
func TestVersionCoverageScenario(t *testing.T) {
	targetVersion := "v1.2.0"

	tests := []struct {
		version   string
		isCovered bool
	}{
		{"v1.3.0", true},  // >= v1.2.0
		{"v1.2.5", true},  // >= v1.2.0
		{"v1.2.0", true},  // == v1.2.0
		{"v1.1.0", false}, // < v1.2.0
		{"v1.0.0", false}, // < v1.2.0
		{"v2.0.0", true},  // >= v1.2.0
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			got := IsVersionGreaterOrEqual(tt.version, targetVersion)
			if got != tt.isCovered {
				t.Errorf("version %s, expected covered=%v, got %v", tt.version, tt.isCovered, got)
			}
		})
	}
}
