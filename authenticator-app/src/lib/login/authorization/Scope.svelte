<script lang="ts">
  import type { Scope } from './models';

  export let scope: Scope;

  let value = true;
  
  // eslint-disable-next-line prefer-destructuring
  $: id = scope.id;
</script>

<div class="flex items-center">
  <input class="w-4 h-4 mx-3" type="checkbox" bind:checked={value} id={id}>
  <label for={id}>
    {scope.description}
  </label>
</div>
