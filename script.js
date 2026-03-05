// ============================================================
// script.js — Todolist Frontend, terhubung ke Go REST API
// Base URL bisa diganti sesuai host server Go
// ============================================================

const API = "http://localhost:8080/api/tasks";

let tasks = [];
let currentFilter = "semua";

// ── Fetch all tasks from Go backend ──────────────────────────
async function loadTasks() {
  try {
    const res = await fetch(API);
    if (!res.ok) throw new Error("Gagal memuat tugas");
    tasks = await res.json();
    render();
  } catch (err) {
    showToast("⚠️ " + err.message, "error");
  }
}

// ── Add new task ──────────────────────────────────────────────
async function addTask() {
  const input = document.getElementById("task-input");
  const catSelect = document.getElementById("cat-select");
  const text = input.value.trim();
  const category = catSelect.value;

  if (!text) {
    showToast("Tulis dulu tugasnya!", "error");
    input.focus();
    return;
  }

  try {
    const res = await fetch(API, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ text, category }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Gagal menambah tugas");

    tasks.push(data);
    input.value = "";
    render();
    showToast("✅ Tugas ditambahkan!");
  } catch (err) {
    showToast("⚠️ " + err.message, "error");
  }
}

// ── Toggle done status ────────────────────────────────────────
async function toggleTask(id) {
  const task = tasks.find((t) => t.id === id);
  if (!task) return;

  try {
    const res = await fetch(`${API}/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ done: !task.done }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Gagal update tugas");

    const idx = tasks.findIndex((t) => t.id === id);
    tasks[idx] = data;
    render();
  } catch (err) {
    showToast("⚠️ " + err.message, "error");
  }
}

// ── Delete single task ────────────────────────────────────────
async function deleteTask(id) {
  try {
    const res = await fetch(`${API}/${id}`, { method: "DELETE" });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Gagal menghapus tugas");

    tasks = tasks.filter((t) => t.id !== id);
    render();
    showToast("🗑️ Tugas dihapus");
  } catch (err) {
    showToast("⚠️ " + err.message, "error");
  }
}

// ── Delete all done tasks ─────────────────────────────────────
async function clearDone() {
  const doneCount = tasks.filter((t) => t.done).length;
  if (doneCount === 0) {
    showToast("Tidak ada tugas selesai untuk dihapus", "error");
    return;
  }

  try {
    const res = await fetch(`${API}/done`, { method: "DELETE" });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Gagal menghapus");

    tasks = tasks.filter((t) => !t.done);
    render();
    showToast(`🗑️ ${data.deleted} tugas selesai dihapus`);
  } catch (err) {
    showToast("⚠️ " + err.message, "error");
  }
}

// ── Filter ────────────────────────────────────────────────────
function setFilter(filter, btn) {
  currentFilter = filter;
  document.querySelectorAll(".filter-btn").forEach((b) => b.classList.remove("active"));
  btn.classList.add("active");
  render();
}

// ── Render ────────────────────────────────────────────────────
function render() {
  const total = tasks.length;
  const done = tasks.filter((t) => t.done).length;
  const left = total - done;
  const pct = total === 0 ? 0 : Math.round((done / total) * 100);

  document.getElementById("total-count").textContent = total;
  document.getElementById("done-count").textContent = done;
  document.getElementById("left-count").textContent = left;
  document.getElementById("progress-pct").textContent = pct + "%";
  document.getElementById("progress-fill").style.width = pct + "%";

  let filtered = tasks;
  if (currentFilter === "aktif") filtered = tasks.filter((t) => !t.done);
  if (currentFilter === "selesai") filtered = tasks.filter((t) => t.done);

  const list = document.getElementById("task-list");
  const empty = document.getElementById("empty-state");

  if (filtered.length === 0) {
    list.innerHTML = "";
    empty.style.display = "flex";
    return;
  }

  empty.style.display = "none";
  list.innerHTML = filtered
    .map(
      (t) => `
    <div class="task-item ${t.done ? "done" : ""}" data-id="${t.id}">
      <button class="check-btn" onclick="toggleTask(${t.id})" title="${t.done ? 'Tandai belum selesai' : 'Tandai selesai'}">
        ${t.done ? "✓" : ""}
      </button>
      <div class="task-content">
        <span class="task-text">${escapeHtml(t.text)}</span>
        <span class="task-meta">
          <span class="task-cat cat-${t.category}">${t.category}</span>
          <span class="task-time">${t.created_at || ""}</span>
        </span>
      </div>
      <button class="delete-btn" onclick="deleteTask(${t.id})" title="Hapus tugas">✕</button>
    </div>
  `
    )
    .join("");
}

// ── Toast notification ────────────────────────────────────────
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

// ── Helpers ───────────────────────────────────────────────────
function escapeHtml(str) {
  return str
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}

// ── Enter key shortcut ────────────────────────────────────────
document.addEventListener("DOMContentLoaded", () => {
  document.getElementById("task-input").addEventListener("keydown", (e) => {
    if (e.key === "Enter") addTask();
  });
  loadTasks();
});
