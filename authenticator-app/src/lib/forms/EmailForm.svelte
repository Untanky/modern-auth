<script lang="ts">
  import Input from "$lib/utils/Input.svelte";

  export let disabled = false;
</script>

<Input
  type="email"
  value=""
  placeholder="Your Email"
  autocomplete="email"
  onInput={() => {}}
  id="email"
  label="email"
  disabled={disabled}
/>
