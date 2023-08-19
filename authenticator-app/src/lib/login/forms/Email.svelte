<script lang="ts">
  import EmailForm from "$lib/forms/EmailForm.svelte";

  export let submit: () => void;
</script>

<form class="flex flex-col flex-1 space-y-1">
  <h2 class="text-xl">
    Registration
  </h2>
  <p>
    Please enter email address.
  </p>
  <EmailForm />
  Select how we use your email:
  <ul>
    <li>
      <input type="checkbox" id="communication" name="profile" checked>
      <label for="communication">Communication about your account.</label>
    </li>
    <li>
      <input type="checkbox" id="profile" name="profile" checked>
      <label for="profile">Associate with profile data. When you login into 3rd party apps with this account, you can choose to share your email.</label>
    </li>
  </ul>
  <div  class="self-end mt-2 space-x-2">
    <button type="button" on:click={submit} class="btn btn-link">
      Skip
    </button>
    <button type="button" on:click={submit} class="btn btn-primary">
      Register
    </button>
  </div>
</form>
