<script>
  import SmartLogin from "./components/SmartLogin.svelte";

</script>
<div class="flex flex-col w-screen h-screen justify-center items-end sm:items-center p-4">
  <main class="sm:w-[420px] card space-y-2">
    <SmartLogin />    
  </main>
  <footer class="absolute text-sm bottom-4 text-stone-500 dark:text-stone-400 w-full px-4 flex justify-between">
    <nav>
      <a href="/imprint" class="rounded px-2 py-1 underline">Imprint</a>
      -
      <a href="/data-privacy" class="rounded px-2 py-1 underline">Data privacy</a>
      -
      <a href="/terms" class="rounded px-2 py-1 underline">Terms</a>
    </nav>
    <div>
      Version: Development
    </div>
  </footer>
</div>
