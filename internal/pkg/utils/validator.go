package utils

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// 自定义验证器
	validate *validator.Validate
)

// InitValidator 初始化验证器
func InitValidator() {
	validate = validator.New()

	// 注册自定义验证器
	validate.RegisterValidation("git_tag", validateGitTag)
	validate.RegisterValidation("git_commit", validateGitCommit)
	validate.RegisterValidation("repository_url", validateRepositoryURL)
	validate.RegisterValidation("environment_type", validateEnvironmentType)
	validate.RegisterValidation("application_type", validateApplicationType)
}

// GetValidator 获取验证器
func GetValidator() *validator.Validate {
	if validate == nil {
		InitValidator()
	}
	return validate
}

// validateGitTag 验证 Git 标签格式
func validateGitTag(fl validator.FieldLevel) bool {
	tag := fl.Field().String()
	// 匹配语义化版本号格式：v1.0.0, v1.0.0-alpha, v1.0.0-beta.1 等
	pattern := `^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	matched, _ := regexp.MatchString(pattern, tag)
	return matched
}

// validateGitCommit 验证 Git 提交哈希格式
func validateGitCommit(fl validator.FieldLevel) bool {
	commit := fl.Field().String()
	// 匹配 40 位十六进制字符的 SHA-1 哈希
	pattern := `^[a-f0-9]{40}$`
	matched, _ := regexp.MatchString(pattern, commit)
	return matched
}

// validateRepositoryURL 验证仓库 URL 格式
func validateRepositoryURL(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	// 匹配 GitHub 仓库 URL
	pattern := `^https://github\.com/[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+(?:\.git)?$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

// validateEnvironmentType 验证环境类型
func validateEnvironmentType(fl validator.FieldLevel) bool {
	envType := fl.Field().String()
	validTypes := []string{"kubernetes", "physical"}
	for _, validType := range validTypes {
		if envType == validType {
			return true
		}
	}
	return false
}

// validateApplicationType 验证应用类型
func validateApplicationType(fl validator.FieldLevel) bool {
	appType := fl.Field().String()
	validTypes := []string{"microservice", "monolith", "frontend", "backend", "api"}
	for _, validType := range validTypes {
		if appType == validType {
			return true
		}
	}
	return false
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
	return GetValidator().Struct(s)
}

// ValidateVar 验证单个字段
func ValidateVar(field interface{}, tag string) error {
	return GetValidator().Var(field, tag)
}

// GetValidationErrors 获取验证错误详情
func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = field + " is required"
			case "min":
				errors[field] = field + " must be at least " + e.Param()
			case "max":
				errors[field] = field + " must be at most " + e.Param()
			case "email":
				errors[field] = field + " must be a valid email address"
			case "url":
				errors[field] = field + " must be a valid URL"
			case "git_tag":
				errors[field] = field + " must be a valid git tag (e.g., v1.0.0)"
			case "git_commit":
				errors[field] = field + " must be a valid git commit hash"
			case "repository_url":
				errors[field] = field + " must be a valid GitHub repository URL"
			case "environment_type":
				errors[field] = field + " must be either 'kubernetes' or 'physical'"
			case "application_type":
				errors[field] = field + " must be a valid application type"
			default:
				errors[field] = field + " is invalid"
			}
		}
	}

	return errors
}
