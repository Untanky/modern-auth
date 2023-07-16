<script lang="ts">
</script>

<div class="flex w-screen h-screen justify-center items-center">
  <main class="w-[420px] card">
    <h1 class="text-stone-950 dark:text-stone-50 text-xl font-medium">
      Login
    </h1>
    <p class="mt-2">
      Please submit your user id.
    </p>
    <form class="flex flex-col flex-1">
      <div class="relative mt-2">
        <label class="dark:bg-stone-800 px-1 left-3 absolute block" for="user-id">
          User Id:
        </label>
        <input class="dark:bg-stone-800 mt-3.5 px-4 py-2 dark:border-stone-600 border rounded-lg w-full" type="text" id="user-id" placeholder="Your user id" />
      </div>
      <button class="self-end mt-2 btn btn-primary">
        Submit
      </button>
    </form>
  </main>
</div>
