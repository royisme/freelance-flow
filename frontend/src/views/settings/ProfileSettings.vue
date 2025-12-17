<script setup lang="ts">
import { ref, computed } from "vue";
import { RefreshCw, Upload } from "lucide-vue-next";
import { useI18n } from "vue-i18n";
import { toast } from "vue-sonner";
import { toTypedSchema } from "@vee-validate/zod";

import { useAuthStore } from "@/stores/auth";
import { dto } from "@/wailsjs/go/models";

import { profileSettingsSchema, changePasswordSchema } from "@/schemas/settings";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import { Separator } from "@/components/ui/separator";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

const authStore = useAuthStore();
const { t } = useI18n();

// Profile Form
const profileFormSchema = toTypedSchema(profileSettingsSchema);
const passwordFormSchema = toTypedSchema(changePasswordSchema);

const fileInputRef = ref<HTMLInputElement | null>(null);
const accountSaving = ref(false);

// Computed initial values that update when currentUser changes
const profileInitialValues = computed(() => ({
  username: authStore.currentUser?.username || "",
  email: authStore.currentUser?.email || "",
  avatarUrl: authStore.currentUser?.avatarUrl || "",
}));

function handleRandomAvatar(setFieldValue: (field: string, value: unknown) => void) {
  const seed = Math.random().toString(36).substring(7);
  // Use Dicebear Avataaars
  const newAvatarUrl = `https://api.dicebear.com/9.x/avataaars/svg?seed=${seed}`;
  setFieldValue('avatarUrl', newAvatarUrl);
}

function handleUploadAvatar() {
  fileInputRef.value?.click();
}

function onFileChange(event: Event, setFieldValue: (field: string, value: unknown) => void) {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (file) {
    const reader = new FileReader();
    reader.onload = (e) => {
      if (typeof e.target?.result === "string") {
        setFieldValue('avatarUrl', e.target.result);
      }
    };
    reader.readAsDataURL(file);
  }
}

async function onProfileSubmit(formValues: unknown) {
  const values = formValues as { username: string; email: string; avatarUrl?: string };
  if (!authStore.currentUser) return;

  try {
    accountSaving.value = true;
    const input = new dto.UpdateUserInput({
      id: authStore.currentUser.id,
      username: values.username,
      email: values.email,
      avatarUrl: values.avatarUrl || "",
      settingsJson: authStore.currentUser.settingsJson, // Keep existing settings
    });

    await authStore.updateProfile(input);
    toast.success(t("settings.profile.messages.saved"));
  } catch (e) {
    toast.error(e instanceof Error ? e.message : t("settings.profile.messages.saveError"));
  } finally {
    accountSaving.value = false;
  }
}

async function onPasswordSubmit(
  formValues: unknown,
  { resetForm }: { resetForm: () => void }
) {
  const values = formValues as { oldPassword: string; newPassword: string };
  if (!authStore.currentUser) return;

  try {
    accountSaving.value = true;
    const input = new dto.ChangePasswordInput({
      id: authStore.currentUser.id,
      oldPassword: values.oldPassword,
      newPassword: values.newPassword,
    });

    await authStore.changePassword(input);
    toast.success(t("settings.profile.messages.saved"));
    resetForm();
  } catch (e) {
    toast.error(e instanceof Error ? e.message : t("settings.profile.messages.saveError"));
  } finally {
    accountSaving.value = false;
  }
}

function getInitials(username: string): string {
  return username?.slice(0, 2).toUpperCase() || "U";
}
</script>

<template>
  <div class="profile-settings space-y-6">
    <Card>
      <CardHeader>
        <CardTitle>{{ t("settings.profile.tabs.account") }}</CardTitle>
      </CardHeader>
      <CardContent>
        <!-- Profile Form -->
        <Form v-slot="{ values, setFieldValue }" :validation-schema="profileFormSchema"
          :initial-values="profileInitialValues" :key="authStore.currentUser?.id" @submit="onProfileSubmit">
          <div class="flex gap-12 items-start mb-6">
            <!-- Left Column: Avatar & Actions -->
            <div class="flex flex-col items-center gap-5 min-w-[160px] pt-2">
              <FormField v-slot="{ value }" name="avatarUrl">
                <div class="relative inline-block">
                  <Avatar class="size-28">
                    <AvatarImage :src="value || ''" :alt="values.username" />
                    <AvatarFallback>{{ getInitials(values.username || '') }}</AvatarFallback>
                  </Avatar>

                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger as-child>
                        <Button variant="outline" size="icon" type="button"
                          class="absolute bottom-0 right-0 size-8 rounded-full shadow-md"
                          @click="handleRandomAvatar(setFieldValue)">
                          <RefreshCw class="size-4" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>
                        {{ t("settings.profile.fields.randomizeAvatar") }}
                      </TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                </div>
              </FormField>

              <Button variant="outline" type="button" @click="handleUploadAvatar" :disabled="accountSaving">
                <Upload class="size-4 mr-2" />
                {{ t("settings.profile.fields.uploadAvatar") }}
              </Button>

              <input type="file" ref="fileInputRef" class="hidden" accept="image/*"
                @change="(e) => onFileChange(e, setFieldValue)" />
            </div>

            <!-- Right Column: Form Fields -->
            <div class="grow max-w-md space-y-4">
              <FormField v-slot="{ componentField }" name="username">
                <FormItem>
                  <FormLabel>{{ t("settings.profile.fields.username") }}</FormLabel>
                  <FormControl>
                    <Input v-bind="componentField" :disabled="accountSaving" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="email">
                <FormItem>
                  <FormLabel>{{ t("settings.profile.fields.email") }}</FormLabel>
                  <FormControl>
                    <Input v-bind="componentField" :disabled="accountSaving" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <div class="flex justify-end pt-2">
                <Button type="submit" :disabled="accountSaving">
                  <span v-if="accountSaving">{{ t("common.saving") }}</span>
                  <span v-else>{{ t("common.save") }}</span>
                </Button>
              </div>
            </div>
          </div>
        </Form>

        <Separator class="my-6" />

        <!-- Password Form -->
        <h3 class="text-lg font-semibold mb-4">{{ t("settings.profile.password.title") }}</h3>
        <Form :validation-schema="passwordFormSchema" @submit="onPasswordSubmit">
          <div class="max-w-md space-y-4">
            <FormField v-slot="{ componentField }" name="oldPassword">
              <FormItem>
                <FormLabel>{{ t("settings.profile.password.current") }}</FormLabel>
                <FormControl>
                  <Input type="password" v-bind="componentField" :disabled="accountSaving" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>

            <FormField v-slot="{ componentField }" name="newPassword">
              <FormItem>
                <FormLabel>{{ t("settings.profile.password.new") }}</FormLabel>
                <FormControl>
                  <Input type="password" v-bind="componentField" :disabled="accountSaving" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>

            <FormField v-slot="{ componentField }" name="confirmPassword">
              <FormItem>
                <FormLabel>{{ t("settings.profile.password.confirm") }}</FormLabel>
                <FormControl>
                  <Input type="password" v-bind="componentField" :disabled="accountSaving" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>

            <div class="flex justify-end pt-2">
              <Button type="submit" variant="destructive" :disabled="accountSaving">
                {{ t("settings.profile.password.changeButton") }}
              </Button>
            </div>
          </div>
        </Form>
      </CardContent>
    </Card>
  </div>
</template>

<style scoped>
/* No scoped styles needed with Tailwind */
</style>
