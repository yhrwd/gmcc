import { formatBytes, formatDateTime, formatPositionLabel, formatUptime } from '@/lib/format'
import type {
  ApiAccount,
  ApiInstance,
  ApiOperationLog,
  ApiResourcesResponse,
  ApiStatusResponse,
} from '@/types/api'
import type {
  AccountAuthStatus,
  ComfortLevel,
  HomeHeroView,
  ViewAccount,
  ViewInstance,
  ViewLogItem,
  ViewResourceSnapshot,
  ViewStatusSummary,
} from '@/types/view'

export function normalizeAccountAuthStatus(authStatus?: string): AccountAuthStatus {
  if (authStatus === 'logged_in') {
    return 'logged_in'
  }
  if (authStatus === 'not_logged_in') {
    return 'not_logged_in'
  }
  if (authStatus === 'auth_invalid') {
    return 'auth_invalid'
  }
  return 'unknown'
}

export function mapAccountStatus(account: Pick<ApiAccount, 'enabled' | 'auth_status' | 'has_token'>) {
  const authStatus = normalizeAccountAuthStatus(account.auth_status)
  if (!account.enabled) {
    return 'disabled' as const
  }
  if (authStatus === 'logged_in') {
    return 'ready' as const
  }
  if (authStatus === 'auth_invalid') {
    return 'invalid' as const
  }
  if (authStatus === 'not_logged_in') {
    return 'pending' as const
  }
  return 'unknown' as const
}

export function mapAccountStatusLabel(status: ReturnType<typeof mapAccountStatus>) {
  switch (status) {
    case 'ready':
      return '在线冒险者'
    case 'pending':
      return '待唤醒'
    case 'disabled':
      return '已休眠'
    case 'invalid':
      return '状态失效'
    default:
      return '状态未识别'
  }
}

export function mapInstanceStatus(status?: string) {
  if (status === 'pending') {
    return 'pending' as const
  }
  if (status === 'starting') {
    return 'starting' as const
  }
  if (status === 'running') {
    return 'running' as const
  }
  if (status === 'reconnecting') {
    return 'reconnecting' as const
  }
  if (status === 'stopped') {
    return 'stopped' as const
  }
  if (status === 'error') {
    return 'error' as const
  }
  return 'unknown' as const
}

export function mapInstanceStatusLabel(status: ReturnType<typeof mapInstanceStatus>) {
  switch (status) {
    case 'pending':
      return '待出勤'
    case 'starting':
      return '启动中'
    case 'running':
      return '出勤中'
    case 'reconnecting':
      return '重连中'
    case 'stopped':
      return '休息中'
    case 'error':
      return '异常中'
    default:
      return '状态不明'
  }
}

export function mapComfortLevel(input: {
  clusterStatus?: string
  cpuPercent?: number | null
  memoryPercent?: number | null
}): ComfortLevel {
  const peak = Math.max(input.cpuPercent ?? 0, input.memoryPercent ?? 0)
  if (input.clusterStatus === 'stopped') {
    return 'quiet'
  }
  if (peak >= 85) {
    return 'tense'
  }
  if (peak >= 60) {
    return 'busy'
  }
  return 'comfort'
}

export function mapStatusSummary(status: ApiStatusResponse): ViewStatusSummary {
  return {
    clusterStatus: status.cluster_status || 'unknown',
    totalInstances: status.total_instances ?? null,
    runningInstances: status.running_instances ?? null,
    uptimeLabel: formatUptime(status.uptime),
  }
}

export function mapResourceSnapshot(resources: ApiResourcesResponse, clusterStatus?: string): ViewResourceSnapshot {
  const memory = resources.memory ?? {}
  return {
    cpuPercent: resources.cpu_percent ?? null,
    memoryPercent: memory.used_percent ?? null,
    usedMemoryLabel: formatBytes(memory.used_bytes),
    totalMemoryLabel: formatBytes(memory.total_bytes),
    collectedAtLabel: formatDateTime(resources.collected_at),
    comfortLevel: mapComfortLevel({
      clusterStatus,
      cpuPercent: resources.cpu_percent ?? null,
      memoryPercent: memory.used_percent ?? null,
    }),
  }
}

export function mapAccount(account: ApiAccount): ViewAccount {
  const authStatus = normalizeAccountAuthStatus(account.auth_status)
  const authStatusTone = mapAccountStatus(account)
  return {
    id: account.id,
    label: account.label?.trim() || account.id,
    note: account.note?.trim() || '',
    enabled: account.enabled,
    authStatus,
    authStatusTone,
    authStatusLabel: mapAccountStatusLabel(authStatusTone),
    hasToken: Boolean(account.has_token),
  }
}

export function mapInstance(instance: ApiInstance): ViewInstance {
  const statusTone = mapInstanceStatus(instance.status)
  return {
    id: instance.id,
    accountId: instance.account_id,
    serverAddress: instance.server_address?.trim() || '未填写地址',
    statusTone,
    statusLabel: mapInstanceStatusLabel(statusTone),
    onlineDurationLabel: instance.online_duration?.trim() || '刚刚开始记录',
    health: typeof instance.health === 'number' ? instance.health : null,
    food: typeof instance.food === 'number' ? instance.food : null,
    positionLabel: formatPositionLabel(instance.position),
  }
}

export function mapLogItem(log: ApiOperationLog): ViewLogItem {
  const target = log.target_instance_id || log.target_account_id || '基地广播'
  return {
    id: log.id,
    timestampLabel: formatDateTime(log.timestamp),
    actionLabel: log.action,
    targetLabel: target,
    details: log.details || log.error_msg || '发生了一次操作',
    tone: log.success ? 'success' : 'error',
  }
}

export function buildHeroView(input: {
  status: ViewStatusSummary | null
  resources: ViewResourceSnapshot | null
  accounts: ViewAccount[]
}): HomeHeroView {
  const pendingCount = input.accounts.filter((account) => account.authStatusTone !== 'ready' && account.authStatusTone !== 'disabled').length
  const total = input.status?.totalInstances ?? 0
  const running = input.status?.runningInstances ?? 0
  const comfortLevel = input.resources?.comfortLevel || 'quiet'

  return {
    title: comfortLevel === 'comfort' ? '系统状态稳定' : comfortLevel === 'busy' ? '系统负载上升' : comfortLevel === 'tense' ? '系统资源紧张' : '系统已停止',
    subtitle: pendingCount > 0 ? `当前有 ${pendingCount} 个账号待处理` : '当前可以继续创建或启动实例',
    comfortLevel,
    pendingAccountCount: pendingCount,
    activeInstanceRatioLabel: total > 0 ? `${running}/${total} 运行中` : '暂无运行中的实例',
    cpuLabel: input.resources?.cpuPercent !== null && input.resources?.cpuPercent !== undefined ? `${Math.round(input.resources.cpuPercent)}%` : '--',
    memoryLabel: input.resources?.memoryPercent !== null && input.resources?.memoryPercent !== undefined ? `${Math.round(input.resources.memoryPercent)}%` : '--',
    sampleLabel: input.resources?.collectedAtLabel || '--',
  }
}
