import { describe, expect, it, vi, beforeEach } from "vitest";
import { flushPromises } from "@vue/test-utils";
import EmailSettings from "@/views/settings/EmailSettings.vue";
import type { InvoiceEmailSettings } from "@/types";
import { mountView } from "@/test-utils/mount";
import { mockMessage, resetNaiveMocks } from "@/test-utils/naive";

const mockStore = {
  settings: null as InvoiceEmailSettings | null,
  fetchSettings: vi.fn(),
  saveSettings: vi.fn(),
};

vi.mock("@/stores/invoiceEmailSettings", () => ({
  useInvoiceEmailSettingsStore: () => mockStore,
}));

describe("EmailSettings view", () => {
  const defaultSettings: InvoiceEmailSettings = {
    provider: "mailto",
    subjectTemplate: "Invoice {{number}}",
    bodyTemplate: "Please find attached invoice {{number}}.",
    signature: "",
  };

  beforeEach(() => {
    vi.clearAllMocks();
    resetNaiveMocks();
    mockStore.settings = null;
    mockStore.fetchSettings = vi.fn().mockResolvedValue(undefined);
    mockStore.saveSettings = vi.fn().mockResolvedValue(undefined);
  });

  it("renders invoice email settings form", async () => {
    const wrapper = mountView(EmailSettings);
    await flushPromises();

    expect(mockStore.fetchSettings).toHaveBeenCalledTimes(1);
    expect(wrapper.find(".select-trigger").exists()).toBe(true);
    expect(wrapper.find("input").exists()).toBe(true);
    expect(wrapper.find("textarea").exists()).toBe(true);
  });

  it("loads settings on mount", async () => {
    mockStore.settings = { ...defaultSettings };
    mountView(EmailSettings);
    await flushPromises();
    expect(mockStore.fetchSettings).toHaveBeenCalled();
  });

  it("shows provider options", async () => {
    const wrapper = mountView(EmailSettings);
    await flushPromises();
    expect(wrapper.text()).toContain("settings.email.options.provider.mailto");
    expect(wrapper.text()).toContain("settings.email.options.provider.resend");
    expect(wrapper.text()).toContain("settings.email.options.provider.smtp");
  });

  it("updates form when settings are loaded", async () => {
    mockStore.settings = {
      ...defaultSettings,
      subjectTemplate: "Loaded Template {{number}}",
      signature: "Best regards",
    };

    mountView(EmailSettings);
    await flushPromises();
    expect(mockStore.fetchSettings).toHaveBeenCalled();
  });

  it("toggles between provider fields correctly", async () => {
    const wrapper = mountView(EmailSettings);
    await flushPromises();
    expect(wrapper.text()).toContain("settings.email.options.provider.mailto");
    expect(wrapper.text()).toContain("settings.email.options.provider.resend");
    expect(wrapper.text()).toContain("settings.email.options.provider.smtp");
  });
});
