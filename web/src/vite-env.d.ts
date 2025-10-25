/// <reference types="vite/client" />

interface ImportMetaEnv {
	readonly VITE_USE_MOCK: string
	// 添加更多环境变量...
}

interface ImportMeta {
	readonly env: ImportMetaEnv
}
