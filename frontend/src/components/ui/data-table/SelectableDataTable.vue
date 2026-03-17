<!--
  SelectableDataTable - 可选择行的数据表格组件
  
  基于 TanStack Table + shadcn-vue Table + Reka UI Checkbox 构建
  
  ============================================================================
  使用方式 (Usage)
  ============================================================================
  
  基础用法 - 带行选择的表格:
  
  ```vue
  <script setup lang="ts">
  import { ref } from 'vue'
  import type { ColumnDef } from '@tanstack/vue-table'
  import { SelectableDataTable } from '@/components/ui/data-table'
  
  interface User {
    id: number
    name: string
    email: string
  }
  
  const users = ref<User[]>([
    { id: 1, name: 'Alice', email: 'alice@example.com' },
    { id: 2, name: 'Bob', email: 'bob@example.com' },
  ])
  
  const selectedUserIds = ref<number[]>([])
  
  const columns: ColumnDef<User>[] = [
    { accessorKey: 'name', header: 'Name' },
    { accessorKey: 'email', header: 'Email' },
  ]
  </script>
  
  <template>
    <SelectableDataTable
      :columns="columns"
      :data="users"
      :selectable="true"
      v-model:selection="selectedUserIds"
    />
  </template>
  ```
  
  ============================================================================
  Props
  ============================================================================
  
  | Prop        | Type                    | Required | Default | Description                    |
  |-------------|-------------------------|----------|---------|--------------------------------|
  | columns     | ColumnDef<TData>[]      | Yes      | -       | TanStack Table 列定义           |
  | data        | TData[]                 | Yes      | -       | 数据数组，每行需要有 id 字段      |
  | selectable  | boolean                 | No       | false   | 是否显示选择列                   |
  | loading     | boolean                 | No       | false   | 是否显示加载状态                 |
  | rowIdField  | keyof TData             | No       | 'id'    | 行 ID 字段名                    |
  
  ============================================================================
  v-model
  ============================================================================
  
  v-model:selection - 双向绑定选中的行 ID 数组
  
  ```vue
  <SelectableDataTable v-model:selection="selectedIds" />
  ```
  
  ============================================================================
  Events
  ============================================================================
  
  | Event           | Payload              | Description          |
  |-----------------|----------------------|----------------------|
  | selectionChange | (ids: number[]) => void | 选择变化时触发       |
  
  ```vue
  <SelectableDataTable @selectionChange="handleSelectionChange" />
  ```
  
  ============================================================================
  Expose
  ============================================================================
  
  通过 ref 可以访问 TanStack Table 实例，用于高级操作:
  
  ```vue
  <script setup>
  const tableRef = ref()
  // tableRef.value.table.toggleAllRowsSelected()
  </script>
  
  <template>
    <SelectableDataTable ref="tableRef" />
  </template>
  ```
  
  ============================================================================
  重要说明 (Important Notes)
  ============================================================================
  
  1. Reka UI Checkbox 使用 modelValue/update:modelValue，不是 checked/update:checked
  2. 数据对象必须包含 id 字段（或通过 rowIdField 指定其他字段）
  3. 选中行会自动添加 data-state="selected" 属性，可用于样式
  4. TData 泛型必须 extends { id: number | string }
-->

<script setup lang="ts" generic="TData extends { id: number | string }, TValue">
import type { ColumnDef, RowSelectionState } from '@tanstack/vue-table'
import {
  FlexRender,
  getCoreRowModel,
  useVueTable,
} from '@tanstack/vue-table'
import { computed, h } from 'vue'
import { Loader2 } from 'lucide-vue-next'

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Checkbox } from '@/components/ui/checkbox'

const props = defineProps<{
  /** TanStack Table 列定义 */
  columns: ColumnDef<TData, TValue>[]
  /** 数据数组，每行需要有 id 字段 */
  data: TData[]
  /** 是否显示加载状态 */
  loading?: boolean
  /** 是否显示选择列（表头全选 + 行选择框） */
  selectable?: boolean
  /** 行 ID 字段名，默认为 'id' */
  rowIdField?: keyof TData
}>()

/** 双向绑定：选中的行 ID 数组 */
const selectedIds = defineModel<(number | string)[]>('selection', { default: () => [] })

const emit = defineEmits<{
  /** 选择变化时触发 */
  (e: 'selectionChange', ids: (number | string)[]): void
}>()

// 将 selectedIds 数组转换为 TanStack Table 的 RowSelectionState 对象
const rowSelectionState = computed<RowSelectionState>(() => {
  const state: RowSelectionState = {}
  selectedIds.value.forEach(id => {
    state[String(id)] = true
  })
  return state
})

// 选择列定义
// 注意：Reka UI Checkbox 使用 modelValue 而非 checked
const selectionColumn: ColumnDef<TData, TValue> = {
  id: 'select',
  header: ({ table }) => h(Checkbox, {
    modelValue: table.getIsAllPageRowsSelected(),
    'onUpdate:modelValue': (value: boolean) => table.toggleAllPageRowsSelected(!!value),
    ariaLabel: 'Select all',
  }),
  cell: ({ row }) => h(Checkbox, {
    modelValue: row.getIsSelected(),
    'onUpdate:modelValue': (value: boolean) => row.toggleSelected(!!value),
    ariaLabel: 'Select row',
  }),
  enableSorting: false,
  enableHiding: false,
}

// 合并选择列与用户定义的列
const allColumns = computed(() => {
  if (props.selectable) {
    return [selectionColumn, ...props.columns]
  }
  return props.columns
})

const table = useVueTable({
  get data() { return props.data },
  get columns() { return allColumns.value },
  getCoreRowModel: getCoreRowModel(),
  getRowId: (row) => String(row[props.rowIdField ?? 'id' as keyof TData]),
  onRowSelectionChange: (updaterOrValue) => {
    const newState = typeof updaterOrValue === 'function'
      ? updaterOrValue(rowSelectionState.value)
      : updaterOrValue
    
    // 将 RowSelectionState 转换回 ID 数组
    const ids = Object.entries(newState)
      .filter(([, selected]) => selected)
      .map(([id]) => {
        const numId = Number(id)
        return isNaN(numId) ? id : numId
      })
    
    selectedIds.value = ids
    emit('selectionChange', ids)
  },
  state: {
    get rowSelection() { return rowSelectionState.value },
  },
})

// 暴露 table 实例供高级用法
defineExpose({ table })
</script>

<template>
  <div class="rounded-md border">
    <Table>
      <TableHeader>
        <TableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
          <TableHead v-for="header in headerGroup.headers" :key="header.id">
            <FlexRender v-if="!header.isPlaceholder" :render="header.column.columnDef.header"
              :props="header.getContext()" />
          </TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <template v-if="loading">
          <TableRow>
            <TableCell :colspan="allColumns.length" class="h-24 text-center">
              <div class="flex items-center justify-center">
                <Loader2 class="h-6 w-6 animate-spin text-muted-foreground" />
              </div>
            </TableCell>
          </TableRow>
        </template>
        <template v-else-if="table.getRowModel().rows?.length">
          <TableRow v-for="row in table.getRowModel().rows" :key="row.id"
            :data-state="row.getIsSelected() ? 'selected' : undefined">
            <TableCell v-for="cell in row.getVisibleCells()" :key="cell.id">
              <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
            </TableCell>
          </TableRow>
        </template>
        <template v-else>
          <TableRow>
            <TableCell :colspan="allColumns.length" class="h-24 text-center">
              No results.
            </TableCell>
          </TableRow>
        </template>
      </TableBody>
    </Table>
  </div>
</template>
