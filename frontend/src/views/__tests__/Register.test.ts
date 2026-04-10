import { describe, expect, it, vi, beforeEach } from "vitest";
import { flushPromises } from "@vue/test-utils";
import Register from "@/views/Register.vue";
import { mountView } from "@/test-utils/mount";

const mockAuthStore = {
  register: vi.fn<[unknown], Promise<void>>(),
  usersList: [],
};

const mockAppStore = {
  setLocale: vi.fn(),
  theme: "light",
};

vi.mock("@/stores/auth", () => ({
  useAuthStore: () => mockAuthStore,
}));

vi.mock("@/stores/app", () => ({
  useAppStore: () => mockAppStore,
}));

const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
};

vi.mock("vue-router", () => ({
  useRouter: () => mockRouter,
}));

vi.mock("vue-i18n", () => ({
  useI18n: () => ({
    t: (key: string) => key,
    locale: { value: "en-US" },
  }),
}));

describe("Register view", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockAuthStore.register.mockResolvedValue({ id: 1 });
    mockAuthStore.usersList = [];
  });

  it("renders registration form", async () => {
    const wrapper = mountView(Register, { global: { stubs: ["router-link"] } });

    await flushPromises();

    expect(wrapper.text()).toContain("auth.stepProfile");
    expect(
      wrapper.find('input[placeholder*="auth.usernamePlaceholder"]').exists()
    ).toBe(true);
    expect(wrapper.findAll("input").length).toBeGreaterThanOrEqual(2);
  });

  it("shows avatar preview", async () => {
    const wrapper = mountView(Register, { global: { stubs: ["router-link"] } });

    await flushPromises();

    const hasAvatar =
      wrapper.findComponent({ name: "Avatar" }).exists() || wrapper.html().includes("api.dicebear.com");
    expect(hasAvatar).toBe(true);
  });

  it("does not advance when passwords mismatch", async () => {
    const wrapper = mountView(Register, { global: { stubs: ["router-link"] } });

    await flushPromises();

    const nextButton1 = wrapper
      .findAll("button")
      .find((b) => b.text().includes("common.next"));
    await nextButton1?.trigger("click");
    await flushPromises();

    expect(wrapper.text()).toContain("auth.setPassword");

    const passwordInput = wrapper.findAll('input[type="password"]').at(0);
    const confirmInput = wrapper.findAll('input[type="password"]').at(1);

    await passwordInput?.setValue("password123");
    await confirmInput?.setValue("differentpassword");

    await flushPromises();

    const nextButton2 = wrapper
      .findAll("button")
      .find((b) => b.text().includes("common.next"));
    await nextButton2?.trigger("click");
    await flushPromises();

    expect(wrapper.text()).toContain("auth.setPassword");
  });

  it("shows final submit button on step 3", async () => {
    const wrapper = mountView(Register, { global: { stubs: ["router-link"] } });

    await flushPromises();

    // Fill valid form
    await wrapper
      .find('input[placeholder*="auth.usernamePlaceholder"]')
      .setValue("validuser");
    await wrapper
      .find('input[placeholder*="auth.emailPlaceholder"]')
      .setValue("valid@example.com");
    await flushPromises();
    const nextButton1 = wrapper
      .findAll("button")
      .find((b) => b.text().includes("common.next"));
    await nextButton1?.trigger("click");
    await flushPromises();

    const passwordInputs = wrapper.findAll('input[type="password"]');
    await passwordInputs.at(0)?.setValue("password123");
    await passwordInputs.at(1)?.setValue("password123");

    const nextButton2 = wrapper
      .findAll("button")
      .find((b) => b.text().includes("common.next"));
    await nextButton2?.trigger("click");
    await flushPromises();

    expect(wrapper.text()).toContain("auth.financialPreferences");
  });

});
