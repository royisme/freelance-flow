import { vi } from "vitest";
import { defineComponent, h } from "vue";
import { mockDialog, mockMessage } from "@/test-utils/naive";
import { applyThemeToRoot } from "@/theme/tokens";

applyThemeToRoot("light");

const Stub = (name: string) =>
  defineComponent({
    name,
    setup(_, { slots }) {
      return () => h("div", { class: name.toLowerCase() }, slots.default?.());
    },
  });

const StubModal = defineComponent({
  name: "NModal",
  props: { show: { type: Boolean, default: false } },
  setup(props, { slots }) {
    return () =>
      props.show
        ? h("div", { class: "n-modal" }, slots.default?.())
        : null;
  },
});

const StubDatePicker = defineComponent({
  name: "NDatePicker",
  props: {
    value: { type: [String, Number, Array], default: "" },
    disabled: Boolean,
  },
  emits: ["update:value"],
  setup(props, { emit }) {
    const getValueString = () => {
      if (Array.isArray(props.value)) {
        return String(props.value[0] ?? "");
      }
      return String(props.value ?? "");
    };
    return () =>
      h("input", {
        type: "date",
        disabled: props.disabled,
        value: getValueString(),
        onInput: (e: Event) =>
          emit("update:value", (e.target as HTMLInputElement).value),
      });
  },
});

const StubForm = defineComponent({
  name: "NForm",
  setup(_, { slots }) {
    return () => h("form", {}, slots.default?.());
  },
});

const StubFormItem = defineComponent({
  name: "NFormItem",
  props: { label: { type: String, default: "" } },
  setup(props, { slots }) {
    return () =>
      h("div", { class: "n-form-item" }, [
        props.label ? h("label", {}, props.label) : null,
        slots.default?.(),
      ]);
  },
});

const StubInput = defineComponent({
  name: "NInput",
  props: {
    value: { type: [String, Number], default: "" },
    disabled: Boolean,
    type: { type: String, default: "text" },
  },
  emits: ["update:value"],
  setup(props, { emit, slots }) {
    return () => {
      if (props.type === "textarea") {
        return h("textarea", {
          disabled: props.disabled,
          value: String(props.value ?? ""),
          onInput: (e: Event) =>
            emit("update:value", (e.target as HTMLTextAreaElement).value),
        });
      }
      return h(
        "input",
        {
          type: props.type || "text",
          disabled: props.disabled,
          value: String(props.value ?? ""),
          onInput: (e: Event) =>
            emit("update:value", (e.target as HTMLInputElement).value),
        },
        slots
      );
    };
  },
});

const StubInputNumber = defineComponent({
  name: "NInputNumber",
  props: { value: { type: Number, default: 0 }, disabled: Boolean },
  emits: ["update:value"],
  setup(props, { emit }) {
    return () =>
      h("input", {
        type: "number",
        disabled: props.disabled,
        value: String(props.value ?? 0),
        onInput: (e: Event) =>
          emit(
            "update:value",
            parseInt((e.target as HTMLInputElement).value, 10)
          ),
      });
  },
});

const StubSwitch = defineComponent({
  name: "NSwitch",
  props: { value: { type: Boolean, default: false }, disabled: Boolean },
  emits: ["update:value"],
  setup(props, { emit }) {
    return () =>
      h("input", {
        type: "checkbox",
        disabled: props.disabled,
        checked: props.value,
        onChange: (e: Event) =>
          emit("update:value", (e.target as HTMLInputElement).checked),
      });
  },
});

const StubSelect = defineComponent({
  name: "NSelect",
  props: {
    value: { type: [String, Number], default: "" },
    options: {
      type: Array as unknown as () => Array<{ label: string; value: string }>,
      default: () => [],
    },
    disabled: Boolean,
  },
  emits: ["update:value"],
  setup(props, { emit }) {
    return () =>
      h(
        "select",
        {
          disabled: props.disabled,
          value: String(props.value ?? ""),
          onChange: (e: Event) =>
            emit("update:value", (e.target as HTMLSelectElement).value),
        },
        props.options.map((opt) =>
          h("option", { value: opt.value }, opt.label)
        )
      );
  },
});

const StubButton = defineComponent({
  name: "NButton",
  props: { disabled: Boolean, loading: Boolean, type: String },
  emits: ["click"],
  setup(props, { emit, slots }) {
    return () =>
      h(
        "button",
        {
          disabled: props.disabled || props.loading,
          type: props.type,
          onClick: () => emit("click"),
        },
        slots.default?.()
      );
  },
});

vi.mock("naive-ui", async () => {
  const actual = await vi.importActual<typeof import("naive-ui")>("naive-ui");
  const overrides = {
    useMessage: () => mockMessage,
    useDialog: () => mockDialog,
    NForm: StubForm,
    NFormItem: StubFormItem,
    NInput: StubInput,
    NInputNumber: StubInputNumber,
    NSwitch: StubSwitch,
    NSelect: StubSelect,
    NButton: StubButton,
    NDatePicker: StubDatePicker,
    NDataTable: Stub("NDataTable"),
    NModal: StubModal,
    NPopconfirm: Stub("NPopconfirm"),
    NIcon: Stub("NIcon"),
    NCard: Stub("NCard"),
    NEmpty: Stub("NEmpty"),
    NTag: Stub("NTag"),
    NDivider: Stub("NDivider"),
    NSpin: Stub("NSpin"),
    NAvatar: Stub("NAvatar"),
    NConfigProvider: Stub("NConfigProvider"),
    NPopover: Stub("NPopover"),
    NTooltip: Stub("NTooltip"),
  };
  return new Proxy({ ...actual, ...overrides }, {
    get(target, prop, receiver) {
      if (prop in target) return Reflect.get(target, prop, receiver);
      if (typeof prop === "string" && prop.startsWith("N")) return Stub(prop);
      return undefined;
    },
  });
});

const iconStub = { render: () => null };

vi.mock("@vicons/antd", () => {
  const handler: ProxyHandler<Record<string, unknown>> = {
    get: () => iconStub,
    has: () => true,
    getOwnPropertyDescriptor: () => ({
      configurable: true,
      enumerable: true,
      value: iconStub,
    }),
  };
  return new Proxy({}, handler);
});

vi.mock("vue-i18n", () => ({
  useI18n: () => ({
    t: (key: string) => key,
    locale: { value: "en-US" },
  }),
}));

vi.mock("@vicons/ionicons5", () => {
  const handler: ProxyHandler<Record<string, unknown>> = {
    get: () => iconStub,
    has: () => true,
    getOwnPropertyDescriptor: () => ({
      configurable: true,
      enumerable: true,
      value: iconStub,
    }),
  };
  return new Proxy({}, handler);
});
