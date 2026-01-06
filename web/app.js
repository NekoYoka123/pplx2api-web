"use strict";

const $ = (id) => document.getElementById(id);

const elements = {
  authKey: $("authKey"),
  authKeyToggle: $("authKeyToggle"),
  newApiKey: $("newApiKey"),
  newApiKeyToggle: $("newApiKeyToggle"),
  apiKeyHint: $("apiKeyHint"),
  sessionsList: $("sessionsList"),
  addSessionBtn: $("addSessionBtn"),
  clearSessionsBtn: $("clearSessionsBtn"),
  sessionsCount: $("sessionsCount"),
  sessionsWarning: $("sessionsWarning"),
  sessionRowTemplate: $("sessionRowTemplate"),
  defaultModel: $("defaultModel"),
  refreshModelsBtn: $("refreshModelsBtn"),
  forceModel: $("forceModel"),
  ignoreModelMonitoring: $("ignoreModelMonitoring"),
  rejectModelMismatch: $("rejectModelMismatch"),
  ignoreSearchResult: $("ignoreSearchResult"),
  searchResultCompatible: $("searchResultCompatible"),
  maxChatHistoryLength: $("maxChatHistoryLength"),
  promptForFile: $("promptForFile"),
  isIncognito: $("isIncognito"),
  noRolePrefix: $("noRolePrefix"),
  isMaxSubscribe: $("isMaxSubscribe"),
  proxy: $("proxy"),
  reloadBtn: $("reloadBtn"),
  saveBtn: $("saveBtn"),
  connStatus: $("connStatus"),
  serverAddress: $("serverAddress"),
  statusMessage: $("statusMessage"),
};

const storageKey = "pplx_admin_key";
elements.serverAddress.textContent = window.location.origin;
let modelList = [];
const buttonLabels = new Map();

function setConnection(label, state) {
  elements.connStatus.textContent = label;
  elements.connStatus.classList.remove("is-good", "is-bad", "is-warn");
  if (state) {
    elements.connStatus.classList.add(state);
  }
}

function setStatus(state, message) {
  elements.statusMessage.textContent = message;
  elements.statusMessage.dataset.state = state || "";
}

function setButtonLoading(button, isLoading, label) {
  if (!button) {
    return;
  }
  if (!buttonLabels.has(button)) {
    buttonLabels.set(button, button.textContent);
  }
  if (isLoading) {
    button.dataset.loading = "true";
    button.setAttribute("aria-busy", "true");
    button.disabled = true;
    if (label) {
      button.textContent = label;
    }
  } else {
    button.removeAttribute("data-loading");
    button.removeAttribute("aria-busy");
    button.disabled = false;
    const original = buttonLabels.get(button);
    if (original) {
      button.textContent = original;
    }
  }
}

function maskKey(key) {
  if (!key) return "未设置";
  if (key.length <= 4) return "****";
  return `****${key.slice(-4)}`;
}

function wireVisibilityToggle(button, input) {
  if (!button || !input) {
    return;
  }
  const updateLabel = (visible) => {
    button.textContent = visible ? "隐藏" : "显示";
    button.setAttribute("aria-pressed", visible ? "true" : "false");
  };
  updateLabel(input.type !== "password");
  button.addEventListener("click", () => {
    const shouldReveal = input.type === "password";
    input.type = shouldReveal ? "text" : "password";
    updateLabel(shouldReveal);
    input.focus({ preventScroll: true });
  });
}

function splitSessions(value) {
  return value
    .split(/[\n,]+/g)
    .map((item) => item.trim())
    .filter(Boolean);
}

function uniqueValues(values) {
  const seen = new Set();
  return values.filter((value) => {
    if (!value || seen.has(value)) {
      return false;
    }
    seen.add(value);
    return true;
  });
}

function getSessionInputs() {
  return Array.from(elements.sessionsList.querySelectorAll(".session-input"));
}

function getSessionValues() {
  return getSessionInputs()
    .map((input) => input.value.trim())
    .filter(Boolean);
}

function renumberSessionRows() {
  const rows = elements.sessionsList.querySelectorAll(".session-row");
  rows.forEach((row, index) => {
    const label = row.querySelector(".session-label");
    if (label) {
      label.textContent = `Cookie ${index + 1}`;
    }
  });
}

function updateSessionsMeta() {
  const sessions = getSessionValues();
  elements.sessionsCount.textContent = `${sessions.length}`;
  const warnings = [];
  if (sessions.length === 0) {
    warnings.push("至少需要一个会话。");
  }
  const duplicates = sessions.length - new Set(sessions).size;
  if (duplicates > 0) {
    warnings.push("发现重复 Cookie。");
  }
  elements.sessionsWarning.textContent = warnings.join(" ");
}

