<script lang="ts">
    import CardWithIcon from '$lib/CardWithIcon.svelte';
    import Email from '$lib/login/icons/Email.svelte';
    import type { InsertPreferences, Preferences } from '$lib/preferences/model';
    import Loading from '$lib/utils/Loading.svelte';
    import { onMount } from 'svelte';
    import CommunicationForm from './CommunicationForm.svelte';

    export let preferences: Promise<Preferences> = new Promise(() => {});

    const loadPreferences = (): Promise<Preferences> => fetch('/v1/communication/preferences')
            .then((res) => res.json())
    
    const updatePreferences = (pref: InsertPreferences): void => {
        fetch('/v1/communication/preferences', {
            body: JSON.stringify(pref),
            method: 'PUT',
        })
            .then(loadPreferences)
            .then((result) => preferences = Promise.resolve(result));
    }

    onMount(() => {
        preferences = loadPreferences();
    });
</script>

<svelte:head>
    <title>Communication Preferences</title>
</svelte:head>

<main class="sm:w-[420px]">
  <CardWithIcon class="rounded-b-md">
    <svelte:fragment slot="icon">
      <Email />
    </svelte:fragment>
    <h1 class="relative z-100">Communication</h1>
    <h4>Update communication data and preferences</h4>
    {#await preferences}
        <Loading />
    {:then pref} 
        <CommunicationForm
            preferences={pref}
            onSave={updatePreferences}
        />
    {/await}
  </CardWithIcon>
  <p class="card mt-2 card-blue rounded-t-md">
    The information you enter here, will  <em>never be shared 
    with third parties</em>. The email entered here is only
    used for communication relating to you account, such as
    account recovery.
  </p>
</main>
