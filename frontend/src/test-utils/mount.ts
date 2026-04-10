import { mount, type MountingOptions } from "@vue/test-utils";
import { createPinia } from "pinia";
import type { Component } from "vue";
import { defineComponent, h } from "vue";

const naiveStubs = {
  NForm: {
    props: ["model", "rules"],
    template: "<form><slot /></form>",
    methods: {
      validate() {
        const maybeModel = this.model as Record<string, unknown> | undefined;
        const maybeRules = this.rules as Record<string, unknown> | undefined;
        const passwordRulePresent =
          !!maybeRules && Object.prototype.hasOwnProperty.call(maybeRules, "password");

        const password = maybeModel?.password;
        if (passwordRulePresent && typeof password === "string" && password.trim() === "") {
          return Promise.reject(new Error("validation failed"));
        }
        return Promise.resolve();
      },
    },
  },
  NFormItem: {
    props: ["label"],
    template:
      '<div class="n-form-item"><label>{{label}}</label><slot /></div>',
  },
  NInput: {
    props: ["value", "disabled", "type", "autosize", "show-password-on"],
    emits: ["update:value"],
    template:
      '<textarea v-if="type === \'textarea\'" :value="value" :disabled="disabled" @input="$emit(\'update:value\', $event.target.value)"></textarea>' +
      '<input v-else :type="type || \'text\'" :value="value" :disabled="disabled" @input="$emit(\'update:value\', $event.target.value)" />',
  },
  NInputNumber: {
    props: ["value", "disabled"],
    emits: ["update:value"],
    template:
      '<input type="number" :value="value" :disabled="disabled" @input="$emit(\'update:value\', parseInt($event.target.value, 10))" />',
  },
  NSwitch: {
    props: ["value", "disabled"],
    emits: ["update:value"],
    template:
      '<input type="checkbox" :checked="value" :disabled="disabled" @change="$emit(\'update:value\', $event.target.checked)" />',
  },
  NSelect: {
    props: ["value", "options", "disabled"],
    emits: ["update:value"],
    template:
      '<select :value="value" :disabled="disabled" @change="$emit(\'update:value\', $event.target.value)"><option v-for="opt in options" :key="opt.value" :value="opt.value">{{opt.label}}</option></select>',
  },
  NSpace: {
    props: ["justify"],
    template: '<div class="n-space" :justify="justify"><slot /></div>',
  },
  NButton: {
    props: ["disabled", "loading", "type"],
    emits: ["click"],
    template:
      '<button :disabled="disabled" :type="type" @click="$emit(\'click\')"><slot /></button>',
  },
  NDatePicker: {
    props: ["value", "disabled"],
    emits: ["update:value"],
    template:
      '<input type="date" :value="value" :disabled="disabled" @input="$emit(\'update:value\', $event.target.value)" />',
  },
  NDataTable: {
    template: '<div class="n-data-table"><slot /></div>',
  },
  NModal: {
    props: ["show"],
    template: '<div class="n-modal" v-if="show"><slot /></div>',
  },
  NPopconfirm: {
    template:
      '<div class="n-popconfirm"><slot /><slot name="action" /></div>',
  },
  NIcon: {
    template: '<i class="n-icon"><slot /></i>',
  },
  NCard: {
    template: '<div class="n-card"><slot /></div>',
  },
  NText: {
    template: '<span class="n-text"><slot /></span>',
  },
  NEmpty: {
    template: '<div class="n-empty"><slot /></div>',
  },
  NTag: {
    template: '<span class="n-tag"><slot /></span>',
  },
  NDivider: {
    template: '<hr class="n-divider" />',
  },
  NSpin: {
    template: '<div class="n-spin"><slot /></div>',
  },
  NAvatar: {
    template: '<div class="n-avatar"></div>',
  },
  NConfigProvider: {
    template: '<div class="n-config-provider"><slot /></div>',
  },
  NPopover: {
    template: '<div class="n-popover"><slot /></div>',
  },
  NTooltip: {
    template: '<div class="n-tooltip"><slot /></div>',
  },
  Dialog: defineComponent({
    name: "Dialog",
    props: { open: Boolean },
    emits: ["update:open"],
    setup(props, { slots }) {
      return () => (props.open ? h("div", { class: "dialog-stub" }, slots.default?.()) : null);
    },
  }),
  DialogContent: { template: '<div class="dialog-content"><slot /></div>' },
  DialogHeader: { template: '<div class="dialog-header"><slot /></div>' },
  DialogTitle: { template: '<div class="dialog-title"><slot /></div>' },
  DialogFooter: { template: '<div class="dialog-footer"><slot /></div>' },
  DialogDescription: { template: '<div class="dialog-description"><slot /></div>' },
  AlertDialog: {
    props: ["open"],
    emits: ["update:open"],
    template: '<div class="alert-dialog"><slot /></div>',
  },
  AlertDialogTrigger: { template: '<div class="alert-dialog-trigger"><slot /></div>' },
  AlertDialogContent: { template: '<div class="alert-dialog-content"><slot /></div>' },
  AlertDialogHeader: { template: '<div class="alert-dialog-header"><slot /></div>' },
  AlertDialogTitle: { template: '<div class="alert-dialog-title"><slot /></div>' },
  AlertDialogDescription: { template: '<div class="alert-dialog-description"><slot /></div>' },
  AlertDialogFooter: { template: '<div class="alert-dialog-footer"><slot /></div>' },
  AlertDialogAction: {
    emits: ["click"],
    template: '<button class="alert-dialog-action" @click="$emit(\'click\')"><slot /></button>',
  },
  AlertDialogCancel: {
    emits: ["click"],
    template: '<button class="alert-dialog-cancel" @click="$emit(\'click\')"><slot /></button>',
  },
  Select: {
    props: ["modelValue", "disabled"],
    emits: ["update:modelValue"],
    template: '<div class="select-stub"><slot /></div>',
  },
  SelectTrigger: { template: '<button class="select-trigger"><slot /></button>' },
  SelectValue: {
    props: ["placeholder"],
    template: '<span class="select-value">{{ placeholder }}</span>',
  },
  SelectContent: { template: '<div class="select-content"><slot /></div>' },
  SelectItem: {
    props: ["value"],
    template: '<div class="select-item" :data-value="value"><slot /></div>',
  },
  SelectGroup: { template: '<div class="select-group"><slot /></div>' },
  Switch: {
    props: ["checked", "disabled"],
    emits: ["update:checked"],
    template:
      '<input type="checkbox" :checked="checked" :disabled="disabled" @change="$emit(\'update:checked\', $event.target.checked)" />',
  },
  Tabs: { template: '<div class="tabs"><slot /></div>' },
  TabsList: { template: '<div class="tabs-list"><slot /></div>' },
  TabsTrigger: { template: '<button class="tabs-trigger"><slot /></button>' },
  TabsContent: { template: '<div class="tabs-content"><slot /></div>' },
  Card: { template: '<div class="card"><slot /></div>' },
  CardHeader: { template: '<div class="card-header"><slot /></div>' },
  CardTitle: { template: '<div class="card-title"><slot /></div>' },
  CardContent: { template: '<div class="card-content"><slot /></div>' },
  Alert: { template: '<div class="alert"><slot /></div>' },
  AlertTitle: { template: '<div class="alert-title"><slot /></div>' },
  AlertDescription: { template: '<div class="alert-description"><slot /></div>' },
};

export function mountView<T extends Component>(
  component: T,
  options: MountingOptions<unknown> = {}
) {
  const globalOptions = options.global ?? {};
  const existingStubs =
    (globalOptions.stubs as Record<string, unknown> | undefined) ?? {};
  return mount(component, {
    ...options,
    global: {
      ...globalOptions,
      plugins: [...((globalOptions.plugins as unknown[]) ?? []), createPinia()],
      stubs: { ...naiveStubs, ...existingStubs },
    },
  });
}
