export type LoadState = 'idle' | 'loading' | 'success' | 'error'

export type ActiveView = 'home' | 'accounts' | 'instances' | 'logs'

export type AccountStatusTone = 'ready' | 'pending' | 'disabled' | 'invalid' | 'unknown'
export type InstanceStatusTone = 'pending' | 'starting' | 'running' | 'reconnecting' | 'stopped' | 'error' | 'unknown'
export type ComfortLevel = 'comfort' | 'busy' | 'tense' | 'quiet'
export type ToastTone = 'success' | 'error'
export type LoginTaskStatus = 'idle' | 'initializing' | 'polling' | 'succeeded' | 'failed' | 'expired' | 'cancelled' | 'replaced'
export type AccountAuthStatus = 'logged_in' | 'not_logged_in' | 'auth_invalid' | 'unknown'

export type ViewStatusSummary = {
  clusterStatus: string
  totalInstances: number | null
  runningInstances: number | null
  uptimeLabel: string
}

export type ViewResourceSnapshot = {
  cpuPercent: number | null
  memoryPercent: number | null
  usedMemoryLabel: string
  totalMemoryLabel: string
  collectedAtLabel: string
  comfortLevel: ComfortLevel
}

export type ViewAccount = {
  id: string
  label: string
  note: string
  enabled: boolean
  authStatus: AccountAuthStatus
  authStatusTone: AccountStatusTone
  authStatusLabel: string
  hasToken: boolean
}

export type ViewInstance = {
  id: string
  accountId: string
  serverAddress: string
  statusTone: InstanceStatusTone
  statusLabel: string
  onlineDurationLabel: string
  health: number | null
  food: number | null
  positionLabel: string
}

export type ViewLogItem = {
  id: string
  timestampLabel: string
  actionLabel: string
  targetLabel: string
  details: string
  tone: ToastTone
}

export type HomeHeroView = {
  title: string
  subtitle: string
  comfortLevel: ComfortLevel
  pendingAccountCount: number
  activeInstanceRatioLabel: string
  cpuLabel: string
  memoryLabel: string
  sampleLabel: string
}

export type ToastMessage = {
  id: string
  tone: ToastTone
  message: string
}

export type LoginTaskView = {
  accountId: string
  taskStatus: LoginTaskStatus
  panelOpen: boolean
  userCode: string
  verificationUri: string
  verificationUriComplete: string
  expiresInSeconds: number
  expiresAt: number | null
  intervalSeconds: number
  lastMessage: string
  transientErrorCount: number
}
