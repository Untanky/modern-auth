<script lang="ts">
  export let submit: (userId: string) => void;

  let userId = '';

  $: canSubmit = userId.length > 0;

  const onClick = (event: SubmitEvent) => {
    event.preventDefault();

    submit(event.target[0].value);
  }
</script>

<form class="flex flex-col flex-1" on:submit={onClick}>
  <h2 class="text-xl">
    Identification
  </h2>
  <p class="mt-2">
    To authenticate, we need to know who you are. Please enter your user id.
  </p>
  <div class="relative mt-4">
    <label class="font-extralight text-sm uppercase dark:bg-stone-800 px-1 left-3 absolute block rounded" for="user-id">
      User ID:
    </label>
    <input
      bind:value={userId}
      class="dark:bg-stone-800 mt-3 px-4 py-2 dark:border-stone-600 border rounded-lg w-full focus:outline-none focus:border-yellow-500"
      type="text"
      id="user-id"
      placeholder="Your user id"
    />
  </div>
  <button type="submit" class="self-end mt-4 btn btn-primary" disabled={!canSubmit}>
    Continue
  </button>
</form>
