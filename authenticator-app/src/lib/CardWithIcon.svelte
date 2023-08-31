<script lang="ts">
    let clazz = '';
    export { clazz as class };
</script>

<div class="relative card {clazz}">
  <div class="absolute -top-10 left-0 right-0 flex justify-center">
    <span class="w-20 h-20 p-2 border bg-stone-50 dark:bg-stone-900 dark:border-stone-600 rounded-full">
    </span>
  </div>
  <div class="absolute -top-0 h-10 left-36 right-36 bg-stone-50 dark:bg-stone-900 flex justify-center fill-stone-950 dark:fill-stone-50">
    <div class="w-20 h-20 p-2 -mt-10">
      <slot name="icon" />
    </div>
  </div>
  <slot />
</div>
