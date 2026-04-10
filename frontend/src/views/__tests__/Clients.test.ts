import { describe, expect, it, vi, beforeEach } from "vitest";
import { flushPromises } from "@vue/test-utils";
import Clients from "@/views/Clients.vue";
import type { Client } from "@/types";
import { mountView } from "@/test-utils/mount";
import { useClientStore } from "@/stores/clients";

const mockApi = vi.hoisted(() => ({
  clients: {
    list: vi.fn<[], Promise<Client[]>>(),
    create: vi.fn<[Omit<Client, "id">], Promise<Client>>(),
    update: vi.fn<[Client], Promise<Client>>(),
    delete: vi.fn<[number], Promise<void>>(),
  },
  projects: {
    list: vi.fn<[], Promise<any[]>>(),
  },
  timesheet: {
    list: vi.fn<[], Promise<any[]>>(),
  },
  invoices: {
    list: vi.fn<[], Promise<any[]>>(),
  },
  statusBar: {
    get: vi.fn().mockResolvedValue({}),
  },
}));

vi.mock("@/api", () => ({ api: mockApi }));

vi.mock("vue-i18n", () => ({
  useI18n: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock("vue-router", () => ({
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
}));

describe("Clients view", () => {
  const mockClients: Client[] = [
    {
      id: 1,
      name: "Client A",
      email: "a@example.com",
      userId: 1,
      createdAt: "",
      updatedAt: "",
    },
    {
      id: 2,
      name: "Client B",
      email: "b@example.com",
      userId: 1,
      createdAt: "",
      updatedAt: "",
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
    mockApi.clients.list.mockResolvedValue(mockClients);
  });

  it("fetches clients on mount", async () => {
    const wrapper = mountView(Clients);

    await flushPromises();

    expect(mockApi.clients.list).toHaveBeenCalledTimes(1);
    const clientStore = useClientStore();
    expect(clientStore.clients.length).toBe(2);
    expect(wrapper.text()).toContain("clients.addClient");
  });
});
