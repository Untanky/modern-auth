<script>

  import CardWithIcon from "$lib/CardWithIcon.svelte";
  import Identification from "$lib/login/icons/Identification.svelte";

</script>

<main class="sm:w-[420px]">
  <CardWithIcon>
    <svelte:fragment slot="icon">
      <Identification />
    </svelte:fragment>
    <h1>
      Register
    </h1>
    <p>
      Create new a new user acocunt
    </p>
    <p>
      You can choose any unique user identifier you want (your name, username, email).
      The identifier will be securely stored. Please enter your unique user identifier:
    </p>
  </CardWithIcon>
</main>
