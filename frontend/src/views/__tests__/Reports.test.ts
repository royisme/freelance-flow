import { describe, expect, it, vi, beforeEach } from "vitest";
import { flushPromises } from "@vue/test-utils";
import Reports from "@/views/Reports.vue";
import type { Client, Project, ReportOutput } from "@/types";
import { mountView } from "@/test-utils/mount";

const mockApi = vi.hoisted(() => ({
  clients: {
    list: vi.fn<[], Promise<Client[]>>(),
  },
  projects: {
    list: vi.fn<[], Promise<Project[]>>(),
  },
  reports: {
    get: vi.fn<[Record<string, unknown>], Promise<ReportOutput>>(),
  },
}));

vi.mock("@/api", () => ({ api: mockApi }));
vi.mock("vue-router", () => ({
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
}));

describe("Reports view", () => {
  beforeEach(() => {
    mockApi.clients.list.mockResolvedValue([]);
    mockApi.projects.list.mockResolvedValue([]);
    mockApi.reports.get.mockResolvedValue({
      totalHours: 0,
      totalIncome: 0,
      rows: [],
      chart: { dates: [], revenue: [], hours: [] },
    });
  });

  it("fetches reports data on mount", async () => {
    const wrapper = mountView(Reports);

    await flushPromises();

    expect(mockApi.reports.get).toHaveBeenCalledTimes(1);
    expect(wrapper.text()).toContain("reports.stats.totalHours");
  });

  it("shows error state when api fails", async () => {
    mockApi.reports.get.mockRejectedValueOnce(new Error("boom"));
    const wrapper = mountView(Reports);

    await flushPromises();

    expect(wrapper.text()).toContain("boom");
  });

  it("apply triggers refetch with filters", async () => {
    const wrapper = mountView(Reports);

    await flushPromises();

    const buttons = wrapper.findAll("button");
    const apply = buttons.find(
      (b) => b.text().includes("reports.filters.apply")
    );
    await apply?.trigger("click");
    await flushPromises();

    expect(mockApi.reports.get.mock.calls.length).toBeGreaterThanOrEqual(2);
    const lastCall =
      mockApi.reports.get.mock.calls[mockApi.reports.get.mock.calls.length - 1];
    const lastArgs = lastCall[0];
    expect(lastArgs).toEqual(expect.any(Object));
  });
});
