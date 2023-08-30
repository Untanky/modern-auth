<script lang="ts">
  import CardWithIcon from '$lib/CardWithIcon.svelte';
  import Person from '$lib/login/icons/Person.svelte';
  import ProfileForm from '$lib/forms/ProfileForm.svelte';

  const onSubmit = (): void => {

  };
</script>

<main class="sm:w-[420px]">
  <CardWithIcon>
    <svelte:fragment slot="icon">
      <Person />
    </svelte:fragment>
    <h1>Profile</h1>
    <p>
      The information you enter here, can be requested by third parties and
      is only <em>shared with your permission</em>.
    </p>
    <form class="flex flex-col" on:submit={onSubmit}>
      <ProfileForm />
      <button type="submit" class="ml-auto btn btn-yellow">Save</button>
    </form>
  </CardWithIcon>
</main>
