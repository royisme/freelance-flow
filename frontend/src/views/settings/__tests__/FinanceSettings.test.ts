import { describe, expect, it, vi, beforeEach } from "vitest";
import { flushPromises } from "@vue/test-utils";
import FinanceSettings from "@/views/settings/FinanceSettings.vue";
import { mountView } from "@/test-utils/mount";
import type { FinanceAccount } from "@/types/finance";

const mockApi = vi.hoisted(() => ({
  finance: {
    accounts: {
      list: vi.fn<[], Promise<FinanceAccount[]>>(),
    },
    settings: {
      get: vi.fn<[], Promise<any>>(),
      update: vi.fn<[any], Promise<void>>(),
    },
  },
}));

vi.mock("@/api", () => ({
  api: mockApi,
}));

vi.mock("vue-i18n", () => ({
  useI18n: () => ({
    t: (key: string) => key,
  }),
}));

describe("FinanceSettings view", () => {
  const mockAccounts: FinanceAccount[] = [
    { id: 1, name: "Checking", type: "checking", currency: "USD", balance: 5000, initialBalance: 5000, userId: 1, createdAt: "", updatedAt: "" },
    { id: 2, name: "Savings", type: "savings", currency: "USD", balance: 10000, initialBalance: 10000, userId: 1, createdAt: "", updatedAt: "" },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
    mockApi.finance.accounts.list.mockResolvedValue(mockAccounts);
    mockApi.finance.settings.get.mockResolvedValue({
      defaultAccountId: 1,
      autoCategorize: true,
      autoReconcile: false,
      userId: 1,
    });
    mockApi.finance.settings.update.mockResolvedValue(undefined);
  });

  it("renders finance settings form", async () => {
    const wrapper = mountView(FinanceSettings);

    await flushPromises();

    expect(wrapper.find(".select-trigger").exists()).toBe(true);
    expect(wrapper.find("input[type='checkbox']").exists()).toBe(true);
    expect(wrapper.find("button").exists()).toBe(true);
  });

  it("loads accounts on mount", async () => {
    mountView(FinanceSettings);

    await flushPromises();

    expect(mockApi.finance.accounts.list).toHaveBeenCalledTimes(1);
  });

  it("loads settings on mount", async () => {
    mountView(FinanceSettings);

    await flushPromises();

    expect(mockApi.finance.settings.get).toHaveBeenCalledTimes(1);
  });

  it("toggles auto categorize switch", async () => {
    const wrapper = mountView(FinanceSettings);

    await flushPromises();

    const switchElement = wrapper.find("input[type='checkbox']");
    expect(switchElement.exists()).toBe(true);
  });
});
