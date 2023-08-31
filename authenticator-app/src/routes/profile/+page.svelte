<script lang="ts">
    import CardWithIcon from '$lib/CardWithIcon.svelte';
    import type { Profile } from '$lib/profile/model';
    import Person from '$lib/login/icons/Person.svelte';
    import { onMount } from 'svelte';
    import Loading from '$lib/utils/Loading.svelte';
    import ProfileForm from './ProfileForm.svelte';

    let profile: Promise<Profile> = new Promise(() => {});

    const loadProfile = (): Promise<Profile> => fetch('/v1/profile')
            .then((res) => res.json())

    const updateProfile = (pref: Profile): void => {
        fetch('/v1/profile', {
            body: JSON.stringify(pref),
            method: 'PUT',
        })
            .then(loadProfile)
            .then((result) => profile = Promise.resolve(result));
    }

    onMount(() => {
        profile = loadProfile();
    });
</script>

<svelte:head>
    <title>Profile</title>
</svelte:head>

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
    {#await profile}
        <Loading />
    {:then p} 
        <ProfileForm
            profile={p}
            onSave={updateProfile}
        />
    {/await}
  </CardWithIcon>
</main>
