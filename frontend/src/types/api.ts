export type ApiResponse = {
  success?: boolean
  message?: string
  error?: string
}

export type ApiStatusResponse = {
  cluster_status?: string
  total_instances?: number
  running_instances?: number
  uptime?: number
}

export type ApiResourcesResponse = {
  cpu_percent?: number
  memory?: {
    total_bytes?: number
    used_bytes?: number
    available_bytes?: number
    used_percent?: number
  }
  collected_at?: string
}

export type ApiAccount = {
  id: string
  player_id?: string
  enabled: boolean
  label?: string
  note?: string
  auth_status?: string
  has_token?: boolean
}

export type ApiAccountsResponse = {
  accounts?: ApiAccount[]
}

export type ApiInstance = {
  id: string
  account_id: string
  player_id?: string
  server_address?: string
  status?: string
  online_duration?: string
  last_seen?: string
  has_token?: boolean
  health?: number
  food?: number
  position?: {
    x?: number
    y?: number
    z?: number
  }
}

export type ApiInstancesResponse = {
  instances?: ApiInstance[]
}

export type ApiInitMicrosoftLoginResponse = ApiResponse & {
  user_code?: string
  verification_uri?: string
  verification_uri_complete?: string
  expires_in?: number
  interval?: number
  account_id?: string
}

export type ApiPollMicrosoftLoginResponse = ApiResponse & {
  status?: string
  account_id?: string
  minecraft_profile?: {
    id?: string
    name?: string
  }
}

export type ApiOperationLog = {
  id: string
  timestamp: string
  action: string
  target_instance_id?: string
  target_account_id?: string
  details?: string
  success: boolean
  error_msg?: string
  client_ip?: string
  user_agent?: string
}

export type ApiOperationLogsResponse = {
  logs?: ApiOperationLog[]
}

export type CreateAccountPayload = {
  id: string
  label?: string
  note?: string
}

export type CreateInstancePayload = {
  id: string
  account_id: string
  server_address: string
  enabled?: boolean
  auto_start?: boolean
}
