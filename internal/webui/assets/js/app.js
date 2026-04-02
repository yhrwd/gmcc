// GMCC Web Dashboard - Main Application
const API_BASE = '/api';

// 状态管理
let currentAction = null;
let currentInstance = null;
let authToken = null;
let refreshInterval = null;

// 初始化
async function init() {
    await loadStatus();
    await loadAccounts();
    
    // 设置自动刷新
    refreshInterval = setInterval(() => {
        loadStatus();
        loadAccounts();
    }, 5000);
}

// 加载集群状态
async function loadStatus() {
    try {
        const res = await fetch(`${API_BASE}/status`);
        if (!res.ok) throw new Error('Failed to load status');
        const data = await res.json();
        updateStatusUI(data);
    } catch (err) {
        console.error('加载状态失败:', err);
        showToast('加载状态失败', 'error');
    }
}

// 加载账号列表
async function loadAccounts() {
    try {
        const res = await fetch(`${API_BASE}/accounts`);
        if (!res.ok) throw new Error('Failed to load accounts');
        const data = await res.json();
        renderAccounts(data.accounts);
    } catch (err) {
        console.error('加载账号失败:', err);
        showToast('加载账号失败', 'error');
    }
}

// 更新状态UI
function updateStatusUI(data) {
    const statusBadge = document.getElementById('cluster-status-badge');
    const statusMap = {
        'running': { text: '运行中', class: 'bg-green-100 text-green-800' },
        'stopped': { text: '已停止', class: 'bg-gray-100 text-gray-800' },
        'partial': { text: '部分运行', class: 'bg-yellow-100 text-yellow-800' },
        'error': { text: '错误', class: 'bg-red-100 text-red-800' }
    };
    
    const status = statusMap[data.cluster_status] || statusMap['stopped'];
    statusBadge.className = `px-3 py-1 rounded-full text-sm font-medium ${status.class}`;
    statusBadge.textContent = status.text;
    
    // 更新统计数据
    document.getElementById('stat-total').textContent = data.total_instances || 0;
    document.getElementById('stat-running').textContent = data.running_instances || 0;
    document.getElementById('stat-stopped').textContent = (data.total_instances - data.running_instances) || 0;
    document.getElementById('stat-uptime').textContent = formatDuration(data.uptime);
}

// 渲染账号列表
function renderAccounts(accounts) {
    const container = document.getElementById('accounts-list');
    
    if (!accounts || accounts.length === 0) {
        container.innerHTML = `
            <div class="px-6 py-8 text-center text-gray-500">
                暂无账号，请点击上方"添加账号"按钮创建
            </div>
        `;
        return;
    }
    
    container.innerHTML = accounts.map(account => `
        <div class="px-6 py-4 flex items-center justify-between hover:bg-gray-50 fade-in">
            <div class="flex items-center space-x-4">
                <div class="flex-shrink-0">
                    <div class="h-10 w-10 rounded-full ${getStatusColor(account.status)} flex items-center justify-center">
                        <span class="text-white font-bold">${account.player_id ? account.player_id[0].toUpperCase() : '?'}</span>
                    </div>
                </div>
                <div>
                    <div class="text-sm font-medium text-gray-900">${account.player_id}</div>
                    <div class="text-sm text-gray-500">${account.server_address}</div>
                    ${account.has_token ? '<span class="text-xs text-blue-600">已存储Token</span>' : ''}
                </div>
            </div>
            <div class="flex items-center space-x-4">
                <span class="px-2 py-1 text-xs font-medium rounded-full ${getStatusBadgeClass(account.status)}">
                    ${getStatusText(account.status)}
                </span>
                <div class="flex space-x-2">
                    ${account.status === 'stopped' || account.status === 'error' ? `
                        <button onclick="startInstance('${account.id}')" 
                                class="px-3 py-1 bg-green-600 text-white text-sm rounded hover:bg-green-700 transition-colors">
                            启动
                        </button>
                    ` : `
                        <button onclick="stopInstance('${account.id}')" 
                                class="px-3 py-1 bg-red-600 text-white text-sm rounded hover:bg-red-700 transition-colors">
                            停止
                        </button>
                        <button onclick="restartInstance('${account.id}')" 
                                class="px-3 py-1 bg-yellow-600 text-white text-sm rounded hover:bg-yellow-700 transition-colors">
                            重启
                        </button>
                    `}
                    <button onclick="deleteInstance('${account.id}')" 
                            class="px-3 py-1 bg-gray-600 text-white text-sm rounded hover:bg-gray-700 transition-colors">
                        删除
                    </button>
                </div>
            </div>
        </div>
    `).join('');
}

// 启动实例
function startInstance(instanceId) {
    currentAction = 'start';
    currentInstance = instanceId;
    showPasswordModal();
}

// 停止实例
function stopInstance(instanceId) {
    currentAction = 'stop';
    currentInstance = instanceId;
    showPasswordModal();
}

// 重启实例
function restartInstance(instanceId) {
    currentAction = 'restart';
    currentInstance = instanceId;
    showPasswordModal();
}

// 删除实例
function deleteInstance(instanceId) {
    if (!confirm('确定要删除此实例吗？此操作不可恢复。')) {
        return;
    }
    currentAction = 'delete';
    currentInstance = instanceId;
    showPasswordModal();
}

// 显示密码Modal
function showPasswordModal() {
    document.getElementById('password-modal').classList.remove('hidden');
    document.getElementById('password-input').value = '';
    document.getElementById('password-error').classList.add('hidden');
    document.getElementById('password-input').focus();
}

