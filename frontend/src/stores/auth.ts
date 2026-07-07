import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { getProfileApi, loginApi } from "@/api/auth";
import type { AdminUser, LoginRequest } from "@/types/system";
import { readStorage, writeStorage } from "@/utils/storage";

const TOKEN_KEY = "mall_admin_token";
const PROFILE_KEY = "mall_admin_profile";

export const useAuthStore = defineStore("auth", () => {
  const token = ref(localStorage.getItem(TOKEN_KEY) || "");
  const profile = ref<AdminUser | null>(readStorage<AdminUser>(PROFILE_KEY));
  const isLoggedIn = computed(() => Boolean(token.value));

  function setSession(nextToken: string, nextProfile: AdminUser): void {
    token.value = nextToken;
    profile.value = nextProfile;
    localStorage.setItem(TOKEN_KEY, nextToken);
    writeStorage(PROFILE_KEY, nextProfile);
  }

  async function login(payload: LoginRequest): Promise<void> {
    const data = await loginApi(payload);
    setSession(data.token, data.profile);
  }

  async function fetchProfile(): Promise<void> {
    profile.value = await getProfileApi();
    writeStorage(PROFILE_KEY, profile.value);
  }

  function logout(): void {
    token.value = "";
    profile.value = null;
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(PROFILE_KEY);
    localStorage.removeItem("mall_admin_tags");
  }

  return {
    token,
    profile,
    isLoggedIn,
    setSession,
    login,
    fetchProfile,
    logout
  };
});
