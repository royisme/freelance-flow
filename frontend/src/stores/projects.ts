import { defineStore } from "pinia";
import { ref } from "vue";
import { api } from "@/api";
import type { Project, CreateProjectInput, UpdateProjectInput } from "@/types";

export const useProjectStore = defineStore("projects", () => {
  const projects = ref<Project[]>([]);
  const loading = ref(false);

  async function fetchProjects() {
    loading.value = true;
    try {
      projects.value = (await api.projects.list()) ?? [];
    } catch (error) {
      console.error("Failed to fetch projects", error);
    } finally {
      loading.value = false;
    }
  }

  async function fetchProjectsByClient(clientId: number) {
    loading.value = true;
    try {
      projects.value = (await api.projects.listByClient(clientId)) ?? [];
    } catch (error) {
      console.error("Failed to fetch projects by client", error);
    } finally {
      loading.value = false;
    }
  }

  async function createProject(input: CreateProjectInput) {
    loading.value = true;
    try {
      await api.projects.create(input);
      await fetchProjects();
    } finally {
      loading.value = false;
    }
  }

  async function updateProject(input: UpdateProjectInput) {
    loading.value = true;
    try {
      await api.projects.update(input);
      await fetchProjects();
    } finally {
      loading.value = false;
    }
  }

  async function deleteProject(id: number) {
    loading.value = true;
    try {
      await api.projects.delete(id);
      await fetchProjects();
    } finally {
      loading.value = false;
    }
  }

  return {
    projects,
    loading,
    fetchProjects,
    fetchProjectsByClient,
    createProject,
    updateProject,
    deleteProject,
  };
});
