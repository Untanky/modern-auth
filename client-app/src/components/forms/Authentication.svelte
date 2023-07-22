<script lang="ts">
  export let submit: () => void;
</script>

<form class="flex flex-col flex-1">
  <h2 class="text-xl">
    Authentication
  <p>
    When you click on authenticate, prepare to be prompted to authenticate with your biometric data.
    Alternatively, you can use a physical hardware token.
  </p>
  <button type="button" on:click={submit} class="self-end mt-2 btn btn-primary">
    Authenticate
  </button>
</form>
