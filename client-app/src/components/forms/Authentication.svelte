<script lang="ts">
  export let submit: () => void;

  let method: 'public-key' | 'password' = 'public-key';
</script>

<form class="flex flex-col flex-1 space-y-2">
  <h2 class="text-xl">
    Authentication
  </h2>
  <p>
    Please select how to authenticate:
  </p>
  <div class="space-y-4" role="radiogroup">
    <div
      class={`radio-select-box ${method === 'public-key' ? 'radio-select-box-active' : ''} flex flex-col hover:cursor-pointer`}
      role="radio"
      aria-checked={method === 'public-key'}
      tabindex={method === 'public-key' ? -1 : 0}
      on:click={() => method = 'public-key'}
      on:keypress={() => method = 'public-key'}
    >
      <h3 class="text-lg">Biometric/physical authentication</h3>
      {#if method === 'public-key'}
      <p>
        When you click on authenticate a system dialog will open and ask you to authenticate with your biometric data or a physical hardware token. Please prepare for the method chosen when setting up this device.
      </p>
      <button type="button" on:click={submit} class="self-end mt-2 btn btn-primary">
        Authenticate
      </button>
      {/if}
    </div>
    <div
      class={`radio-select-box ${method === 'password' ? 'radio-select-box-active' : ''} flex flex-col hover:cursor-pointer`}
      role="radio"
      aria-checked={method === 'password'}
      tabindex={method === 'password' ? -1 : 0}
      on:click={() => method = 'password'}
      on:keypress={() => method = 'password'}
    >
      <h3 class="text-lg">Password authentication</h3>
      {#if method === 'password'}
      <p>
        Please enter your password:
      </p>
      <input 
        class="dark:bg-stone-800 mt-3 px-4 py-2 dark:border-stone-600 border rounded-lg w-full"
        type="password"
        tabindex="-1"
      >
      <button type="button" on:click={submit} class="self-end mt-2 btn btn-primary">
        Continue
      </button>
      {/if}
    </div>
  </div>
</form>

<style scoped>
  .radio-select-box {
    @apply px-4 py-4 border space-y-2 border-stone-500 flex flex-col rounded-lg;
  }

  .radio-select-box-active {
    @apply border-yellow-500;
  }
</style>
