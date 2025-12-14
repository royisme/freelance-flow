# Modules System (Frontend)

This project treats each main sidebar entry as a **module**. A module owns:

- Navigation entry (sidebar)
- Routes (and route guards)
- Optional Settings sub-page(s)
- Optional i18n messages
- Optional enable/disable toggle (per user)

The goal is to make adding new functionality (e.g. **Tax Prep**) predictable and low-friction.

## 1. Key Files

- Module registry: `frontend/src/modules/registry.ts`
- Module types: `frontend/src/modules/types.ts`
- Module i18n merging: `frontend/src/locales/index.ts`
- Sidebar rendering: `frontend/src/layouts/MainLayout.vue`
- Router wiring + module guard: `frontend/src/router/index.ts`
- Module toggles UI: `frontend/src/views/settings/GeneralSettings.vue`

## 1.1 Copal Memory Shortcut (Codex CLI)

This decision is also stored in CoPal memory.

- Memory: `copal memory show 24635ac8-3c31-4d11-8602-5bd78e92ebc8`
- Codex prompt: `.codex/prompts/copal/memory-modules.md`

## 2. Data Model (Persistence)

Module enable/disable state is stored in `users.settings_json` as part of `UserSettings`:

- Field: `moduleOverrides: Record<string, boolean>`
  - Key: module id (e.g. `"finance"`)
  - Value: `true` forces enabled, `false` forces disabled
  - When a key is absent, the module falls back to `enabledByDefault` in the module definition

Backend types:

- `internal/dto/settings.go`
- `internal/models/user.go`
- `internal/services/settings_service.go` (normalization)

Wails bindings are regenerated during build; do not edit `frontend/src/wailsjs/**` manually.

## 3. Module Definition

Each module is an `AppModule`:

- `id`: stable module identifier
- `enabledByDefault`: whether it shows up without overrides
- `toggleable`: whether it can be enabled/disabled by the user
- `nav`: sidebar menu entry (supports nested children)
- `routes`: Vue Router route records for the module (must include `meta.moduleID`)
- `settingsPages` (optional): pages that should appear under Settings
- `messages` (optional): i18n messages owned by the module

Reference implementation: Finance module

- `frontend/src/modules/finance/module.ts`
- `frontend/src/modules/finance/i18n.ts`

## 4. Adding a New Module (Checklist)

Example: add `tax` module.

1. Create module folder:
   - `frontend/src/modules/tax/module.ts`
   - `frontend/src/modules/tax/i18n.ts` (optional)
2. Implement `taxModule: AppModule`:
   - Provide `routes` under `/tax` (or any prefix)
   - Set `meta: { requiresAuth: true, layout: "main", moduleID: "tax" }` on the route record (and any child routes)
   - Provide `nav` entry:
     - `labelKey` should live in your module i18n messages
   - If you want a settings page:
     - add `settingsPages: [{ key: "tax", labelKey: "settings.tax.title", component: () => import(...), order: 60, moduleID: "tax" }]`
3. Register the module:
   - Add it into `nonSettingsModules` in `frontend/src/modules/registry.ts`
   - Settings pages are aggregated automatically (Settings module is generated from all modules’ `settingsPages`)
4. Add i18n messages (optional but recommended):
   - Put your module messages in `frontend/src/modules/tax/i18n.ts`
   - Export them as `messages` from the module definition
5. If you want the module toggleable:
   - set `toggleable: true`
   - set your desired `enabledByDefault`
   - It will appear in Settings → General → Modules automatically
6. Run checks:
   - `cd frontend && bun run check`

## 5. Enable/Disable Behavior

When a module is disabled:

- It is removed from the sidebar menu (including its Settings sub-page entry)
- Navigation to its routes is blocked by a router guard and redirects to `/dashboard`

Notes:

- Settings itself is not toggleable.
- “Disabled” is per-user because it’s stored in `settings_json` (multi-user safe).

## 6. i18n Ownership Rules

- Base app i18n lives in:
  - `frontend/src/locales/en-US.ts`
  - `frontend/src/locales/zh-CN.ts`
- Module i18n lives inside the module folder and is merged into the global messages in:
  - `frontend/src/locales/index.ts`

Recommended convention:

- Navigation keys under `nav.*` (e.g. `nav.tax`)
- Module-specific keys under a dedicated namespace (e.g. `tax.*`)
- Settings page keys under `settings.<module>.*` (e.g. `settings.tax.title`)

## 7. Route Meta Requirements

To participate in module gating, routes must set:

- `meta.moduleID = "<module id>"`

If you forget this, the route will not be blocked when disabled.
