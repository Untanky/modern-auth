<script>
  import Footer from "$lib/Footer.svelte";
  import "../app.postcss";
</script>

<div class="flex flex-col w-screen h-screen justify-center items-end sm:items-center p-4">
  <slot></slot>
  <Footer />
</div>