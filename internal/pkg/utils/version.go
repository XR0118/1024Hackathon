package utils

import (
	"strconv"
	"strings"
)

// CompareVersions 比较两个版本号
// 返回值: -1 表示 v1 < v2, 0 表示 v1 == v2, 1 表示 v1 > v2
// 支持格式: v1.2.3, 1.2.3, v1.2.3-alpha, v1.2.3-beta.1 等
func CompareVersions(v1, v2 string) int {
	// 移除 'v' 前缀
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	// 分离主版本号和预发布标签
	v1Main, v1Pre := splitVersion(v1)
	v2Main, v2Pre := splitVersion(v2)

	// 比较主版本号
	mainCmp := compareMainVersion(v1Main, v2Main)
	if mainCmp != 0 {
		return mainCmp
	}

	// 主版本号相同，比较预发布标签
	return comparePrereleaseVersion(v1Pre, v2Pre)
}

// splitVersion 分离主版本号和预发布标签
func splitVersion(version string) (string, string) {
	parts := strings.SplitN(version, "-", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

// compareMainVersion 比较主版本号
func compareMainVersion(v1, v2 string) int {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	// 确保至少有3个部分
	for len(v1Parts) < 3 {
		v1Parts = append(v1Parts, "0")
	}
	for len(v2Parts) < 3 {
		v2Parts = append(v2Parts, "0")
	}

	// 逐个比较
	for i := 0; i < 3; i++ {
		n1, _ := strconv.Atoi(v1Parts[i])
		n2, _ := strconv.Atoi(v2Parts[i])

		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
	}

	return 0
}

// comparePrereleaseVersion 比较预发布版本标签
func comparePrereleaseVersion(pre1, pre2 string) int {
	if pre1 == "" && pre2 == "" {
		return 0
	}
	if pre1 == "" && pre2 != "" {
		return 1
	}
	if pre1 != "" && pre2 == "" {
		return -1
	}
	if pre1 < pre2 {
		return -1
	}
	if pre1 > pre2 {
		return 1
	}
	return 0
}

// IsVersionGreaterOrEqual 判断 v1 >= v2
func IsVersionGreaterOrEqual(v1, v2 string) bool {
	return CompareVersions(v1, v2) >= 0
}
