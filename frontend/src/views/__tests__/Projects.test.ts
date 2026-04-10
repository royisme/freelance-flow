import { describe, expect, it, vi, beforeEach } from "vitest";
import { flushPromises } from "@vue/test-utils";
import Projects from "@/views/Projects.vue";
import type { Client, Project } from "@/types";
import { mountView } from "@/test-utils/mount";

const mockApi = vi.hoisted(() => ({
  clients: {
    list: vi.fn<[], Promise<Client[]>>(),
  },
  projects: {
    list: vi.fn<[], Promise<Project[]>>(),
    create: vi.fn<[Omit<Project, "id">], Promise<Project>>(),
    update: vi.fn<[Project], Promise<Project>>(),
    delete: vi.fn<[number], Promise<void>>(),
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

describe("Projects view", () => {
  const mockClients: Client[] = [
    { id: 1, name: "Client A", email: "a@example.com", userId: 1, createdAt: "", updatedAt: "" },
  ];

  const mockProjects: Project[] = [
    { id: 1, name: "Project A", clientId: 1, hourlyRate: 50, status: "active", tags: [], deadline: null, userId: 1, createdAt: "", updatedAt: "" },
    { id: 2, name: "Project B", clientId: 1, hourlyRate: 75, status: "active", tags: [], deadline: null, userId: 1, createdAt: "", updatedAt: "" },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
    mockApi.clients.list.mockResolvedValue(mockClients);
    mockApi.projects.list.mockResolvedValue(mockProjects);
  });

  it("fetches projects and clients on mount", async () => {
    const wrapper = mountView(Projects);

    await flushPromises();

    expect(mockApi.projects.list).toHaveBeenCalledTimes(1);
    expect(mockApi.clients.list).toHaveBeenCalledTimes(1);
    expect(wrapper.text()).toContain("projects.addProject");
  });

});
