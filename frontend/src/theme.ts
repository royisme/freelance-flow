import type { GlobalThemeOverrides } from "naive-ui";
import type { ThemeMode } from "@/theme/tokens";
import { THEME_TOKENS } from "@/theme/tokens";

/**
 * FreelanceFlow Theme Overrides
 *
 * Architecture:
 * - Keep a semantic token layer, but Naive UI tokens must be literal colors.
 *   (Naive UI uses seemly to parse colors; CSS var(...) cannot be parsed.)
 * - Light/Dark switching is controlled by AppStore + Naive `darkTheme`.
 * - CSS variables are applied from the same TS tokens for non-Naive areas.
 */

// 1. Shared Configuration (Geometry, Typography) - "The Skeleton"
const sharedConfig = {
  fontFamily: '"Inter", "Varela Round", sans-serif',
  borderRadius: {
    base: "12px",
    card: "16px",
    button: "8px",
  },
  typography: {
    cardTitleWeight: "700",
    cardTitleSize: "18px",
  },
};

type SemanticPalette = {
  primary: string;
  primaryHover: string;
  primaryPressed: string;
  bgApp: string;
  bgSurface: string;
  bgInset: string;
  bgOverlay: string;
  text1: string;
  text2: string;
  text3: string;
  borderSubtle: string;
  borderDefault: string;
  divider: string;
  hover: string;
  pressed: string;
  shadow: string;
};

function paletteFromTokens(mode: ThemeMode): SemanticPalette {
  const t = THEME_TOKENS[mode];
  return {
    primary: t.primitives["--c-primary-500"],
    primaryHover: t.primitives["--c-primary-400"],
    primaryPressed: t.primitives["--c-primary-600"],

    bgApp: t.semantics["--bg-app"],
    bgSurface: t.semantics["--bg-surface"],
    bgInset: t.semantics["--bg-inset"],
    bgOverlay: t.semantics["--bg-overlay"],

    text1: t.semantics["--text-primary"],
    text2: t.semantics["--text-secondary"],
    text3: t.semantics["--text-tertiary"],

    borderSubtle: t.semantics["--border-subtle"],
    borderDefault: t.semantics["--border-default"],
    divider: t.semantics["--divider"],

    hover: t.effects["--bg-hover"],
    pressed: t.effects["--bg-pressed"],
    shadow: t.effects["--shadow-elevated"],
  };
}

function createTheme(mode: ThemeMode): GlobalThemeOverrides {
  const p = paletteFromTokens(mode);
  return {
    common: {
      primaryColor: p.primary,
      primaryColorHover: p.primaryHover,
      primaryColorPressed: p.primaryPressed,
      primaryColorSuppl: p.primary,

      bodyColor: p.bgApp,
      cardColor: p.bgSurface,
      modalColor: p.bgOverlay,
      popoverColor: p.bgOverlay,
      tableColor: p.bgSurface,
      inputColor: p.bgInset,

      textColor1: p.text1,
      textColor2: p.text2,
      textColor3: p.text3,

      borderColor: p.borderDefault,
      dividerColor: p.divider,

      hoverColor: p.hover,
      pressedColor: p.pressed,

      fontFamily: sharedConfig.fontFamily,
      borderRadius: sharedConfig.borderRadius.base,
    },
    Layout: {
      color: p.bgApp,
      siderColor: p.bgSurface,
      headerColor: p.bgSurface,
      footerColor: p.bgSurface,
      siderBorderColor: p.borderSubtle,
      headerBorderColor: p.borderSubtle,
      footerBorderColor: p.borderSubtle,
    },
    Card: {
      borderRadius: sharedConfig.borderRadius.card,
      borderColor: p.borderSubtle,
      boxShadow: p.shadow,
      titleTextColor: p.text1,
      titleFontWeight: sharedConfig.typography.cardTitleWeight,
    },
    Button: {
      borderRadius: sharedConfig.borderRadius.button,
      fontWeight: "600",
      textColorQuaternary: p.text2,
      textColorHoverQuaternary: p.primary,
    },
    Menu: {
      itemBorderRadius: "8px",
      fontSize: "14px",
      itemHeight: "40px",

      itemColorHover: p.hover,
      itemTextColorHover: p.primary,
      itemIconColorHover: p.primary,

      itemColorActive: p.pressed,
      itemTextColorActive: p.primary,
      itemIconColorActive: p.primary,
    },
    Statistic: {
      labelTextColor: p.text2,
      labelFontWeight: "500",
      valueTextColor: p.text1,
      valueFontWeight: "700",
    },
  };
}

export const lightThemeOverrides = createTheme("light");
export const darkThemeOverrides = createTheme("dark");
