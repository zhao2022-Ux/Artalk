<template>
  <div>
    <h3>Custom Options Editor</h3>

    <form @submit.prevent="handleSubmit">
      <div v-for="(value, path) in formValues" :key="path" class="form-row">
        <label>{{ path }}</label>
        <input v-model="formValues[path]" type="text" />
      </div>
      <button type="submit">Save</button>
    </form>

    <h3>Generated JSON</h3>
    <code>{{ generatedJson }}</code>
  </div>
</template>

<script setup lang="ts">
// 接受组件的属性
const props = defineProps<{
  defaultValue: Record<string, any>
  schema: Record<string, any>
}>()

const $emit = defineEmits<{
  save: [value: Record<string, any>]
}>()

// 扁平化 JSON Schema 的函数
function flattenSchema(schema: any, prefix = ''): Record<string, any> {
  let result: Record<string, any> = {}

  for (let key in schema.properties) {
    const newKey = prefix ? `${prefix}.${key}` : key

    if (schema.properties[key].type === 'object') {
      Object.assign(result, flattenSchema(schema.properties[key], newKey))
    } else {
      result[newKey] = '' // 给每个 final property 的默认值：空字符串
    }
  }

  return result
}

// 初始化表单数据
const formValues = reactive<Record<string, any>>({})

// 当组件挂载时，初始化表单
onMounted(() => {
  updateSchemaAndDefaultValues(props.schema, props.defaultValue)
})

// 监控 props.schema 或者 props.defaultValue 的变化
watch([() => props.schema, () => props.defaultValue], ([newSchema, newDefaultValue]) => {
  updateSchemaAndDefaultValues(newSchema, newDefaultValue)
})

// 初始化 schema 和 defaultValue 的处理逻辑，并更新表单值
const updateSchemaAndDefaultValues = (
  schema: Record<string, any>,
  defaultValues: Record<string, any>,
) => {
  // 生成 schema 默认的扁平化结构
  const flattenedSchema = flattenSchema(schema)

  // 使用 defaultValue 初始化表单值
  for (let path in flattenedSchema) {
    formValues[path] =
      defaultValues && path in defaultValues ? defaultValues[path] : flattenedSchema[path]
  }
}

// 将表单内容转换为 JSON 结构
const generatedJson = computed(() => {
  const result: any = {}

  const setNestedValue = (obj: Record<string, any>, path: string[], value: any) => {
    let current = obj
    for (let i = 0; i < path.length - 1; i++) {
      if (!(path[i] in current) || typeof current[path[i]] !== 'object') {
        current[path[i]] = {}
      }
      current = current[path[i]]
    }
    current[path[path.length - 1]] = value
  }

  // 复原嵌套结构
  for (let path in formValues) {
    const keys = path.split('.')
    setNestedValue(result, keys, formValues[path])
  }

  return result
})

// 提交表单
const handleSubmit = () => {
  console.log('Form Data Submitted:', formValues)

  // 提交更新到父组件
  $emit('save', generatedJson.value)
}
</script>

<style lang="scss" scoped>
.form-row {
  margin-bottom: 10px;
}
label {
  margin-right: 10px;
  font-weight: bold;
}
input {
  padding: 5px;
}
button {
  margin-top: 20px;
}
</style>
