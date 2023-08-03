<script lang="ts">
  export let label, value, id: string;
  export let placeholder: string = '';
  export let autocomplete: string = '';
  export let type: string = 'text';
  export let onInput: (event: Event) => void;
  export let disabled = false;
</script>

<div class="relative mt-4">
  <label class="font-extralight text-sm uppercase {disabled && 'dark:text-stone-400'} dark:bg-stone-800 px-1 left-3 absolute block rounded" for={id}>
    {label}
  </label>
  <input
    value={value}
    on:input={onInput}
    class="dark:bg-stone-800 mt-3 px-4 py-2 disabled:dark:text-stone-400  {disabled && 'dark:border-stone-700'} border rounded-lg w-full"
    type={type}
    id={id}
    placeholder={placeholder}
    autocomplete={autocomplete}
    disabled={disabled}
  />
</div>