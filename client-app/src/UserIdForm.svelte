<script lang="ts">
  export let submit: () => void;
</script>

<form class="flex flex-col flex-1">
  <p>
    Please submit your user id.
  </p>
  <div class="relative mt-2">
    <label class="dark:bg-stone-800 px-1 left-3 absolute block" for="user-id">
      User Id:
    </label>
    <input class="dark:bg-stone-800 mt-3.5 px-4 py-2 dark:border-stone-600 border rounded-lg w-full" type="text" id="user-id" placeholder="Your user id" />
  </div>
  <button type="button" on:click={submit} class="self-end mt-2 btn btn-primary">
    Submit
  </button>
</form>
