<script lang="ts">
  export let submit: () => void;
  export let userId: string;
</script>

<form class="flex flex-col flex-1">
  <p>
    There is no credential for user id <span class="font-mono">{userId}</span>. Please create one.
  </p>
  <button type="button" on:click={submit} class="self-end mt-2 btn btn-primary">
    Create credential
  </button>
</form>
