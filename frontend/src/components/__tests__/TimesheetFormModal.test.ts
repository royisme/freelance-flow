import { describe, expect, it } from "vitest";
import { shallowMount } from "@vue/test-utils";
import TimesheetFormModal from "@/components/TimesheetFormModal.vue";
import type { Project } from "@/types";
import { defineComponent } from "vue";

describe("TimesheetFormModal", () => {
  const DialogStub = defineComponent({
    name: "Dialog",
    props: { open: { type: Boolean, default: false } },
    emits: ["update:open"],
    template: '<div class="dialog-stub" v-if="open"><slot /></div>',
  });

  const projects: Project[] = [
    {
      id: 1,
      clientId: 1,
      name: "Project A",
      description: "",
      hourlyRate: 100,
      currency: "USD",
      status: "active",
      deadline: "",
      tags: [],
      userId: 1,
      createdAt: "",
      updatedAt: "",
    },
  ];

  it("forwards Dialog update:open=true", async () => {
    const wrapper = shallowMount(TimesheetFormModal, {
      props: {
        show: true,
        entry: null,
        projects,
      },
      global: {
        stubs: {
          Dialog: DialogStub,
        },
      },
    });

    wrapper.findComponent({ name: "Dialog" }).vm.$emit("update:open", true);

    expect(wrapper.emitted("update:show")?.[0]).toEqual([true]);
  });

  it("forwards Dialog update:open=false", async () => {
    const wrapper = shallowMount(TimesheetFormModal, {
      props: {
        show: true,
        entry: null,
        projects,
      },
      global: {
        stubs: {
          Dialog: DialogStub,
        },
      },
    });

    wrapper.findComponent({ name: "Dialog" }).vm.$emit("update:open", false);

    expect(wrapper.emitted("update:show")?.[0]).toEqual([false]);
  });
});