function handleSessionPaste(event) {
  const text = event.clipboardData ? event.clipboardData.getData("text") : "";
  if (!text) {
    return;
  }
  const parts = splitSessions(text);
  if (parts.length <= 1) {
    return;
  }
  const input = event.target;
  if (input.value.trim() !== "") {
    return;
  }
  event.preventDefault();
  input.value = parts[0];
  for (let i = 1; i < parts.length; i += 1) {
    addSessionRow(parts[i], { suppressMeta: true });
  }
  renumberSessionRows();
  updateSessionsMeta();
}

function addSessionRow(value, options = {}) {
  const row = elements.sessionRowTemplate.content.firstElementChild.cloneNode(true);
  const input = row.querySelector(".session-input");
  input.value = value || "";
  input.addEventListener("input", updateSessionsMeta);
  input.addEventListener("paste", handleSessionPaste);
  wireVisibilityToggle(row.querySelector(".toggle-visibility"), input);

  const removeBtn = row.querySelector(".remove-session");
  removeBtn.addEventListener("click", () => {
    row.remove();
    if (elements.sessionsList.children.length === 0) {
      addSessionRow("");
    } else {
      renumberSessionRows();
      updateSessionsMeta();
    }
  });

  elements.sessionsList.appendChild(row);
  renumberSessionRows();
  if (options.focus) {
    input.focus();
  }
  if (!options.suppressMeta) {
    updateSessionsMeta();
  }
}

function setSessionsList(values) {
  elements.sessionsList.innerHTML = "";
  if (!values || values.length === 0) {
    addSessionRow("", { focus: true });
    return;
  }
  values.forEach((value) => addSessionRow(value, { suppressMeta: true }));
  renumberSessionRows();
  updateSessionsMeta();
}

function setModelOptions(models, selected) {
  const select = elements.defaultModel;
  select.innerHTML = "";

  if (!models || models.length === 0) {
    const placeholder = new Option("未加载模型列表", "");
    placeholder.disabled = true;
    placeholder.selected = true;
    select.appendChild(placeholder);
    if (selected) {
      const fallback = new Option(`${selected}（当前）`, selected);
      select.appendChild(fallback);
      select.value = selected;
    }
    return;
  }

  const uniqueModels = uniqueValues(models);
  if (selected && !uniqueModels.includes(selected)) {
    select.appendChild(new Option(`${selected}（当前）`, selected));
  }
  uniqueModels.forEach((id) => {
    select.appendChild(new Option(id, id));
  });

  if (selected) {
    select.value = selected;
  }
  if (!select.value && uniqueModels.length > 0) {
    select.value = uniqueModels[0];
  }
}

async function loadModels(preferred, options = {}) {
  const key = elements.authKey.value.trim();
  if (!key) {
    if (!options.silent) {
      setStatus("error", "请输入管理密钥（config.json 的 apikey）后加载模型列表。");
    }
    return false;
  }

  const shouldAnnounce = !options.silent;
  if (shouldAnnounce) {
    setStatus("info", "正在加载模型列表...");
    setButtonLoading(elements.refreshModelsBtn, true, "加载中...");
  }

  try {
    const response = await fetch("/v1/models", {
      headers: {
        Authorization: `Bearer ${key}`,
      },
    });

    if (!response.ok) {
      const errorPayload = await response.json().catch(() => ({}));
      throw new Error(errorPayload.error || "加载模型列表失败。");
    }

    const data = await response.json();
    const list = Array.isArray(data.data) ? data.data.map((item) => item.id).filter(Boolean) : [];
    modelList = uniqueValues(list);
    setModelOptions(modelList, preferred);
    if (shouldAnnounce) {
      setStatus("success", "模型列表已更新。");
    }
    return true;
  } catch (err) {
    setModelOptions([], preferred);
    if (shouldAnnounce) {
      setStatus("error", err.message || "加载模型列表失败。");
    }
    return false;
  } finally {
    if (shouldAnnounce) {
      setButtonLoading(elements.refreshModelsBtn, false);
    }
  }
}

async function loadConfig() {
  const key = elements.authKey.value.trim();
  if (!key) {
    setStatus("error", "请输入管理密钥（config.json 的 apikey）后加载配置。");
    setConnection("未连接", "is-bad");
    return;
  }

  setButtonLoading(elements.reloadBtn, true, "加载中...");
  setButtonLoading(elements.saveBtn, true, "请稍候");
  setStatus("info", "正在加载配置...");
  try {
    const response = await fetch("/admin/config", {
      headers: {
        Authorization: `Bearer ${key}`,
      },
    });

    if (!response.ok) {
      const errorPayload = await response.json().catch(() => ({}));
      throw new Error(errorPayload.error || "加载配置失败。");
    }

    const data = await response.json();
    elements.apiKeyHint.textContent = data.api_key_hint || (data.api_key_set ? "已设置" : "未设置");
    elements.proxy.value = data.proxy || "";
    elements.isIncognito.checked = Boolean(data.is_incognito);
    elements.maxChatHistoryLength.value = data.max_chat_history_length || 1;
    elements.noRolePrefix.checked = Boolean(data.no_role_prefix);
    elements.searchResultCompatible.checked = Boolean(data.search_result_compatible);
    elements.promptForFile.value = data.prompt_for_file || "";
    elements.ignoreSearchResult.checked = Boolean(data.ignore_search_result);
    elements.ignoreModelMonitoring.checked = Boolean(data.ignore_model_monitoring);
    elements.rejectModelMismatch.checked = Boolean(data.reject_model_mismatch);
    elements.isMaxSubscribe.checked = Boolean(data.is_max_subscribe);
    const defaultModel = data.default_model || "claude-3.7-sonnet";
    elements.forceModel.value = data.force_model || "";
    setSessionsList(data.sessions || []);

    const modelsLoaded = await loadModels(defaultModel, { silent: true });
    if (!modelsLoaded) {
      setModelOptions(modelList, defaultModel);
    }

    setConnection("已连接", "is-good");
    setStatus("success", modelsLoaded ? "配置已加载。" : "配置已加载（模型列表未更新）。");
  } catch (err) {
    setConnection("未连接", "is-bad");
    setStatus("error", err.message || "加载配置失败。");
  } finally {
    setButtonLoading(elements.reloadBtn, false);
    setButtonLoading(elements.saveBtn, false);
  }
}

