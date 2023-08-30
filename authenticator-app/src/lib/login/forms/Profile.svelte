<script lang="ts">
  import ProfileForm from '$lib/forms/ProfileForm.svelte';

  export let submit: () => void;
</script>

<form class="flex flex-col flex-1 space-y-1">
  <h2 class="text-xl">
    Profile
  </h2>
  <p>
    When you login to 3rd party apps through this account, you can select to share this information with the app. You may complete your profile:
  </p>
  <ProfileForm />
  <div class="self-end mt-2 space-x-2">
    <button type="button" on:click={submit} class="btn btn-link">
      Skip
    </button>
    <button type="button" on:click={submit} class="btn btn-yellow">
      Register
    </button>
  </div>
</form>
