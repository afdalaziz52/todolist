const API = "http://localhost:8080/api";

let tasks = [];
let currentFilter = "semua";

// ambil token dari localStorage
function getToken() {
    const token = localStorage.getItem("token");
    if (!token) {
        window.location.href = "login.html";
        return null;
    }
    return token;
}

// ── Load tasks ──────────────────────────────────────────────
async function loadTasks() {
    const token = getToken();
    if (!token) return;

    try {
        const res = await fetch(`${API}/tasks`, {
            headers: { "Authorization": `Bearer ${token}` }
        });

        if (res.status === 401) {
            localStorage.removeItem("token");
            window.location.href = "login.html";
            return;
        }

        const data = await res.json();
        tasks = data.data || [];
        render();
    } catch (err) {
        showToast("⚠️ " + err.message, "error");
    }
}

// ── Add task ────────────────────────────────────────────────
async function addTask() {
    const token = getToken();
    if (!token) return;

    const input = document.getElementById("task-input");
    const catSelect = document.getElementById("cat-select");
    const title = input.value.trim();
    const category = catSelect.value;

    if (!title) {
        showToast("Tulis dulu tugasnya!", "error");
        input.focus();
        return;
    }

    try {
        const body = { title, category };

        // kalau other, tambah custom_category
        if (category === "other") {
            const custom = prompt("Masukkan nama kategori:");
            if (!custom) return;
            body.custom_category = custom;
        }

        const res = await fetch(`${API}/tasks`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify(body)
        });

        const data = await res.json();
        if (!res.ok) throw new Error(data.message || "Gagal menambah tugas");

        input.value = "";
        await loadTasks();
        showToast("✅ Tugas ditambahkan!");
    } catch (err) {
        showToast("⚠️ " + err.message, "error");
    }
}

// ── Toggle status ───────────────────────────────────────────
async function toggleTask(id, currentStatus) {
    const token = getToken();
    if (!token) return;

    const newStatus = currentStatus === "pending" ? "done" : "pending";

    try {
        const res = await fetch(`${API}/tasks/${id}/status`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({ status: newStatus })
        });

        if (!res.ok) throw new Error("Gagal update status");
        await loadTasks();
    } catch (err) {
        showToast("⚠️ " + err.message, "error");
    }
}

// ── Delete task ─────────────────────────────────────────────
async function deleteTask(id) {
    const token = getToken();
    if (!token) return;

    try {
        const res = await fetch(`${API}/tasks/${id}`, {
            method: "DELETE",
            headers: { "Authorization": `Bearer ${token}` }
        });

        if (!res.ok) throw new Error("Gagal menghapus tugas");
        tasks = tasks.filter(t => t.id !== id);
        render();
        showToast("🗑️ Tugas dihapus");
    } catch (err) {
        showToast("⚠️ " + err.message, "error");
    }
}

// ── Clear done tasks ────────────────────────────────────────
async function clearDone() {
    const doneTasks = tasks.filter(t => t.status === "done");
    if (doneTasks.length === 0) {
        showToast("Tidak ada tugas selesai untuk dihapus", "error");
        return;
    }

    try {
        // hapus satu-satu karena tidak ada endpoint bulk delete
        await Promise.all(doneTasks.map(t => deleteTask(t.id)));
        showToast(`🗑️ ${doneTasks.length} tugas selesai dihapus`);
    } catch (err) {
        showToast("⚠️ " + err.message, "error");
    }
}

// ── Filter ──────────────────────────────────────────────────
function setFilter(filter, btn) {
    currentFilter = filter;
    document.querySelectorAll(".filter-btn").forEach(b => b.classList.remove("active"));
    btn.classList.add("active");
    render();
}

// ── Render ──────────────────────────────────────────────────
function render() {
    const total = tasks.length;
    const done = tasks.filter(t => t.status === "done").length;
    const left = total - done;
    const pct = total === 0 ? 0 : Math.round((done / total) * 100);

    document.getElementById("total-count").textContent = total;
    document.getElementById("done-count").textContent = done;
    document.getElementById("left-count").textContent = left;
    document.getElementById("progress-pct").textContent = pct + "%";
    document.getElementById("progress-fill").style.width = pct + "%";

    let filtered = tasks;
    if (currentFilter === "aktif") filtered = tasks.filter(t => t.status === "pending");
    if (currentFilter === "selesai") filtered = tasks.filter(t => t.status === "done");

    const list = document.getElementById("task-list");
    const empty = document.getElementById("empty-state");

    if (filtered.length === 0) {
        list.innerHTML = "";
        empty.style.display = "flex";
        return;
    }

    empty.style.display = "none";
    list.innerHTML = filtered.map(t => `
        <div class="task-item ${t.status === "done" ? "done" : ""}" data-id="${t.id}">
            <button class="check-btn" onclick="toggleTask(${t.id}, '${t.status}')" 
                title="${t.status === "done" ? 'Tandai belum selesai' : 'Tandai selesai'}">
                ${t.status === "done" ? "✓" : ""}
            </button>
            <div class="task-content">
                <span class="task-text">${escapeHtml(t.title)}</span>
                <span class="task-meta">
                    <span class="task-cat cat-${t.category}">
                        ${t.custom_category ? t.custom_category : t.category}
                    </span>
                </span>
            </div>
            <button class="delete-btn" onclick="deleteTask(${t.id})" title="Hapus tugas">✕</button>
        </div>
    `).join("");
}

// ── Toast ───────────────────────────────────────────────────
function showToast(msg, type = "success") {
    let toast = document.getElementById("toast");
    if (!toast) {
        toast = document.createElement("div");
        toast.id = "toast";
        document.body.appendChild(toast);
    }
    toast.textContent = msg;
    toast.className = "toast show " + (type === "error" ? "toast-error" : "toast-success");
    clearTimeout(toast._timer);
    toast._timer = setTimeout(() => toast.classList.remove("show"), 2800);
}

// ── Helpers ─────────────────────────────────────────────────
function escapeHtml(str) {
    return str
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;");
}

// ── Init ────────────────────────────────────────────────────
document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("task-input").addEventListener("keydown", (e) => {
        if (e.key === "Enter") addTask();
    });
    loadTasks();
});