function buildPayload() {
  const maxChatHistoryLength = Number.parseInt(elements.maxChatHistoryLength.value, 10);
  if (!Number.isFinite(maxChatHistoryLength) || maxChatHistoryLength < 1) {
    setStatus("error", "最大聊天历史长度至少为 1。");
    return null;
  }

  const defaultModel = elements.defaultModel.value.trim();
  if (!defaultModel) {
    setStatus("error", "请选择默认模型（先刷新模型列表）。");
    return null;
  }

  const sessions = getSessionValues();
  if (sessions.length === 0) {
    setStatus("error", "至少需要一个会话 token。");
    return null;
  }

  const payload = {
    proxy: elements.proxy.value.trim(),
    is_incognito: elements.isIncognito.checked,
    max_chat_history_length: maxChatHistoryLength,
    no_role_prefix: elements.noRolePrefix.checked,
    search_result_compatible: elements.searchResultCompatible.checked,
    prompt_for_file: elements.promptForFile.value,
    ignore_search_result: elements.ignoreSearchResult.checked,
    ignore_model_monitoring: elements.ignoreModelMonitoring.checked,
    reject_model_mismatch: elements.rejectModelMismatch.checked,
    is_max_subscribe: elements.isMaxSubscribe.checked,
    default_model: defaultModel,
    force_model: elements.forceModel.value.trim(),
    sessions,
  };

  const newKey = elements.newApiKey.value.trim();
  if (newKey) {
    payload.apikey = newKey;
  }

  return payload;
}

async function saveConfig() {
  const key = elements.authKey.value.trim();
  if (!key) {
    setStatus("error", "请先输入管理密钥（config.json 的 apikey）再保存。");
    return;
  }

  const payload = buildPayload();
  if (!payload) {
    return;
  }

  setButtonLoading(elements.saveBtn, true, "保存中...");
  setStatus("info", "正在保存配置...");
  try {
    const response = await fetch("/admin/config", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${key}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      const errorPayload = await response.json().catch(() => ({}));
      throw new Error(errorPayload.error || "保存配置失败。");
    }

    const data = await response.json();
    if (payload.apikey) {
      elements.authKey.value = payload.apikey;
      localStorage.setItem(storageKey, payload.apikey);
      elements.newApiKey.value = "";
      elements.apiKeyHint.textContent = maskKey(payload.apikey);
    }

    setConnection("已连接", "is-good");
    const changed = Array.isArray(data.changed) && data.changed.length > 0 ? data.changed.join(", ") : "无变更字段";
    setStatus("success", `保存成功（${changed}）。`);
  } catch (err) {
    setConnection("未连接", "is-bad");
    setStatus("error", err.message || "保存配置失败。");
  } finally {
    setButtonLoading(elements.saveBtn, false);
  }
}

elements.reloadBtn.addEventListener("click", loadConfig);
elements.saveBtn.addEventListener("click", saveConfig);
elements.addSessionBtn.addEventListener("click", () => addSessionRow("", { focus: true }));
elements.clearSessionsBtn.addEventListener("click", () => setSessionsList([]));
elements.refreshModelsBtn.addEventListener("click", () => loadModels(elements.defaultModel.value, { silent: false }));
elements.authKey.addEventListener("input", (event) => {
  const value = event.target.value.trim();
  if (value) {
    localStorage.setItem(storageKey, value);
  } else {
    localStorage.removeItem(storageKey);
  }
});

const storedKey = localStorage.getItem(storageKey);
if (storedKey) {
  elements.authKey.value = storedKey;
}
wireVisibilityToggle(elements.authKeyToggle, elements.authKey);
wireVisibilityToggle(elements.newApiKeyToggle, elements.newApiKey);
setSessionsList([]);
setModelOptions([], "");
setConnection("未连接", "is-bad");
setStatus("info", "请输入管理密钥（config.json 的 apikey）后加载配置。");
