# Module Architecture Evolution

> **Status**: Discussion / Future Planning  
> **Last Updated**: 2024-12-16

This document discusses the architectural evolution path for FreelanceFlow's module system, considering future extensibility, plugin support, and commercial licensing models.

---

## 1. Current State

The current module system is documented in [MODULES.md](./MODULES.md). Key characteristics:

| Aspect        | Current Implementation                           |
| ------------- | ------------------------------------------------ |
| **Loading**   | Compile-time (static imports)                    |
| **Toggle**    | Runtime enable/disable via `moduleOverrides`     |
| **Packaging** | All modules bundled into main app                |
| **API**       | Deep integration, shared runtime                 |
| **Examples**  | Finance module (`frontend/src/modules/finance/`) |

This is essentially a **Built-in Module** or **Feature Flag** pattern, suitable for:

- Official first-party features
- Features requiring deep integration
- High stability requirements

---

## 2. Reference: Obsidian's Plugin Architecture

Obsidian is an excellent reference for desktop app plugin systems:

```
┌─────────────────────────────────────────┐
│            Obsidian Core                │  ← Free, core editor
├─────────────────────────────────────────┤
│         Core Plugins (built-in)         │  ← Toggleable, shipped with app
├─────────────────────────────────────────┤
│      Community Plugins (3rd party)      │  ← Runtime loaded, user installed
├─────────────────────────────────────────┤
│     Sync / Publish (paid services)      │  ← Subscription, server-dependent
└─────────────────────────────────────────┘
```

### Obsidian's Key Design Decisions

| Aspect             | Obsidian's Approach                        | Rationale                                          |
| ------------------ | ------------------------------------------ | -------------------------------------------------- |
| **Core vs Plugin** | Core features built-in, plugins enhance    | Stable base experience                             |
| **Loading**        | Runtime from `.obsidian/plugins/`          | No repackaging needed                              |
| **Plugin Format**  | `main.js` + `manifest.json` + `styles.css` | Standardized, easy to distribute                   |
| **API Boundary**   | Explicit `Plugin` base class with hooks    | Control what plugins can do                        |
| **Business Model** | Free core + paid hosting services          | Plugin ecosystem drives traffic, services monetize |

---

## 3. FreelanceFlow's Potential Layering

### Proposed Architecture Layers

```
┌────────────────────────────────────────────┐
│              Core App                       │
│   (Dashboard, Clients, Projects, Time...)   │  ← Always built-in
├────────────────────────────────────────────┤
│          Built-in Modules (toggleable)      │
│          (Finance, Advanced Reports...)     │  ← Compile-time, toggleable
├────────────────────────────────────────────┤
│          Premium Modules (licensed)         │
│          (Enterprise features...)           │  ← Compile-time, license-gated
├────────────────────────────────────────────┤
│          Runtime Plugins (future)           │
│          (3rd party integrations)           │  ← Runtime loaded
└────────────────────────────────────────────┘
```

### Layer Characteristics

| Layer                | Loading | Gating       | Use Case                            |
| -------------------- | ------- | ------------ | ----------------------------------- |
| **Core App**         | Static  | None         | Essential features every user needs |
| **Built-in Modules** | Static  | Feature flag | Optional official features          |
| **Premium Modules**  | Static  | License key  | Paid advanced features              |
| **Runtime Plugins**  | Dynamic | API contract | 3rd party extensions                |

---

## 4. Key Architectural Decisions

### 4.1 Do We Need Third-Party Plugins?

**Considerations:**

- Target user base size (freelancers, not developers generally)
- Developer community potential
- Cost of maintaining plugin API stability

**Freelance Tool Characteristics:**

- Users are freelancers (mostly non-developers)
- Requirements are relatively standardized (time, clients, invoices, reports)
- May be better suited for **official modules + paid premium features** than open plugins

### 4.2 Runtime Loading vs Compile-Time Bundling

| Approach         | Suitable For                       | Technical Challenges                                     |
| ---------------- | ---------------------------------- | -------------------------------------------------------- |
| **Compile-time** | Official modules, core features    | Simple, type-safe, no security concerns                  |
| **Runtime**      | 3rd party plugins, rapid iteration | Requires sandbox, version compat, signature verification |

**Why Obsidian Can Do Runtime Loading:**

- Electron environment, can `require()` arbitrary JS
- Well-defined Plugin API with types
- Active developer community

**FreelanceFlow (Wails/Tauri) Challenges:**

- Frontend is pure Vue, no Node.js runtime
- Would need `<script type="module">` or eval for external code
- Higher security risks

### 4.3 Business Model Mapping

| Model                  | Suitable Features                        | Examples                     |
| ---------------------- | ---------------------------------------- | ---------------------------- |
| **Free Core**          | Basic time tracking, client management   | Dashboard, Clients, Projects |
| **One-time Purchase**  | Advanced reports, bulk invoicing         | Advanced Reports Module      |
| **Subscription**       | Cloud sync, multi-device, collaboration  | Cloud Sync Service           |
| **Plugin Marketplace** | 3rd party integrations, custom templates | QuickBooks Integration       |

---

## 5. Recommended Evolution Path

### Phase 1: Built-in Modules (Current)

- Compile-time bundling
- `toggleable` flag for optional modules
- `isModuleEnabled()` checks at runtime
- **No license gating yet**

### Phase 2: Premium Modules (When Needed)

- Still compile-time bundling
- Add `LicenseService` to check license keys
- Extend `isModuleEnabled()` to check license status
- Module definition adds `requiresLicense?: string` field

```typescript
// Future extension to AppModule
type AppModule = {
  // ... existing fields
  requiresLicense?: "basic" | "pro" | "enterprise";
};
```

### Phase 3: Runtime Plugins (If Necessary)

Only implement if there's clear demand for 3rd party ecosystem:

- Dedicated `plugins/` directory
- Plugin manifest format (`manifest.json`)
- Plugin loader with security constraints
- Limited API exposure via `PluginContext`
- Optional: Plugin marketplace / distribution

---

## 6. Current Extension Points

The current architecture already provides these hooks for future expansion:

| Extension Point     | Location          | Purpose                        |
| ------------------- | ----------------- | ------------------------------ |
| `ModuleID` type     | `types.ts`        | Add new module identifiers     |
| `isModuleEnabled()` | `registry.ts`     | Add license/permission checks  |
| `settingsPages`     | Module definition | Modules contribute settings UI |
| `messages`          | Module definition | Modules contribute i18n        |
| `routes`            | Module definition | Modules contribute routes      |

No immediate changes needed—the architecture is ready for Phase 2 when required.

---

## 7. Open Questions

- [ ] What's the target pricing model? (Subscription vs one-time purchase)
- [ ] Which features should be premium vs free?
- [ ] Is there demand for 3rd party integrations (accounting software, etc.)?
- [ ] Cloud sync / multi-device support plans?

---

## References

- [Obsidian Developer Docs](https://docs.obsidian.md/Plugins/Getting+started/Build+a+plugin)
- [Electron Plugin Patterns](https://www.electronjs.org/docs/latest/)
- Current implementation: [MODULES.md](./MODULES.md)
