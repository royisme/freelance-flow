import { describe, expect, it, vi, beforeEach } from "vitest";
import { flushPromises } from "@vue/test-utils";
import FinanceOverview from "@/views/finance/index.vue";
import { mountView } from "@/test-utils/mount";
import type { FinanceSummary } from "@/types/finance";

const mockApi = vi.hoisted(() => ({
  finance: {
    summary: {
      get: vi.fn<[], Promise<FinanceSummary>>(),
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

describe("Finance view", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockApi.finance.summary.get.mockResolvedValue({
      totalBalance: 5000,
      totalIncome: 1000,
      totalExpense: 200,
      cashFlow: 800,
      byAccount: [],
      monthlyTrends: [],
    });
  });

  it("renders finance overview page", async () => {
    const wrapper = mountView(FinanceOverview);

    await flushPromises();

    expect(wrapper.find(".finance-overview").exists()).toBe(true);
    expect(wrapper.findAll(".card").length).toBeGreaterThan(0);
  });

  it("loads finance summary on mount", async () => {
    mountView(FinanceOverview);

    await flushPromises();

    expect(mockApi.finance.summary.get).toHaveBeenCalled();
  });

  it("displays summary cards with correct values", async () => {
    mockApi.finance.summary.get.mockResolvedValue({
      totalBalance: 5000,
      totalIncome: 1000,
      totalExpense: 200,
      cashFlow: 800,
      byAccount: [],
      monthlyTrends: [],
    });

    const wrapper = mountView(FinanceOverview);

    await flushPromises();

    const cards = wrapper.findAll(".card");
    expect(cards.length).toBeGreaterThanOrEqual(4);
  });
});
