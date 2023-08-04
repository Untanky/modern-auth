<script>
  import SmartLogin from "./components/SmartLogin.svelte";
  import Identification from "./components/icons/Identification.svelte";

</script>
<div class="flex flex-col w-screen h-screen justify-center items-end sm:items-center p-4">
  <main class="relative sm:w-[420px] card space-y-2">
    <div class="absolute -top-10 left-0 right-0 flex justify-center">
      <span class="w-20 h-20 p-2 border bg-stone-50 dark:bg-stone-900 dark:border-stone-600 rounded-full">
      </span>
    </div>
    <div class="absolute top-[-8px] h-10 left-24 right-24 bg-stone-50 dark:bg-stone-900 flex justify-center fill-stone-950 dark:fill-stone-50">
      <div class="w-20 h-20 p-2 -mt-10"><Identification /></div>
    </div>
    <SmartLogin />    
  </main>
  <footer class="absolute text-sm bottom-4 text-stone-500 dark:text-stone-400 w-full px-4 flex justify-between">
    <nav>
      <a href="/imprint" class="rounded mx-2 my-1 underline">Imprint</a>
      -
      <a href="/data-privacy" class="rounded mx-2 my-1 underline">Data privacy</a>
      -
      <a href="/terms" class="rounded mx-2 my-1 underline">Terms</a>
    </nav>
    <div>
      Version: Development
    </div>
  </footer>
</div>
