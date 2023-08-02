<script>
  import DoneIcon from './icons/Done.svelte';
  import IdentificationIcon from './icons/Identification.svelte';
  import PersonIcon from './icons/Person.svelte';

  export let progress = 0;

  const gradient = `bg-gradient-to-r from-green-400 via-yellow-400 to-stone-50`;
  const fullStone = 'bg-stone-50';
  const halfYellow = `${gradient} from-0% via-0% to-50%`;
  const fullGreen = 'bg-green-400';

  $: first = progress === 0 ? halfYellow : fullGreen;
  $: second = progress < 0.33333 ? fullStone : progress > 0.33333 ? fullGreen : halfYellow;
  $: third = progress < 0.66667 ? fullStone : progress > 0.66667 ? fullGreen : halfYellow;
</script>

<div class="relative py-2">
  <div class="flex justify-between">
    <div class={`p-px w-[58px] h-[58px] rounded-full ${first}`}/>
    <div class={`p-px w-[58px] h-[58px] rounded-full ${second}`}/>
    <div class={`p-px w-[58px] h-[58px] rounded-full ${third}`}/>
  </div>
  <div class={`absolute top-[34px] left-4 right-4 p-px progressive-bg ${gradient} from-0% via-[${progress * 100}%] to-[${progress * 100}%]`}>
    <div class="h-1.5 dark:bg-stone-800"></div>
  </div>
  <div class="absolute left-0 right-0 top-2 flex justify-between">
    <div class="p-px">
      <div class="p-2 w-14 h-14 dark:fill-stone-50 dark:bg-stone-800 rounded-full">
        <PersonIcon />
      </div>
    </div>
    <div class="p-px">
      <div class="p-2 w-14 h-14 dark:fill-stone-50 dark:bg-stone-800 rounded-full">
        <IdentificationIcon />
      </div>
    </div>
    <div class="p-px">
      <div class="p-2 w-14 h-14 dark:fill-stone-50 dark:bg-stone-800 rounded-full">
        <DoneIcon />
      </div>
    </div>
  </div>
</div>
