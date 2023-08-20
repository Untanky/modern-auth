<script lang="ts">
  let clazz = '';

  export { clazz as class };
  export let open = false;
  export let click: () => void;

  const onKeyPress = (event: KeyboardEvent) => {
      if (event.key === 'Space' || event.key === 'Enter') {
          click();
      }
  };
</script>

<div
  class="px-4 py-4 border space-y-2 { open ? 'border-yellow-500' : 'border-stone-500' } flex flex-col rounded-lg hover:cursor-pointer {clazz}"
  role="radio"
  aria-checked={open}
  tabindex={open ? -1 : 0}
  on:click={click}
  on:keypress={onKeyPress}
>
  <slot />
</div>
