import { describe, expect, it, vi, beforeEach } from "vitest";
import { flushPromises } from "@vue/test-utils";
import GeneralSettings from "@/views/settings/GeneralSettings.vue";
import { mountView } from "@/test-utils/mount";
const mockPreferencesStore = {
  preferences: {
    currency: "USD",
    language: "en-US",
    theme: "light",
    dateFormat: "2006-01-02",
    timezone: "UTC",
    moduleOverrides: {},
  },
  fetchPreferences: vi.fn(),
  savePreferences: vi.fn(),
};
const mockAppStore = {
  theme: "light",
  setTheme: vi.fn(),
  setLocale: vi.fn(),
};

vi.mock("@/stores/userPreferences", () => ({
  useUserPreferencesStore: () => mockPreferencesStore,
}));

vi.mock("@/stores/app", () => ({
  useAppStore: () => mockAppStore,
}));

describe("GeneralSettings view", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockPreferencesStore.fetchPreferences.mockResolvedValue(undefined);
    mockPreferencesStore.savePreferences.mockResolvedValue(undefined);
  });

  it("loads preferences on mount", async () => {
    const wrapper = mountView(GeneralSettings);

    await flushPromises();
    expect(mockPreferencesStore.fetchPreferences).toHaveBeenCalledTimes(1);
    expect(wrapper.text()).toContain("settings.general.cardTitle");
  });
});
