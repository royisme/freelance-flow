import { describe, expect, it, vi, beforeEach } from "vitest";
import { flushPromises } from "@vue/test-utils";
import Login from "@/views/Login.vue";
import { mountView } from "@/test-utils/mount";

const mockAuthStore = {
  usersList: [
    { id: 1, username: "user1", avatarUrl: "", settings: {} },
    { id: 2, username: "user2", avatarUrl: "", settings: {} },
  ],
  login: vi.fn<[{ username: string; password: string }], Promise<void>>().mockResolvedValue(undefined),
};

vi.mock("@/stores/auth", () => ({
  useAuthStore: () => mockAuthStore,
}));

const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
  currentRoute: { value: { path: "/login" } },
};

vi.mock("vue-router", () => ({
  useRouter: () => mockRouter,
}));

vi.mock("vue-i18n", () => ({
  useI18n: () => ({
    t: (key: string) => key,
  }),
}));

describe("Login view", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockAuthStore.usersList = [
      { id: 1, username: "user1", avatarUrl: "", settings: {} },
      { id: 2, username: "user2", avatarUrl: "", settings: {} },
    ];
  });

  it("renders user selection grid when users exist", async () => {
    const wrapper = mountView(Login, {
      global: { stubs: ["router-link"] },
    });

    await flushPromises();

    expect(wrapper.text()).toContain("auth.welcome");
    expect(wrapper.text()).toContain("auth.selectUser");
    expect(wrapper.text()).toContain("user1");
    expect(wrapper.text()).toContain("user2");
  });

  it("redirects to register if no users exist", async () => {
    mockAuthStore.usersList = [];

    mountView(Login, {
      global: { stubs: ["router-link"] },
    });

    await flushPromises();

    expect(mockRouter.replace).toHaveBeenCalledWith("/register");
  });

  it("selects a user and shows password input", async () => {
    const wrapper = mountView(Login, {
      global: { stubs: ["router-link"] },
    });

    await flushPromises();

    const vm = wrapper.vm as unknown as { selectedUserId: string };
    vm.selectedUserId = "1";
    await flushPromises();

    expect(wrapper.find('input[type="password"]').exists()).toBe(true);
    expect(wrapper.text()).toContain("auth.login");
  });

  it("goes to register page", async () => {
    const wrapper = mountView(Login, {
      global: { stubs: ["router-link"] },
    });

    await flushPromises();

    const buttons = wrapper.findAll("button");
    const registerButton = buttons.find((b) => b.text().includes("auth.createAccount"));
    await registerButton?.trigger("click");

    expect(mockRouter.push).toHaveBeenCalledWith("/register");
  });

});
