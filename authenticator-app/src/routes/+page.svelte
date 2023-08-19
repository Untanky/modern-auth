<script>
  import CardWithIcon from "$lib/CardWithIcon.svelte";
  import SmartLogin from "$lib/login/SmartLogin.svelte";
  import Identification from "$lib/login/icons/Identification.svelte";
</script>

<main class="sm:w-[420px]">
  <CardWithIcon>
    <svelte:fragment slot="icon">
      <Identification />
    </svelte:fragment>
    <SmartLogin />    
  </CardWithIcon>
</main>
<a href="/communication">Communication</a>
