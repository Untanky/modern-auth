<script>
  import { browser } from "$app/environment";
  import Footer from "$lib/Footer.svelte";
  import { initializeStoreLocally } from "$lib/authorization/store";
  import { onMount } from "svelte";
  import "../app.postcss";

  onMount(() => {
    if (browser) {
      initializeStoreLocally();
    }
  });
</script>

<div class="flex flex-col w-screen h-screen justify-center items-end sm:items-center p-4">
  <slot></slot>
  <Footer />
</div>