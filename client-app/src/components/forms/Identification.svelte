<script lang="ts">
  export let submit: (userId: string) => void;

  const onClick = (event: Event) => {
    event.preventDefault();

    const inputElement = document.getElementById('user-id') as HTMLInputElement;

    submit(inputElement.value);
  }
</script>

<form class="flex flex-col flex-1">
  <h2 class="text-xl">
    Identification
  </h2>
  <p>
    To authenticate, we need to know who you are. Please enter your user id.
  </p>
  <div class="relative mt-2">
    <label class="dark:bg-stone-800 px-1 left-3 absolute block" for="user-id">
      User Id:
    </label>
    <input class="dark:bg-stone-800 mt-3.5 px-4 py-2 dark:border-stone-600 border rounded-lg w-full" type="text" id="user-id" placeholder="Your user id" />
  </div>
  <button type="button" on:click={onClick} class="self-end mt-2 btn btn-primary">
    Submit
  </button>
</form>
