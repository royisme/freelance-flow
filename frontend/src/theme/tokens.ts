export type ThemeMode = "light" | "dark";

type ThemeTokens = {
  // Primitives
  primitives: Record<string, string>;
  // Semantics
  semantics: Record<string, string>;
  // Misc (still semantic but not strictly colors)
  effects: Record<string, string>;
};

const lightTokens: ThemeTokens = {
  primitives: {
    // Neutral Scale (Sand & Stone)
    "--c-neutral-50": "#F8F7F4",
    "--c-neutral-100": "#F3F2EE",
    "--c-neutral-200": "#EBE9E4",
    "--c-neutral-300": "#D6D3D1",
    "--c-neutral-400": "#A8A29E",
    "--c-neutral-500": "#78716C",
    "--c-neutral-600": "#5F5C56",
    "--c-neutral-700": "#44403C",
    "--c-neutral-800": "#2C2B29",
    "--c-neutral-900": "#1F1E1C",

    // Neutral Scale (Dark references; defined for completeness)
    "--c-neutral-dark-50": "#E7E5E4",
    "--c-neutral-dark-100": "#D6D3D1",
    "--c-neutral-dark-200": "#A8A29E",
    "--c-neutral-dark-300": "#78716C",
    "--c-neutral-dark-600": "#57534E",
    "--c-neutral-dark-700": "#44403C",
    "--c-neutral-dark-800": "#292524",
    "--c-neutral-dark-900": "#1C1917",
    "--c-neutral-dark-950": "#0C0A09",

    // Brand
    "--c-primary-500": "#8D7B68",
    "--c-primary-400": "#A4907C",
    "--c-primary-600": "#5D5145",

    // Accent palette (from splash_bg.jpg)
    "--color-warm-orange": "#FB923C",
    "--color-warm-amber": "#F59E0B",
    "--color-warm-rose": "#FB7185",
    "--color-warm-pink": "#F472B6",
    "--color-warm-purple": "#818CF8",
    "--color-deep-purple": "#7C3AED",
    "--color-sky-blue": "#38BDF8",

    // Status
    "--color-success": "#10B981",
    "--color-success-bg": "#ECFDF5",
    "--color-warning": "#F59E0B",
    "--color-warning-bg": "#FFF7ED",
    "--color-error": "#EF4444",
    "--color-error-bg": "#FEF2F2",
    "--color-info": "#3B82F6",
    "--color-info-bg": "#EFF6FF",
  },
  semantics: {
    // Background hierarchy
    "--bg-app": "#F3F2EE",
    "--bg-surface": "#FFFFFF",
    "--bg-surface-muted": "#F8F7F4",
    "--bg-inset": "#F8F7F4",
    "--bg-overlay": "#FFFFFF",
    "--bg-card": "rgba(255, 255, 255, 0.85)",
    "--bg-card-solid": "#FFFFFF",

    // Text hierarchy
    "--text-primary": "#2C2B29",
    "--text-secondary": "#5F5C56",
    "--text-tertiary": "#78716C",
    "--text-on-gradient": "#FFFFFF",

    // Borders / dividers
    "--border-subtle": "#EBE9E4",
    "--border-default": "#D6D3D1",
    "--divider": "#EBE9E4",
    "--border-light": "rgba(255, 255, 255, 0.2)",
  },
  effects: {
    "--bg-hover": "#F8F7F4",
    "--bg-pressed": "#EBE9E4",
    "--focus-ring": "rgba(141, 123, 104, 0.35)",

    "--shadow-elevated":
      "0 2px 10px rgba(44, 43, 41, 0.08), 0 1px 4px rgba(44, 43, 41, 0.05)",

    "--gradient-warm":
      "linear-gradient(135deg, #FB923C 0%, #F472B6 50%, #818CF8 100%)",
    "--gradient-warm-subtle":
      "linear-gradient(135deg, rgba(251, 146, 60, 0.8) 0%, rgba(244, 114, 182, 0.8) 50%, rgba(129, 140, 248, 0.8) 100%)",
    "--gradient-glass-border":
      "linear-gradient(135deg, rgba(255,255,255,0.3) 0%, rgba(255,255,255,0.1) 100%)",
  },
};

const darkTokens: ThemeTokens = {
  primitives: {
    ...lightTokens.primitives,
    // Brand (Dark)
    "--c-primary-500": "#C8B6A6",
    "--c-primary-400": "#D9C8B8",
    "--c-primary-600": "#E7E5E4",
  },
  semantics: {
    // Background hierarchy
    "--bg-app": "#1C1917",
    "--bg-surface": "#292524",
    "--bg-surface-muted": "#201D1B",
    "--bg-inset": "#2A2725",
    "--bg-overlay": "#292524",
    "--bg-card": "rgba(41, 37, 36, 0.85)",
    "--bg-card-solid": "#292524",

    // Text hierarchy
    "--text-primary": "#E7E5E4",
    "--text-secondary": "#D6D3D1",
    "--text-tertiary": "#A8A29E",
    "--text-on-gradient": "#FFFFFF",

    // Borders / dividers
    "--border-subtle": "#44403C",
    "--border-default": "#57534E",
    "--divider": "#44403C",
    "--border-light": "rgba(255, 255, 255, 0.12)",
  },
  effects: {
    "--bg-hover": "#2A2725",
    "--bg-pressed": "#44403C",
    "--focus-ring": "rgba(200, 182, 166, 0.45)",
    "--shadow-elevated":
      "0 10px 30px rgba(0, 0, 0, 0.35), 0 4px 12px rgba(0, 0, 0, 0.25)",
    "--gradient-warm":
      "linear-gradient(135deg, #FB923C 0%, #F472B6 50%, #818CF8 100%)",
    "--gradient-warm-subtle":
      "linear-gradient(135deg, rgba(251, 146, 60, 0.8) 0%, rgba(244, 114, 182, 0.8) 50%, rgba(129, 140, 248, 0.8) 100%)",
    "--gradient-glass-border":
      "linear-gradient(135deg, rgba(255,255,255,0.12) 0%, rgba(255,255,255,0.06) 100%)",
  },
};

export const THEME_TOKENS: Record<ThemeMode, ThemeTokens> = {
  light: lightTokens,
  dark: darkTokens,
};

export function getCssVarsForMode(mode: ThemeMode): Record<string, string> {
  const tokens = THEME_TOKENS[mode];
  return {
    ...tokens.primitives,
    ...tokens.semantics,
    ...tokens.effects,
  };
}

export function applyThemeToRoot(mode: ThemeMode): void {
  if (typeof document === "undefined") return;
  const root = document.documentElement;
  root.dataset.theme = mode;
  root.classList.toggle("dark", mode === "dark");
  const vars = getCssVarsForMode(mode);
  for (const [name, value] of Object.entries(vars)) {
    root.style.setProperty(name, value);
  }
}
