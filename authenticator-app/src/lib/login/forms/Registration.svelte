<script lang="ts">
  export let submit: () => void;
</script>

<form class="flex flex-col flex-1 space-y-1">
  <h2 class="text-xl">
    Registration
  </h2>
  <p>
    When you click on register, prepare to be prompted to authenticate with your biometric data.
    Alternatively, you can use a physical hardware token.
  </p>
  <p>
    Once authenticated, we will create new user account for you.
  </p>
  <button type="button" on:click={submit} class="self-end mt-2 btn btn-yellow">
    Register
  </button>
</form>
