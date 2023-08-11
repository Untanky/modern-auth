<script>
  import { browser } from "$app/environment";
  import CardWithIcon from "$lib/CardWithIcon.svelte";
  import { initializeStoreLocally } from "$lib/authorization/store";
  import SmartLogin from "$lib/login/SmartLogin.svelte";
  import Identification from "$lib/login/icons/Identification.svelte";
  import { onMount } from "svelte";

  onMount(() => {
    if (browser) {
      initializeStoreLocally();
    }
  });
</script>

<main class="sm:w-[420px]">
  <CardWithIcon>
    <svelte:fragment slot="icon">
      <Identification />
    </svelte:fragment>
    <SmartLogin />    
  </CardWithIcon>
</main>
