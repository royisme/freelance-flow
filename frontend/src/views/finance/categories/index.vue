<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Plus } from 'lucide-vue-next'
import { api } from '@/api'
import { toast } from 'vue-sonner'

const { t } = useI18n()

const loading = ref(false)
const categories = ref<any[]>([])
const showCreateModal = ref(false)

const newCategory = ref({
  name: '',
  type: 'expense',
  color: '#000000',
  icon: '',
})

async function loadCategories() {
  loading.value = true
  try {
    categories.value = await api.finance.categories.list()
  } catch (e) {
    toast.error('Failed to load categories')
  } finally {
    loading.value = false
  }
}

async function createCategory() {
  try {
    await api.finance.categories.create(newCategory.value)
    toast.success(t('common.saved'))
    showCreateModal.value = false
    loadCategories()
    newCategory.value = { name: '', type: 'expense', color: '#000000', icon: '' }
  } catch (e) {
    toast.error(t('common.error'))
  }
}

onMounted(() => {
  loadCategories()
})
</script>

<template>
  <div class="categories-page space-y-6">
    <div class="page-header flex justify-between items-center">
      <h1 class="page-title text-2xl font-bold">{{ t('finance.categories.title') }}</h1>
      <Dialog v-model:open="showCreateModal">
        <DialogTrigger as-child>
          <Button>
            <Plus class="mr-2 h-4 w-4" />
            {{ t('finance.categories.addCategory') || 'Add Category' }}
          </Button>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{{ t('finance.categories.addCategory') || 'Add Category' }}</DialogTitle>
          </DialogHeader>
          <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="name" class="text-right">{{ t('common.name') }}</Label>
              <Input id="name" v-model="newCategory.name" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="type" class="text-right">Type</Label>
              <Select v-model="newCategory.type">
                <SelectTrigger class="col-span-3">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="income">Income</SelectItem>
                  <SelectItem value="expense">Expense</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="color" class="text-right">Color</Label>
              <Input id="color" type="color" v-model="newCategory.color" class="col-span-3 h-10" />
            </div>
          </div>
          <DialogFooter>
            <Button @click="createCategory">{{ t('common.save') }}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>

    <Card>
      <CardContent class="p-0">
        <div class="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{{ t('common.name') }}</TableHead>
                <TableHead>Type</TableHead>
                <TableHead>Color</TableHead>
                <TableHead class="w-[100px]">{{ t('common.actions') }}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <template v-if="categories.length > 0">
                <TableRow v-for="category in categories" :key="category.id">
                  <TableCell class="font-medium">{{ category.name }}</TableCell>
                  <TableCell>
                     <Badge :variant="category.type === 'income' ? 'default' : 'secondary'">{{ category.type }}</Badge>
                  </TableCell>
                  <TableCell>
                    <div class="w-6 h-6 rounded-full border" :style="{ backgroundColor: category.color }"></div>
                  </TableCell>
                  <TableCell>
                     <Button size="sm" variant="ghost" class="h-8 px-2 text-destructive hover:text-destructive/90">
                        {{ t('common.delete') }}
                      </Button>
                  </TableCell>
                </TableRow>
              </template>
              <TableRow v-else>
                <TableCell colspan="4" class="h-24 text-center">
                  {{ loading ? t('common.loading') : t('common.noData') }}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
