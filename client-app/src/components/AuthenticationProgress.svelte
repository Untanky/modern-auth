<script>
  import DoneIcon from './icons/Done.svelte';
  import IdentificationIcon from './icons/Identification.svelte';
  import PersonIcon from './icons/Person.svelte';

</script>

<div class="relative py-2">
  <div class="flex justify-between">
    <div class="p-px w-[58px] h-[58px] rounded-full bg-stone-500 dark:bg-stone-400"/>
    <div class="p-px w-[58px] h-[58px] rounded-full bg-stone-500 dark:bg-stone-400"/>
    <div class="p-px w-[58px] h-[58px] rounded-full bg-stone-500 dark:bg-stone-400"/>
  </div>
  <div class="absolute top-[34px] left-4 right-4 p-px bg-stone-500 dark:bg-stone-400">
    <div class="h-1.5 bg-stone-50 dark:bg-stone-900"></div>
  </div>
  <div class="absolute left-0 right-0 top-2 flex justify-between">
    <div class="p-px">
      <div class="p-2 w-14 h-14 fill-stone-950 bg-stone-50 dark:fill-stone-50 dark:bg-stone-900 rounded-full">
        <PersonIcon />
      </div>
    </div>
    <div class="p-px">
      <div class="p-2 w-14 h-14 fill-stone-950 bg-stone-50 dark:fill-stone-50 dark:bg-stone-900 rounded-full">
        <IdentificationIcon />
      </div>
    </div>
    <div class="p-px">
      <div class="p-2 w-14 h-14 fill-stone-950 bg-stone-50 dark:fill-stone-50 dark:bg-stone-900 rounded-full">
        <DoneIcon />
      </div>
    </div>
  </div>
</div>
