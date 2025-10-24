import { create } from 'zustand'
import type {
  Version,
  Application,
  Environment,
  Deployment,
  DeploymentDetail,
} from '@/types'

interface AppState {
  versions: Version[]
  applications: Application[]
  environments: Environment[]
  deployments: Deployment[]
  currentDeployment: DeploymentDetail | null
  
  setVersions: (versions: Version[]) => void
  setApplications: (applications: Application[]) => void
  setEnvironments: (environments: Environment[]) => void
  setDeployments: (deployments: Deployment[]) => void
  setCurrentDeployment: (deployment: DeploymentDetail | null) => void
  
  addDeployment: (deployment: Deployment) => void
  updateDeployment: (id: string, updates: Partial<Deployment>) => void
}

export const useAppStore = create<AppState>((set) => ({
  versions: [],
  applications: [],
  environments: [],
  deployments: [],
  currentDeployment: null,
  
  setVersions: (versions) => set({ versions }),
  setApplications: (applications) => set({ applications }),
  setEnvironments: (environments) => set({ environments }),
  setDeployments: (deployments) => set({ deployments }),
  setCurrentDeployment: (currentDeployment) => set({ currentDeployment }),
  
  addDeployment: (deployment) =>
    set((state) => ({
      deployments: [deployment, ...state.deployments],
    })),
  
  updateDeployment: (id, updates) =>
    set((state) => ({
      deployments: state.deployments.map((d) =>
        d.id === id ? { ...d, ...updates } : d
      ),
    })),
}))

interface UIState {
  sidebarCollapsed: boolean
  toggleSidebar: () => void
}

export const useUIStore = create<UIState>((set) => ({
  sidebarCollapsed: false,
  toggleSidebar: () =>
    set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),
}))
