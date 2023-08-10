<script>
  import SmartLogin from "$lib/login/SmartLogin.svelte";
  import Identification from "$lib/login/icons/Identification.svelte";
</script>

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
