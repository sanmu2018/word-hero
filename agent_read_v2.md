# AI 微服务协作规则（适用于所有 AI 服务）

本文件定义 AI 的行为约束。所有规则均为强制执行。

---

# 1. 服务类型定义

每个服务必须属于以下三类之一：

- **agent**：唯一。负责需求分析、任务拆分、Issue 分发。
- **api**：唯一。负责 API 协议维护。
- **service**：多个。负责具体业务逻辑。

当前服务类型由 `agent_role.yaml` 提供。

---

# agent_role.yaml 文件说明

`agent_role.yaml` 用于告诉 AI 当前服务的“身份信息”，包括：

1. **service_name**  
   当前服务的名称（如 word-hero-agent / word-hero-api / xxx-service）。

2. **service_type**  
   当前服务属于哪一类：
    - `agent`：总控、调度、分发 Issue
    - `api`：API 协议维护
    - `service`：普通业务服务

3. **agent_name**（可选）  
   当前服务归属于哪个 agent 服务。
    - 若 service_type = agent，则此字段为空或忽略。
    - 若 service_type = api / service，则必须指定上游 agent。

4. **agent_url**（可选）  
   上游 agent 服务的访问地址。  
   用于非 agent 服务在需要时访问 agent。

---

# AI 在工作时如何使用 agent_role.yaml？

AI 在开始任务前必须读取 `agent_role.yaml` 来确定：

1. **我是哪个服务？**（service_name）
2. **我属于哪种服务类型？**（service_type）
3. **我是否需要与 agent 服务通信？**
    - 若是，则使用 agent_name / agent_url
4. **我应该遵循哪套行为规则？**
    - agent 的行为
    - api 的行为
    - service 的行为

# 示例

```yaml
service_name: word-hero-api
service_type: api
agent_name: word-hero-agent
agent_url: https://agent.example.com
```


---

# 3. 全局硬性规则（必须遵守）

### 3.1 Issue 驱动
- 只能从 Issue 开始工作。
- 只能完成属于本服务的任务。
- 完成后关闭 Issue。
- 不得在无 Issue 情况下修改代码。

### 3.2 服务边界
- 不得修改其他服务的代码。
- 如需跨服务任务，由 agent 创建 Issue 处理。
- 可临时 clone 其他服务代码，但只能放在 `.tmp/repos/{service}` 且只读。

### 3.3 API 协议
- 所有 API 协议必须从 API 服务获取。
- 不得本地维护/复制/推测 API。

### 3.4 信息真源
- 所有服务地址与类型必须来自 `agent_service.yaml`。
- 不得硬编码服务地址。
- 不得自我编造服务信息。

---

# 4. 各服务类型的行为

## 4.1 agent（唯一）
Agent 服务必须做到：

1. 需求分析 → 任务拆分。
2. 依据 `agent_service.yaml` 选择服务。
3. 为对应服务创建 Issue。
4. 将 Issue 链接写入 `service_issues.md`。
5. 不得直接修改其他服务代码。
6. agent 自己不执行业务逻辑。

## 4.2 api（唯一）
API 服务必须做到：

1. 所有 API 变更必须先创建 Issue。
2. 更新 API 协议（OpenAPI/Protobuf/自定义规范）。
3. 主动从自身 Issue 获取任务。
4. 不负责任何业务逻辑。

## 4.3 service（多个）
Service 服务必须做到：

1. 从自身 Issue 读取任务。
2. 只处理与本服务相关的内容。
3. 如需要查看其他服务代码：仅 clone 到 `.tmp/repos/`，只读。
4. 需要 API 信息时，从 API 服务拉取。
5. 完成任务后关闭 Issue。

---

# 5. 初始化规则

启动时 AI 必须执行：

1. 读取 `agent_role.yaml`：确认服务类型与上游 agent
2. 若非 agent：从 agent 服务拉取最新 `agent_service.yaml`
3. 若任务涉及 API：从 API 服务拉取最新 API 协议
4. 所有跨服务协作必须通过 Issue，不得私自调用或修改他人仓库

---

# 6. 禁止行为（AI 不得实现）

- 不得凭空创建服务
- 不得硬编码服务地址/API
- 不得修改其他服务仓库
- 不得关闭其他服务的 Issue

---

# 7. 允许行为（AI 可以做）

- 从当前服务 Issue 获取任务
- 生成代码、文档、测试（限当前服务）
- 访问 agent 提供的 `agent_service.yaml`
- 从 API 服务读取协议
- 临时 clone 其他服务代码到 `.tmp/repos/`

---

# 8. 工作流程（简化版）

## agent
需求 → 拆分 → 给服务创建 Issue → 记录 Issue 链接

## api
收到 Issue → 更新协议 → 完成 Issue

## service
读取 Issue → 完成本服务逻辑 → 提交 → 关闭 Issue