// 关闭密码Modal
function closePasswordModal() {
    document.getElementById('password-modal').classList.add('hidden');
    currentAction = null;
    currentInstance = null;
}

// 提交密码
async function submitPassword() {
    const password = document.getElementById('password-input').value;
    if (!password) {
        showPasswordError('请输入密码');
        return;
    }
    
    const errorEl = document.getElementById('password-error');
    errorEl.classList.add('hidden');
    
    try {
        let endpoint;
        switch (currentAction) {
            case 'start':
                endpoint = `/instances/${currentInstance}/start`;
                break;
            case 'stop':
                endpoint = `/instances/${currentInstance}/stop`;
                break;
            case 'restart':
                endpoint = `/instances/${currentInstance}/restart`;
                break;
            case 'delete':
                endpoint = `/instances/${currentInstance}`;
                break;
            default:
                throw new Error('Unknown action');
        }
        
        const options = {
            method: currentAction === 'delete' ? 'DELETE' : 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ password })
        };
        
        const res = await fetch(`${API_BASE}${endpoint}`, options);
        const data = await res.json();
        
        if (data.success) {
            closePasswordModal();
            showToast(`${getActionText(currentAction)}成功`, 'success');
            loadAccounts();
        } else {
            showPasswordError(data.error || '操作失败');
        }
    } catch (err) {
        showPasswordError('网络错误，请重试');
        console.error(err);
    }
}

// 显示密码错误
function showPasswordError(message) {
    const errorEl = document.getElementById('password-error');
    errorEl.textContent = message;
    errorEl.classList.remove('hidden');
}

// 添加账号相关
function showAddAccountModal() {
    document.getElementById('add-account-modal').classList.remove('hidden');
    document.getElementById('account-id').value = '';
    document.getElementById('account-player-id').value = '';
    document.getElementById('account-server').value = '';
    document.getElementById('account-official-auth').checked = false;
    document.getElementById('account-password').value = '';
    document.getElementById('add-account-error').classList.add('hidden');
}

function closeAddAccountModal() {
    document.getElementById('add-account-modal').classList.add('hidden');
}

async function submitAddAccount() {
    const id = document.getElementById('account-id').value.trim();
    const playerId = document.getElementById('account-player-id').value.trim();
    const server = document.getElementById('account-server').value.trim();
    const useOfficialAuth = document.getElementById('account-official-auth').checked;
    const password = document.getElementById('account-password').value;
    
    if (!id || !playerId || !server || !password) {
        showAddAccountError('请填写所有必填项');
        return;
    }
    
    try {
        const res = await fetch(`${API_BASE}/accounts`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                password,
                id,
                player_id: playerId,
                server_address: server,
                use_official_auth: useOfficialAuth
            })
        });
        
        const data = await res.json();
        
        if (data.success) {
            closeAddAccountModal();
            showToast('账号添加成功', 'success');
            loadAccounts();
        } else {
            showAddAccountError(data.error || '添加失败');
        }
    } catch (err) {
        showAddAccountError('网络错误，请重试');
        console.error(err);
    }
}

function showAddAccountError(message) {
    const errorEl = document.getElementById('add-account-error');
    errorEl.textContent = message;
    errorEl.classList.remove('hidden');
}

// 工具函数
function getStatusColor(status) {
    const colors = {
        running: 'bg-green-500',
        stopped: 'bg-gray-400',
        error: 'bg-red-500',
        pending: 'bg-yellow-500',
        starting: 'bg-blue-500',
        reconnecting: 'bg-orange-500'
    };
    return colors[status] || 'bg-gray-400';
}

function getStatusBadgeClass(status) {
    const classes = {
        running: 'bg-green-100 text-green-800',
        stopped: 'bg-gray-100 text-gray-800',
        error: 'bg-red-100 text-red-800',
        pending: 'bg-yellow-100 text-yellow-800',
        starting: 'bg-blue-100 text-blue-800',
        reconnecting: 'bg-orange-100 text-orange-800'
    };
    return classes[status] || 'bg-gray-100 text-gray-800';
}

function getStatusText(status) {
    const texts = {
        running: '运行中',
        stopped: '已停止',
        error: '错误',
        pending: '启动中',
        starting: '正在启动',
        reconnecting: '重连中'
    };
    return texts[status] || status;
}

function getActionText(action) {
    const texts = {
        start: '启动',
        stop: '停止',
        restart: '重启',
        delete: '删除'
    };
    return texts[action] || action;
}

function formatDuration(duration) {
    if (!duration) return '-';
    
    // 解析Go duration字符串（如 "3h24m30s"）
    const match = duration.toString().match(/(\d+)h(\d+)m/);
    if (match) {
        return `${match[1]}小时${match[2]}分`;
    }
    
    // 如果是数字（纳秒），转换为可读格式
    const nanoseconds = parseInt(duration);
    if (!isNaN(nanoseconds)) {
        const seconds = Math.floor(nanoseconds / 1e9);
        const minutes = Math.floor(seconds / 60);
        const hours = Math.floor(minutes / 60);
        const days = Math.floor(hours / 24);
        
        if (days > 0) return `${days}天${hours % 24}小时`;
        if (hours > 0) return `${hours}小时${minutes % 60}分`;
        if (minutes > 0) return `${minutes}分${seconds % 60}秒`;
        return `${seconds}秒`;
    }
    
    return duration.toString();
}

// Toast通知
function showToast(message, type = 'success') {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    
    container.appendChild(toast);
    
    setTimeout(() => {
        toast.style.opacity = '0';
        toast.style.transform = 'translateX(100%)';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

// 键盘事件
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        closePasswordModal();
        closeAddAccountModal();
    }
});

// 初始化
init();
