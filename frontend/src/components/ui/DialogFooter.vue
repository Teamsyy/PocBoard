<template>
  <div class="dialog-footer">
    <slot />
  </div>
</template>

<style scoped>
.dialog-footer {
  padding: 1rem 1.5rem 1.5rem 1.5rem;
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .dialog-footer {
    flex-direction: column;
  }
}
</style>
