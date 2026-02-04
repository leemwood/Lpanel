document.addEventListener('DOMContentLoaded', () => {
    // 视图切换逻辑
    const navItems = document.querySelectorAll('.nav-item');
    const views = document.querySelectorAll('.view');
    
    navItems.forEach(item => {
        item.addEventListener('click', () => {
            const viewName = item.getAttribute('data-view');
            
            // 更新导航状态
            navItems.forEach(i => i.classList.remove('active'));
            item.classList.add('active');
            
            // 更新视图显示
            views.forEach(v => v.style.display = 'none');
            document.getElementById(`${viewName}View`).style.display = 'block';
            
            // 触发对应视图的加载
            if (viewName === 'dashboard') startSystemMonitoring();
            if (viewName === 'sites') fetchSites();
            if (viewName === 'files') fetchFiles('E:\\');
        });
    });

    // --- 仪表盘逻辑 ---
    let monitorInterval;
    const startSystemMonitoring = () => {
        if (monitorInterval) clearInterval(monitorInterval);
        const update = async () => {
            try {
                const res = await fetch('/api/system/status');
                const data = await res.json();
                
                document.getElementById('cpuBar').style.width = `${data.cpu_usage}%`;
                document.getElementById('cpuText').innerText = `${data.cpu_usage.toFixed(1)}%`;
                
                document.getElementById('memBar').style.width = `${data.mem_percent}%`;
                const memUsedGB = (data.mem_used / 1024 / 1024 / 1024).toFixed(1);
                const memTotalGB = (data.mem_total / 1024 / 1024 / 1024).toFixed(1);
                document.getElementById('memText').innerText = `${data.mem_percent.toFixed(1)}% (${memUsedGB}GB / ${memTotalGB}GB)`;
                
                document.getElementById('diskBar').style.width = `${data.disk_percent}%`;
                const diskUsedGB = (data.disk_used / 1024 / 1024 / 1024).toFixed(1);
                const diskTotalGB = (data.disk_total / 1024 / 1024 / 1024).toFixed(1);
                document.getElementById('diskText').innerText = `${data.disk_percent.toFixed(1)}% (${diskUsedGB}GB / ${diskTotalGB}GB)`;
            } catch (e) { console.error('监控更新失败', e); }
        };
        update();
        monitorInterval = setInterval(update, 3000);
    };

    // --- 网站管理逻辑 ---
    const siteFormModal = document.getElementById('siteFormModal');
    document.getElementById('showAddSiteBtn').onclick = () => siteFormModal.style.display = 'block';
    document.getElementById('closeModalBtn').onclick = () => siteFormModal.style.display = 'none';

    const fetchSites = async () => {
        const res = await fetch('/api/sites');
        const sites = await res.json();
        const tbody = document.getElementById('siteTableBody');
        tbody.innerHTML = sites.map(site => `
            <tr>
                <td>${site.name}</td>
                <td>${site.domain}</td>
                <td>${site.port}</td>
                <td><span class="status-${site.status}">${site.status === 'running' ? '运行中' : '已停止'}</span></td>
                <td>
                    <button class="btn-success" onclick="toggleSite(${site.ID})">${site.status === 'running' ? '停止' : '启动'}</button>
                    <button class="btn-danger" onclick="deleteSite(${site.ID})">删除</button>
                </td>
            </tr>
        `).join('');
    };

    document.getElementById('addSiteForm').onsubmit = async (e) => {
        e.preventDefault();
        const formData = {
            name: document.getElementById('name').value,
            domain: document.getElementById('domain').value,
            port: parseInt(document.getElementById('port').value),
            path: document.getElementById('path').value,
            proxy_pass: document.getElementById('proxyPass').value
        };
        await fetch('/api/sites', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(formData)
        });
        siteFormModal.style.display = 'none';
        fetchSites();
    };

    window.toggleSite = async (id) => {
        await fetch(`/api/sites/${id}/toggle`, {method: 'POST'});
        fetchSites();
    };

    window.deleteSite = async (id) => {
        if (confirm('确定删除？')) {
            await fetch(`/api/sites/${id}`, {method: 'DELETE'});
            fetchSites();
        }
    };

    // --- 文件管理逻辑 ---
    let currentPath = 'E:\\';

    const fetchFiles = async (path) => {
        try {
            currentPath = path;
            renderBreadcrumbs(path);
            
            const res = await fetch(`/api/files/list?path=${encodeURIComponent(path)}`);
            const files = await res.json();
            
            const tbody = document.getElementById('fileTableBody');
            if (!Array.isArray(files)) {
                tbody.innerHTML = `<tr><td colspan="4" style="color: red; text-align: center;">加载失败: ${files.error || '未知错误'}</td></tr>`;
                return;
            }

            // 添加“返回上一级”选项
            let html = '';
            if (path !== 'E:\\' && path !== 'E:') {
                const parentPath = path.substring(0, path.lastIndexOf('\\')) || 'E:\\';
                html += `
                    <tr class="parent-dir" onclick="fetchFiles('${parentPath.replace(/\\/g, '\\\\')}')">
                        <td colspan="4">回退到上一级 ...</td>
                    </tr>
                `;
            }

            html += files.map(file => `
                <tr>
                    <td>${file.is_dir ? '📁' : '📄'} ${file.name}</td>
                    <td>${file.is_dir ? '-' : (file.size / 1024).toFixed(1) + ' KB'}</td>
                    <td>${file.mod_time}</td>
                    <td>
                        ${file.is_dir ? `<button onclick="fetchFiles('${(path.endsWith('\\') ? path + file.name : path + '\\' + file.name).replace(/\\/g, '\\\\')}')">打开</button>` : ''}
                        <button class="btn-danger" onclick="deleteFile('${(path.endsWith('\\') ? path + file.name : path + '\\' + file.name).replace(/\\/g, '\\\\')}')">删除</button>
                    </td>
                </tr>
            `).join('');
            
            tbody.innerHTML = html;
        } catch (error) {
            console.error('获取文件列表失败:', error);
            document.getElementById('fileTableBody').innerHTML = `<tr><td colspan="4" style="color: red; text-align: center;">网络错误或服务器异常</td></tr>`;
        }
    };

    const renderBreadcrumbs = (path) => {
        const parts = path.split('\\').filter(p => p !== '');
        const breadcrumbContainer = document.getElementById('currentPath');
        breadcrumbContainer.innerHTML = '';
        
        let currentBuildPath = '';
        parts.forEach((part, index) => {
            currentBuildPath += (index === 0 ? part : '\\' + part);
            const span = document.createElement('span');
            span.className = 'breadcrumb-item';
            span.innerText = part;
            const targetPath = currentBuildPath + (index === 0 && part.endsWith(':') ? '\\' : '');
            span.onclick = () => fetchFiles(targetPath);
            breadcrumbContainer.appendChild(span);
            if (index < parts.length - 1) {
                const separator = document.createElement('span');
                separator.innerText = ' \\ ';
                breadcrumbContainer.appendChild(separator);
            }
        });
    };
    window.fetchFiles = fetchFiles;

    window.deleteFile = async (path) => {
        if (confirm('确定删除文件/文件夹？')) {
            await fetch(`/api/files?path=${encodeURIComponent(path)}`, {method: 'DELETE'});
            const parentPath = path.substring(0, path.lastIndexOf('\\')) || 'E:\\';
            fetchFiles(parentPath);
        }
    };

    // --- 目录选择器逻辑 ---
    const dirPickerModal = document.getElementById('dirPickerModal');
    const pickerList = document.getElementById('pickerList');
    const pickerCurrentPathDisplay = document.getElementById('pickerCurrentPath');
    let pickerCurrentPath = 'E:\\';
    let selectedDirPath = '';

    const fetchPickerFiles = async (path) => {
        try {
            pickerCurrentPath = path;
            pickerCurrentPathDisplay.innerText = path;
            selectedDirPath = path; // 默认选择当前目录

            const res = await fetch(`/api/files/list?path=${encodeURIComponent(path)}`);
            const files = await res.json();

            if (!Array.isArray(files)) {
                pickerList.innerHTML = `<li style="color: red;">加载失败: ${files.error || '未知错误'}</li>`;
                return;
            }

            let html = '';
            // 添加“返回上一级”
            if (path !== 'E:\\' && path !== 'E:') {
                const parentPath = path.substring(0, path.lastIndexOf('\\')) || 'E:\\';
                html += `
                    <li onclick="fetchPickerFiles('${parentPath.replace(/\\/g, '\\\\')}')">
                        <span>⬅️</span> <span>回退到上一级 ...</span>
                    </li>
                `;
            }

            // 只显示目录
            const dirs = files.filter(f => f.is_dir);
            html += dirs.map(dir => {
                const fullPath = (path.endsWith('\\') ? path + dir.name : path + '\\' + dir.name).replace(/\\/g, '\\\\');
                return `
                    <li onclick="handlePickerDirClick(event, '${fullPath}')">
                        <span>📁</span> <span>${dir.name}</span>
                    </li>
                `;
            }).join('');

            pickerList.innerHTML = html || '<li>(空目录)</li>';
        } catch (error) {
            console.error('获取目录列表失败:', error);
            pickerList.innerHTML = `<li style="color: red;">网络错误或服务器异常</li>`;
        }
    };

    window.handlePickerDirClick = (event, path) => {
        // 双击进入目录，单击选中
        if (event.detail === 2) {
            fetchPickerFiles(path);
        } else {
            selectedDirPath = path;
            // 更新选中状态
            const items = pickerList.querySelectorAll('li');
            items.forEach(item => item.classList.remove('selected'));
            event.currentTarget.classList.add('selected');
        }
    };

    document.getElementById('selectPathBtn').onclick = () => {
        dirPickerModal.style.display = 'block';
        fetchPickerFiles(pickerCurrentPath);
    };

    document.getElementById('confirmDirBtn').onclick = () => {
        document.getElementById('path').value = selectedDirPath;
        dirPickerModal.style.display = 'none';
    };

    document.getElementById('closeDirPickerBtn').onclick = () => {
        dirPickerModal.style.display = 'none';
    };

    // 初始加载
    startSystemMonitoring();
});
