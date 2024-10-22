<template>
  <div class="modal-overlay" @click="handleOverlayClick">
    <div class="modal-content" @click.stop>
      <header class="modal-header">
        <h2>{{ title }}</h2>
        <button class="close-button" @click="close">x</button>
      </header>
      <div class="modal-body">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
const props = defineProps<{
  title?: string
}>()

const emit = defineEmits<{
  close: [value: void]
}>()

const close = () => {
  emit('close')
}

const handleOverlayClick = (event: MouseEvent) => {
  close()
}
</script>

<style lang="scss" scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 12px;
  padding: 20px;
  width: 400px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
  position: relative;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
  margin-bottom: 20px;
}

.close-button {
  background-color: transparent;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  font-weight: bold;
}

.modal-body {
  font-size: 1rem;
}
</style>
