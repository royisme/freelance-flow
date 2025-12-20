import { defineStore } from "pinia";
import { ref, computed } from "vue";
import {
  Login,
  Register,
  GetAllUsers,
  GetUserByID,
  HasUsers,
  UpdateUser,
  ChangePassword,
} from "../wailsjs/go/services/AuthService";
import { dto } from "../wailsjs/go/models";
import { useAppStore } from "./app";

export const useAuthStore = defineStore("auth", () => {
  const appStore = useAppStore();
  // State
  const currentUser = ref<dto.UserOutput | null>(null);
  const usersList = ref<dto.UserListItem[]>([]);
  const isInitialized = ref(false);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  // Computed
  const isAuthenticated = computed(() => currentUser.value !== null);
  const userId = computed(() => currentUser.value?.id ?? 0);
  const username = computed(() => currentUser.value?.username ?? "");
  const avatarUrl = computed(() => currentUser.value?.avatarUrl ?? "");

  // Actions
  async function initialize() {
    isLoading.value = true;
    error.value = null;
    try {
      const savedUserId = localStorage.getItem("currentUserId");
      const hasExistingUsers = await HasUsers();

      if (hasExistingUsers) {
        const usersPromise = GetAllUsers();
        const userPromise =
          savedUserId != null
            ? (async () => {
                const id = parseInt(savedUserId, 10);
                if (isNaN(id)) return null;
                try {
                  return await GetUserByID(id);
                } catch {
                  localStorage.removeItem("currentUserId");
                  return null;
                }
              })()
            : Promise.resolve(null);

        const [users, user] = await Promise.all([usersPromise, userPromise]);
        usersList.value = users;
        if (user) {
          currentUser.value = user;
          await appStore.loadUserSettings(user.id);
        }
      }

      isInitialized.value = true;
    } catch (e) {
      error.value =
        e instanceof Error ? e.message : "Failed to initialize auth";
      console.error(e);
    } finally {
      isLoading.value = false;
    }
  }

  async function register(input: dto.RegisterInput) {
    isLoading.value = true;
    error.value = null;
    try {
      const user = await Register(input);
      currentUser.value = user;
      localStorage.setItem("currentUserId", String(user.id));
      usersList.value = await GetAllUsers(); // Refresh list including new user

      // Load settings in a separate try-catch to not affect auth flow
      try {
        await appStore.loadUserSettings(user.id);
      } catch (settingsError) {
        console.error(
          "Failed to load settings after registration:",
          settingsError
        );
        // Don't throw - settings will use defaults
      }

      return user;
    } catch (e) {
      error.value = e instanceof Error ? e.message : "Registration failed";
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  async function login(input: dto.LoginInput) {
    isLoading.value = true;
    error.value = null;
    try {
      const user = await Login(input);
      currentUser.value = user;
      localStorage.setItem("currentUserId", String(user.id));
      return user;
    } catch (e) {
      // Clean up local state on failure just in case
      currentUser.value = null;
      error.value = e instanceof Error ? e.message : "Login failed";
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  async function fetchAllUsers() {
    try {
      usersList.value = await GetAllUsers();
    } catch (e) {
      console.error("Failed to fetch users:", e);
    }
  }

  function logout() {
    currentUser.value = null;
    localStorage.removeItem("currentUserId");
    localStorage.removeItem("lastAuthedRoute");
  }

  function clearError() {
    error.value = null;
  }

  function switchUser() {
    // Clear current session but keep users list (logic is same as logout basically)
    currentUser.value = null;
    localStorage.removeItem("currentUserId");
    localStorage.removeItem("lastAuthedRoute");
  }

  async function updateProfile(input: dto.UpdateUserInput) {
    isLoading.value = true;
    error.value = null;
    try {
      const user = await UpdateUser(input);
      currentUser.value = user;
      await fetchAllUsers(); // Refresh list to update avatar/username
      return user;
    } catch (e) {
      error.value = e instanceof Error ? e.message : "Failed to update profile";
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  async function changePassword(input: dto.ChangePasswordInput) {
    isLoading.value = true;
    error.value = null;
    try {
      const success = await ChangePassword(input);
      return success;
    } catch (e) {
      error.value =
        e instanceof Error ? e.message : "Failed to change password";
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  return {
    // State
    currentUser,
    usersList,
    isInitialized,
    isLoading,
    error,
    // Computed
    isAuthenticated,
    userId,
    username,
    avatarUrl,
    // Actions
    initialize,
    register,
    login,
    logout,
    clearError,
    switchUser,
    fetchAllUsers,
    updateProfile,
    changePassword,
  };
});
