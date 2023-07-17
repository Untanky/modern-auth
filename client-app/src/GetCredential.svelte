<script lang="ts">
  export let submit: () => void;
  export let userId: string;
</script>

<form class="flex flex-col flex-1">
  <p>
    Found credential for <span class="font-mono">{userId}</span>. After pressing the button, you will be asked to provide your hardware token or biometric signature.
  </p>
  <button type="button" on:click={submit} class="self-end mt-2 btn btn-primary">
    Login
  </button>
</form>
